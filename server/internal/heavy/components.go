//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package heavy

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/network"

	"github.com/insolar/insolar/ledger/heavy/exporter"
	"google.golang.org/grpc"

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/server/internal"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/genesis"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/genesisdataprovider"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/handler"
	"github.com/insolar/insolar/ledger/heavy/pulsemanager"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/platformpolicy"
)

type components struct {
	cmp         component.Manager
	NodeRef     string
	NodeRole    string
	rollback    *executor.DBRollback
	stateKeeper *executor.InitialStateKeeper
	inRouter    *watermillMsg.Router
	outRouter   *watermillMsg.Router

	replicator executor.HeavyReplicator
}

func newComponents(ctx context.Context, cfg configuration.Configuration, genesisCfg insolar.GenesisHeavyConfig) (*components, error) {
	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
		CertManager   insolar.CertificateManager
	)
	{
		var err error
		// Private key storage.
		ks, err := keystore.NewKeyStore(cfg.KeysPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load KeyStore")
		}
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()

		c := component.Manager{}
		c.Inject(CryptoService, CryptoScheme, KeyProcessor, ks)

		publicKey, err := CryptoService.GetPublicKey()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve node public key")
		}

		// Node certificate.
		CertManager, err = certificate.NewManagerReadCertificate(publicKey, KeyProcessor, cfg.CertificatePath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Certificate")
		}
	}

	c := &components{}
	c.cmp = component.Manager{}
	c.NodeRef = CertManager.GetCertificate().GetNodeRef().String()
	c.NodeRole = CertManager.GetCertificate().GetRole().String()

	logger := inslogger.FromContext(ctx)

	// Watermill stuff.
	var (
		wmLogger   *log.WatermillLogAdapter
		publisher  watermillMsg.Publisher
		subscriber watermillMsg.Subscriber
	)
	{
		wmLogger = log.NewWatermillLogAdapter(logger)
		pubsub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		subscriber = pubsub
		publisher = pubsub
		// Wrapped watermill publisher for introspection.
		publisher = internal.PublisherWrapper(ctx, &c.cmp, cfg.Introspection, publisher)
	}

	// Network.
	var (
		NetworkService *servicenetwork.ServiceNetwork
		NodeNetwork    network.NodeNetwork
		Termination    insolar.TerminationHandler
	)
	{
		var err error
		// External communication.
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg, &c.cmp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}

		Termination = termination.NewHandler(NetworkService)

		// Node info.
		NodeNetwork, err = nodenetwork.NewNodeNetwork(cfg.Host.Transport, CertManager.GetCertificate())
		if err != nil {
			return nil, errors.Wrap(err, "failed to start NodeNetwork")
		}

	}

	// API.
	var (
		Requester       insolar.ContractRequester
		GenesisProvider insolar.GenesisDataProvider
		API             insolar.APIRunner
	)
	{
		var err error
		Requester, err = contractrequester.New(nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ContractRequester")
		}

		GenesisProvider, err = genesisdataprovider.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start GenesisDataProvider")
		}

		API, err = api.NewRunner(&cfg.APIRunner)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ApiRunner")
		}
	}

	// Storage.
	var (
		Coordinator jet.Coordinator
		Pulses      *pulse.DB
		Nodes       *node.Storage
		DB          *store.BadgerDB
		Jets        *jet.DBStore
	)
	{
		var err error
		fullDataDirectoryPath, err := filepath.Abs(cfg.Ledger.Storage.DataDirectory)
		if err != nil {
			panic(errors.Wrap(err, "failed to get absolute path for DataDirectory"))
		}
		DB, err = store.NewBadgerDB(badger.DefaultOptions(fullDataDirectoryPath))
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize DB"))
		}
		Nodes = node.NewStorage()
		Pulses = pulse.NewDB(DB)
		Jets = jet.NewDBStore(DB)

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.OriginProvider = NodeNetwork
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	// Communication.
	var (
		Tokens  insolar.DelegationTokenFactory
		Parcels message.ParcelFactory
		Bus     insolar.MessageBus
		WmBus   *bus.Bus
	)
	{
		var err error
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Parcels = messagebus.NewParcelFactory()
		Bus, err = messagebus.NewMessageBus(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start MessageBus")
		}
		WmBus = bus.NewBus(cfg.Bus, publisher, Pulses, Coordinator, CryptoScheme)
	}

	metricsHandler, err := metrics.NewMetrics(
		ctx,
		cfg.Metrics,
		metrics.GetInsolarRegistry(c.NodeRole),
		c.NodeRole,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start Metrics")
	}

	var (
		PulseManager insolar.PulseManager
		Handler      *handler.Handler
		Genesis      *genesis.Genesis
		Records      *object.RecordDB
		JetKeeper    executor.JetKeeper
	)
	{
		Records = object.NewRecordDB(DB)
		indexes := object.NewIndexDB(DB, Records)
		drops := drop.NewDB(DB)
		JetKeeper = executor.NewJetKeeper(Jets, DB, Pulses)
		c.rollback = executor.NewDBRollback(JetKeeper, Pulses, drops, Records, indexes, Jets, Pulses)
		c.stateKeeper = executor.NewInitialStateKeeper(JetKeeper, Jets, Coordinator, indexes, drops)

		sp := pulse.NewStartPulse()

		backupMaker, err := executor.NewBackupMaker(ctx, DB, cfg.Ledger.Backup, JetKeeper.TopSyncPulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed create backuper")
		}

		pm := pulsemanager.NewPulseManager()
		pm.NodeNet = NodeNetwork
		pm.NodeSetter = Nodes
		pm.Nodes = Nodes
		pm.PulseAppender = Pulses
		pm.PulseAccessor = Pulses
		pm.JetModifier = Jets
		pm.StartPulse = sp
		pm.FinalizationKeeper = executor.NewFinalizationKeeperDefault(JetKeeper, Pulses, cfg.Ledger.LightChainLimit)

		replicator := executor.NewHeavyReplicatorDefault(Records, indexes, CryptoScheme, Pulses, drops, JetKeeper, backupMaker, Jets)
		c.replicator = replicator

		h := handler.New(cfg.Ledger)
		h.RecordAccessor = Records
		h.RecordModifier = Records
		h.JetCoordinator = Coordinator
		h.IndexAccessor = indexes
		h.IndexModifier = indexes
		h.DropModifier = drops
		h.PCS = CryptoScheme
		h.PulseAccessor = Pulses
		h.PulseCalculator = Pulses
		h.StartPulse = sp
		h.JetModifier = Jets
		h.JetAccessor = Jets
		h.JetTree = Jets
		h.DropDB = drops
		h.JetKeeper = JetKeeper
		h.InitialStateReader = c.stateKeeper
		h.BackupMaker = backupMaker
		h.Sender = WmBus
		h.Replicator = replicator

		PulseManager = pm
		Handler = h

		artifactManager := &artifact.Scope{
			PulseNumber:    insolar.FirstPulseNumber,
			PCS:            CryptoScheme,
			RecordAccessor: Records,
			RecordModifier: Records,
			IndexModifier:  indexes,
			IndexAccessor:  indexes,
		}
		Genesis = &genesis.Genesis{
			ArtifactManager: artifactManager,
			BaseRecord: &genesis.BaseRecord{
				DB:             DB,
				DropModifier:   drops,
				PulseAppender:  Pulses,
				PulseAccessor:  Pulses,
				RecordModifier: Records,
				IndexModifier:  indexes,
			},

			DiscoveryNodes:  genesisCfg.DiscoveryNodes,
			ContractsConfig: genesisCfg.ContractsConfig,
		}
	}

	// Exporter
	var (
		recordExporter *exporter.RecordServer
		pulseExporter  *exporter.PulseServer
	)
	{
		recordExporter = exporter.NewRecordServer(Pulses, Records, Records, JetKeeper)
		pulseExporter = exporter.NewPulseServer(Pulses, JetKeeper)

		grpcServer := grpc.NewServer()
		exporter.RegisterRecordExporterServer(grpcServer, recordExporter)
		exporter.RegisterPulseExporterServer(grpcServer, pulseExporter)

		lis, err := net.Listen("tcp", cfg.Exporter.Addr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open port for Exporter")
		}

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				panic(fmt.Errorf("exporter failed to serve: %s", err))
			}
		}()
	}

	c.cmp.Inject(
		DB,
		WmBus,
		Handler,
		PulseManager,
		Jets,
		Pulses,
		Coordinator,
		metricsHandler,
		Bus,
		Requester,
		Tokens,
		Parcels,
		artifacts.NewClient(WmBus),
		GenesisProvider,
		API,
		KeyProcessor,
		Termination,
		CryptoScheme,
		CryptoService,
		CertManager,
		NodeNetwork,
		NetworkService,
		publisher,
	)
	err = c.cmp.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init components")
	}

	if !genesisCfg.Skip {
		if err := Genesis.Start(ctx); err != nil {
			logger.Fatalf("genesis failed on heavy with error: %v", err)
		}
	}

	c.startWatermill(ctx, wmLogger, subscriber, WmBus, NetworkService.SendMessageHandler, Handler.Process)

	return c, nil
}

