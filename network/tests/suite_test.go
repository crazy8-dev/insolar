//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

// +build networktest

package tests

import (
	"context"
	"crypto"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/node"

	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/keystore"

	"github.com/insolar/insolar/network/servicenetwork"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

var (
	testNetworkPort uint32 = 10000
)

var (
	suiteLogger = inslogger.FromContext(initLogger(context.Background(), insolar.ErrorLevel))
)

const (
	UseFakeTransport = true
	UseFakeBootstrap = true

	reqTimeoutMs     int32 = 2000
	pulseDelta       int32 = 1
	consensusMin           = 5 // minimum count of participants that can survive when one node leaves
	maxPulsesForJoin       = 3
)

const cacheDir = "network_cache/"

func initLogger(ctx context.Context, level insolar.LogLevel) context.Context {
	logger := inslogger.FromContext(ctx).WithCaller(false)
	logger, _ = logger.WithLevelNumber(level)
	logger, _ = logger.WithFormat(insolar.TextFormat)
	ctx = inslogger.SetLogger(ctx, logger)
	return ctx
}

// testSuite is base test suite
type testSuite struct {
	bootstrapCount int
	nodesCount     int
	ctx            context.Context
	bootstrapNodes []*networkNode
	networkNodes   []*networkNode
	pulsar         TestPulsar
	t              *testing.T

	discoveriesAreBootstrapped uint32
}

type consensusSuite struct {
	testSuite
}

func newTestSuite(t *testing.T, bootstrapCount, nodesCount int) testSuite {
	return testSuite{
		bootstrapCount: bootstrapCount,
		nodesCount:     nodesCount,
		t:              t,
		ctx:            initLogger(inslogger.TestContext(t), insolar.DebugLevel),
		bootstrapNodes: make([]*networkNode, 0),
		networkNodes:   make([]*networkNode, 0),
	}
}

func newConsensusSuite(t *testing.T, bootstrapCount, nodesCount int) *consensusSuite {
	if bootstrapCount < consensusMin {
		panic("incorrect bootstrapCount, it should 5 or more")
	}

	return &consensusSuite{
		testSuite: newTestSuite(t, bootstrapCount, nodesCount),
	}
}

// SetupSuite creates and run network with bootstrap and common nodes once before run all tests in the suite
func (s *consensusSuite) SetupTest() {
	var err error
	s.pulsar, err = NewTestPulsar(reqTimeoutMs, pulseDelta)
	require.NoError(s.t, err)

	suiteLogger.Info("SetupTest")

	for i := 0; i < s.bootstrapCount; i++ {
		role := insolar.StaticRoleVirtual
		if i == 0 {
			role = insolar.StaticRoleHeavyMaterial
		}
		s.bootstrapNodes = append(s.bootstrapNodes, s.newNetworkNodeWithRole(fmt.Sprintf("bootstrap_%d", i), role))
	}

	for i := 0; i < s.nodesCount; i++ {
		s.networkNodes = append(s.networkNodes, s.newNetworkNode(fmt.Sprintf("node_%d", i)))
	}

	pulseReceivers := make([]string, 0)
	for _, n := range s.bootstrapNodes {
		pulseReceivers = append(pulseReceivers, n.host)
	}

	suiteLogger.Info("Setup bootstrap nodes")
	s.SetupNodesNetwork(s.bootstrapNodes)
	if UseFakeBootstrap {
		bnodes := make([]insolar.NetworkNode, 0)
		for _, n := range s.bootstrapNodes {
			o := n.serviceNetwork.NodeKeeper.GetOrigin()
			dig, sig := o.(node.MutableNode).GetSignature()
			require.NotNil(s.t, dig)
			require.NotNil(s.t, sig.Bytes())

			bnodes = append(bnodes, o)
		}
		for _, n := range s.bootstrapNodes {
			n.serviceNetwork.ConsensusMode = consensus.ReadyNetwork
			n.serviceNetwork.NodeKeeper.SetInitialSnapshot(bnodes)
			err := n.serviceNetwork.PulseAppender.AppendPulse(s.ctx, *insolar.GenesisPulse)
			require.NoError(s.t, err)
			n.serviceNetwork.Gatewayer.SwitchState(s.ctx, insolar.CompleteNetworkState, *insolar.GenesisPulse)
			pulseReceivers = append(pulseReceivers, n.host)
		}
	}

	s.StartNodesNetwork(s.bootstrapNodes)

	expectedBootstrapsCount := len(s.bootstrapNodes)
	retries := 10
	for {
		activeNodes := s.bootstrapNodes[0].GetActiveNodes()
		if expectedBootstrapsCount == len(activeNodes) {
			break
		}

		retries--
		if retries == 0 {
			break
		}

		time.Sleep(2 * time.Second)
	}

	activeNodes := s.bootstrapNodes[0].GetActiveNodes()
	require.Equal(s.t, len(s.bootstrapNodes), len(activeNodes))

	if len(s.networkNodes) > 0 {
		suiteLogger.Info("Setup network nodes")
		s.SetupNodesNetwork(s.networkNodes)
		s.StartNodesNetwork(s.networkNodes)

		s.waitForConsensus(2)

		// active nodes count verification
		activeNodes1 := s.networkNodes[0].GetActiveNodes()
		activeNodes2 := s.networkNodes[0].GetActiveNodes()

		require.Equal(s.t, s.getNodesCount(), len(activeNodes1))
		require.Equal(s.t, s.getNodesCount(), len(activeNodes2))
	}
	suiteLogger.Info("Start test pulsar")
	err = s.pulsar.Start(initLogger(s.ctx, insolar.ErrorLevel), pulseReceivers)
	require.NoError(s.t, err)
}

func (s *testSuite) waitResults(results chan error, expected int) {
	count := 0
	for count < expected {
		err := <-results
		require.NoError(s.t, err)
		count++
	}
}

func (s *testSuite) SetupNodesNetwork(nodes []*networkNode) {
	for _, n := range nodes {
		s.preInitNode(n)
	}

	results := make(chan error, len(nodes))
	initNode := func(node *networkNode) {
		err := node.init()
		results <- err
	}

	suiteLogger.Info("Init nodes")
	for _, n := range nodes {
		go initNode(n)
	}
	s.waitResults(results, len(nodes))
}

func (s *testSuite) StartNodesNetwork(nodes []*networkNode) {
	suiteLogger.Info("Start nodes")

	results := make(chan error, len(nodes))
	startNode := func(node *networkNode) {
		err := node.componentManager.Start(node.ctx)
		node.serviceNetwork.RegisterConsensusFinishedNotifier(func(ctx context.Context, report network.Report) {
			node.consensusResult <- report.PulseNumber
		})
		results <- err
	}

	for _, n := range nodes {
		go startNode(n)
	}
	s.waitResults(results, len(nodes))
	atomic.StoreUint32(&s.discoveriesAreBootstrapped, 1)
}

// stopNetwork shutdowns all nodes in network
func (s *consensusSuite) stopNetwork() {
	suiteLogger.Info("=================== stopNetwork()")
	suiteLogger.Info("Stop network nodes")
	for _, n := range s.networkNodes {
		err := n.componentManager.Stop(n.ctx)
		require.NoError(s.t, err)
	}
	suiteLogger.Info("Stop bootstrap nodes")
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Stop(n.ctx)
		require.NoError(s.t, err)
	}
	suiteLogger.Info("Stop test pulsar")
	err := s.pulsar.Stop(s.ctx)
	require.NoError(s.t, err)
}
func (s *consensusSuite) waitForNodeJoin(ref insolar.Reference, pulsesCount int) bool {
	for i := 0; i < pulsesCount; i++ {
		pn := s.waitForConsensus(1)
		if s.isNodeInActiveLists(ref, pn) {
			return true
		}
	}
	return false
}

func (s *consensusSuite) waitForConsensus(consensusCount int) insolar.PulseNumber {
	var pn insolar.PulseNumber
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.bootstrapNodes {
			select {
			case pn = <-n.consensusResult:
				continue
			case <-time.After(time.Second * 12):
				panic("waitForConsensus timeout")
			}
		}

		for _, n := range s.networkNodes {
			<-n.consensusResult
		}
	}
	return pn
}

func (s *consensusSuite) waitForConsensusExcept(consensusCount int, exception insolar.Reference) {
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.bootstrapNodes {
			if n.id.Equal(exception) {
				continue
			}
			<-n.consensusResult
		}

		for _, n := range s.networkNodes {
			if n.id.Equal(exception) {
				continue
			}
			<-n.consensusResult
		}
	}
}

// nodesCount returns count of nodes in network without testNode
func (s *testSuite) getNodesCount() int {
	return len(s.bootstrapNodes) + len(s.networkNodes)
}

func (s *testSuite) isNodeInActiveLists(ref insolar.Reference, p insolar.PulseNumber) bool {
	for _, n := range s.bootstrapNodes {
		a := n.serviceNetwork.NodeKeeper.GetAccessor(p)
		if a.GetActiveNode(ref) == nil {
			return false
		}
	}
	return true
}