func (c *components) Start(ctx context.Context) error {
	err := c.rollback.Start(ctx)
	if err != nil {
		return errors.Wrapf(err, "rollback.Start return error: %s", err.Error())
	}

	err = c.stateKeeper.Start(ctx)
	if err != nil {
		return errors.Wrapf(err, "stateKeeper.Start return error: %s", err.Error())
	}
	return c.cmp.Start(ctx)
}

func (c *components) Stop(ctx context.Context) error {
	err := c.inRouter.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error("Error while closing router", err)
	}
	err = c.outRouter.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error("Error while closing router", err)
	}
	c.replicator.Stop()
	return c.cmp.Stop(ctx)
}

func (c *components) startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	sub watermillMsg.Subscriber,
	b *bus.Bus,
	outHandler, inHandler watermillMsg.NoPublishHandlerFunc,
) {
	inRouter, err := watermillMsg.NewRouter(watermillMsg.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	outRouter, err := watermillMsg.NewRouter(watermillMsg.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	outRouter.AddNoPublisherHandler(
		"OutgoingHandler",
		bus.TopicOutgoing,
		sub,
		outHandler,
	)

	inRouter.AddMiddleware(
		middleware.InstantAck,
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		sub,
		inHandler,
	)

	startRouter(ctx, inRouter)
	c.inRouter = inRouter
	startRouter(ctx, outRouter)
	c.outRouter = outRouter
}

func startRouter(ctx context.Context, router *watermillMsg.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