func (s *testSuite) InitNode(node *networkNode) {
	if node.componentManager != nil {
		err := node.init()
		require.NoError(s.t, err)
	}
}

func (s *testSuite) StartNode(node *networkNode) {
	if node.componentManager != nil {
		err := node.componentManager.Start(node.ctx)
		require.NoError(s.t, err)
	}
}

func (s *testSuite) StopNode(node *networkNode) {
	if node.componentManager != nil {
		err := node.componentManager.Stop(s.ctx)
		require.NoError(s.t, err)
	}
}

func (s *testSuite) GracefulStop(node *networkNode) {
	if node.componentManager != nil {
		err := node.componentManager.GracefulStop(s.ctx)
		require.NoError(s.t, err)

		err = node.componentManager.Stop(s.ctx)
		require.NoError(s.t, err)
	}
}

type networkNode struct {
	id                  insolar.Reference
	role                insolar.StaticRole
	privateKey          crypto.PrivateKey
	cryptographyService insolar.CryptographyService
	host                string
	ctx                 context.Context

	componentManager   *component.Manager
	serviceNetwork     *servicenetwork.ServiceNetwork
	terminationHandler *testutils.TerminationHandlerMock
	consensusResult    chan insolar.PulseNumber
}

func (s *testSuite) newNetworkNode(name string) *networkNode {
	return s.newNetworkNodeWithRole(name, insolar.StaticRoleVirtual)
}
func (s *testSuite) startNewNetworkNode(name string) *networkNode {
	testNode := s.newNetworkNode(name)
	s.preInitNode(testNode)
	s.InitNode(testNode)
	s.StartNode(testNode)
	return testNode
}

// newNetworkNode returns networkNode initialized only with id, host address and key pair
func (s *testSuite) newNetworkNodeWithRole(name string, role insolar.StaticRole) *networkNode {
	key, err := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NoError(s.t, err)
	address := "127.0.0.1:" + strconv.Itoa(incrementTestPort())

	n := &networkNode{
		id:                  gen.Reference(),
		role:                role,
		privateKey:          key,
		cryptographyService: cryptography.NewKeyBoundCryptographyService(key),
		host:                address,
		consensusResult:     make(chan insolar.PulseNumber, 1),
	}

	nodeContext, _ := inslogger.WithFields(s.ctx, map[string]interface{}{
		"node_name": name,
	})

	n.ctx = nodeContext
	return n
}

func incrementTestPort() int {
	result := atomic.AddUint32(&testNetworkPort, 1)
	return int(result)
}

// init calls Init for node component manager and wraps PhaseManager
func (n *networkNode) init() error {
	err := n.componentManager.Init(n.ctx)
	// n.serviceNetwork.PhaseManager = &phaseManagerWrapper{original: n.serviceNetwork.PhaseManager, result: n.consensusResult}
	return err
}

func (n *networkNode) GetActiveNodes() []insolar.NetworkNode {
	p, err := n.serviceNetwork.PulseAccessor.GetLatestPulse(n.ctx)
	if err != nil {
		panic(err)
	}
	return n.serviceNetwork.NodeKeeper.GetAccessor(p.PulseNumber).GetActiveNodes()
}

func (n *networkNode) GetWorkingNodes() []insolar.NetworkNode {
	p, err := n.serviceNetwork.PulseAccessor.GetLatestPulse(n.ctx)
	if err != nil {
		panic(err)
	}
	return n.serviceNetwork.NodeKeeper.GetAccessor(p.PulseNumber).GetWorkingNodes()
}

func (s *testSuite) initCrypto(node *networkNode) (*certificate.CertificateManager, insolar.CryptographyService) {
	pubKey, err := node.cryptographyService.GetPublicKey()
	require.NoError(s.t, err)

	// init certificate

	proc := platformpolicy.NewKeyProcessor()
	publicKey, err := proc.ExportPublicKeyPEM(pubKey)
	require.NoError(s.t, err)

	cert := &certificate.Certificate{}
	cert.PublicKey = string(publicKey[:])
	cert.Reference = node.id.String()
	cert.Role = node.role.String()
	cert.BootstrapNodes = make([]certificate.BootstrapNode, 0)
	cert.MinRoles.HeavyMaterial = 1
	cert.MinRoles.Virtual = 4

	for _, b := range s.bootstrapNodes {
		pubKey, _ := b.cryptographyService.GetPublicKey()
		pubKeyBuf, err := proc.ExportPublicKeyPEM(pubKey)
		require.NoError(s.t, err)

		bootstrapNode := certificate.NewBootstrapNode(
			pubKey,
			string(pubKeyBuf[:]),
			b.host,
			b.id.String(),
			b.role.String(),
		)

		sign, err := certificate.SignCert(b.cryptographyService, cert.PublicKey, cert.Role, cert.Reference)
		require.NoError(s.t, err)
		bootstrapNode.NodeSign = sign.Bytes()

		cert.BootstrapNodes = append(cert.BootstrapNodes, *bootstrapNode)
	}

	// dump cert and read it again from json for correct private files initialization
	jsonCert, err := cert.Dump()
	require.NoError(s.t, err)

	cert, err = certificate.ReadCertificateFromReader(pubKey, proc, strings.NewReader(jsonCert))
	require.NoError(s.t, err)
	return certificate.NewCertificateManager(cert), node.cryptographyService
}

type PublisherMock struct{}

func (p *PublisherMock) Publish(topic string, messages ...*message.Message) error {
	return nil
}

func (p *PublisherMock) Close() error {
	return nil
}

// preInitNode inits previously created node with mocks and external dependencies
func (s *testSuite) preInitNode(node *networkNode) {
	cfg := configuration.NewConfiguration()
	cfg.Pulsar.PulseTime = pulseDelta * 1000
	cfg.Host.Transport.Address = node.host
	cfg.Service.CacheDirectory = cacheDir + node.host

	node.componentManager = &component.Manager{}
	node.componentManager.Register(platformpolicy.NewPlatformCryptographyScheme())
	serviceNetwork, err := servicenetwork.NewServiceNetwork(cfg, node.componentManager)
	require.NoError(s.t, err)

	certManager, cryptographyService := s.initCrypto(node)

	realKeeper, err := nodenetwork.NewNodeNetwork(cfg.Host.Transport, certManager.GetCertificate())
	require.NoError(s.t, err)
	terminationHandler := testutils.NewTerminationHandlerMock(s.t)
	terminationHandler.LeaveMock.Set(func(p context.Context, p1 insolar.PulseNumber) {})
	terminationHandler.OnLeaveApprovedMock.Set(func(p context.Context) {})
	terminationHandler.AbortMock.Set(func(reason string) { log.Error(reason) })

	keyProc := platformpolicy.NewKeyProcessor()
	pubMock := &PublisherMock{}
	if UseFakeTransport {
		// little hack: this Register will override transport.Factory
		// in servicenetwork internal component manager with fake factory
		node.componentManager.Register(transport.NewFakeFactory(cfg.Host.Transport))
	} else {
		node.componentManager.Register(transport.NewFactory(cfg.Host.Transport))
	}

	pulseManager := testutils.NewPulseManagerMock(s.t)
	pulseManager.SetMock.Set(func(ctx context.Context, pulse insolar.Pulse) (err error) {
		return nil
	})
	node.componentManager.Inject(
		realKeeper,
		pulseManager,
		pubMock,
		certManager,
		cryptographyService,
		keystore.NewInplaceKeyStore(node.privateKey),
		serviceNetwork,
		keyProc,
		terminationHandler,
		testutils.NewContractRequesterMock(s.t),
		// pulse.NewStorageMem(),
	)
	node.serviceNetwork = serviceNetwork
	node.terminationHandler = terminationHandler

	nodeContext, _ := inslogger.WithFields(s.ctx, map[string]interface{}{
		"node_id":      realKeeper.GetOrigin().ShortID(),
		"node_address": realKeeper.GetOrigin().Address(),
		"node_role":    realKeeper.GetOrigin().Role().String(),
	})

	node.ctx = nodeContext
}

func (s *testSuite) AssertActiveNodesCountDelta(delta int) {
	activeNodes := s.bootstrapNodes[1].GetActiveNodes()
	require.Equal(s.t, s.getNodesCount()+delta, len(activeNodes))
}

func (s *testSuite) AssertWorkingNodesCountDelta(delta int) {
	workingNodes := s.bootstrapNodes[0].GetWorkingNodes()
	require.Equal(s.t, s.getNodesCount()+delta, len(workingNodes))
}
