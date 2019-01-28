/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
)

type middleware struct {
	objectStorage  storage.ObjectStorage
	jetStorage     storage.JetStorage
	jetCoordinator core.JetCoordinator
	messageBus     core.MessageBus
	pulseStorage   core.PulseStorage
	hotDataWaiter  HotDataWaiter
	conf           *configuration.Ledger
	handler        *MessageHandler
	seqMutex       sync.Mutex
	sequencer      map[core.RecordID]*struct {
		sync.Mutex
		done bool
	}
}

func newMiddleware(
	h *MessageHandler,
) *middleware {
	return &middleware{
		objectStorage:  h.ObjectStorage,
		jetStorage:     h.JetStorage,
		jetCoordinator: h.JetCoordinator,
		messageBus:     h.Bus,
		pulseStorage:   h.PulseStorage,
		hotDataWaiter:  h.HotDataWaiter,
		handler:        h,
		conf:           h.conf,
		sequencer: map[core.RecordID]*struct {
			sync.Mutex
			done bool
		}{},
	}
}

// In Build it should follow after checkJet, since it expects jet key
func (m *middleware) addFieldsToLogger(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		context, _ := inslogger.WithField(ctx, "target", parcel.DefaultTarget().String())

		val := ctx.Value(jetKey{})
		j, ok := val.(core.RecordID)
		if ok {
			context, _ = inslogger.WithField(context, "jetid", j.DebugString())
		}

		return handler(context, parcel)
	}
}

type jetKey struct{}

func contextWithJet(ctx context.Context, jetID core.RecordID) context.Context {
	return context.WithValue(ctx, jetKey{}, jetID)
}

func jetFromContext(ctx context.Context) core.RecordID {
	val := ctx.Value(jetKey{})
	j, ok := val.(core.RecordID)
	if !ok {
		panic("failed to extract jet from context")
	}

	return j
}

func (m *middleware) zeroJetForHeavy(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		return handler(contextWithJet(ctx, *jet.NewID(0, nil)), parcel)
	}
}

func (m *middleware) checkJet(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		msg := parcel.Message()
		if msg.DefaultTarget() == nil {
			return nil, errors.New("unexpected message")
		}

		logger := inslogger.FromContext(ctx)
		logger.Debugf("checking jet for %v", parcel.Type().String())

		// FIXME: @andreyromancev. 17.01.19. Temporary allow any genesis request. Remove it.
		if parcel.Pulse() == core.FirstPulseNumber {
			logger.Debugf("genesis pulse shortcut")
			return handler(contextWithJet(ctx, *jet.NewID(0, nil)), parcel)
		}

		// Check token jet.
		token := parcel.DelegationToken()
		if token != nil {
			logger.Debugf("received token. returning any jet")
			// Calculate jet for target pulse.
			target := *msg.DefaultTarget().Record()
			pulse := target.Pulse()
			switch tm := msg.(type) {
			case *message.GetObject:
				pulse = tm.State.Pulse()
			case *message.GetChildren:
				if tm.FromChild == nil {
					return nil, errors.New("fetching children without child pointer is forbidden")
				}
				pulse = tm.FromChild.Pulse()
			}
			tree, err := m.jetStorage.GetJetTree(ctx, pulse)
			if err != nil {
				return nil, err
			}

			jetID, actual := tree.Find(target)
			if !actual {
				inslogger.FromContext(ctx).Errorf(
					"got message of type %s with redirect token,"+
						" but jet %s for pulse %d is not actual",
					msg.Type(), jetID.DebugString(), pulse,
				)
			}

			return handler(contextWithJet(ctx, *jetID), parcel)
		}

		// Calculate jet for current pulse.
		var jetID core.RecordID
		if msg.DefaultTarget().Record().Pulse() == core.PulseNumberJet {
			logger.Debugf("special pulse number (jet). returning jet from message")
			jetID = *msg.DefaultTarget().Record()
		} else {
			j, actual, err := m.fetchJet(ctx, *msg.DefaultTarget().Record(), parcel.Pulse())
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch jet tree")
			}
			if !actual {
				return &reply.JetMiss{JetID: *j}, nil
			}
			jetID = *j
		}

		// Check if jet is ours.
		node, err := m.jetCoordinator.LightExecutorForJet(ctx, jetID, parcel.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate executor for jet")
		}
		if *node != m.jetCoordinator.Me() {
			return &reply.JetMiss{JetID: jetID}, nil
		}

		return handler(contextWithJet(ctx, jetID), parcel)
	}
}

func (m *middleware) saveParcel(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		jetID := jetFromContext(ctx)
		pulse, err := m.pulseStorage.Current(ctx)
		if err != nil {
			return nil, err
		}
		logger.Debugf("saveParcel, pulse - %v", pulse.PulseNumber)
		err = m.objectStorage.SetMessage(ctx, jetID, pulse.PulseNumber, parcel)
		if err != nil {
			return nil, err
		}

		return handler(ctx, parcel)
	}
}

func (m *middleware) checkHeavySync(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		// TODO: @andreyromancev. 10.01.2019. Uncomment to enable backpressure for writing requests.
		// Currently disabled due to big initial difference in pulse numbers, which prevents requests from being accepted.
		// jetID := jetFromContext(ctx)
		// replicated, err := m.db.GetReplicatedPulse(ctx, jetID)
		// if err != nil {
		// 	return nil, err
		// }
		// if parcel.Pulse()-replicated >= m.conf.LightChainLimit {
		// 	return nil, errors.New("failed to write data (waiting for heavy replication)")
		// }

		return handler(ctx, parcel)
	}
}

func (m *middleware) fetchJet(
	ctx context.Context, target core.RecordID, pulse core.PulseNumber,
) (*core.RecordID, bool, error) {
	// Look in the local tree. Return if the actual jet found.
	tree, err := m.jetStorage.GetJetTree(ctx, pulse)
	if err != nil {
		return nil, false, err
	}
	jetID, actual := tree.Find(target)
	if actual {
		inslogger.FromContext(ctx).Debugf(
			"we believe object %s is in JET %s", target.String(), jetID.DebugString(),
		)
		return jetID, actual, nil
	}

	inslogger.FromContext(ctx).Debugf(
		"jet %s is not actual in our tree, asking neighbors for jet of object %s",
		jetID.DebugString(), target.String(),
	)

	m.seqMutex.Lock()
	if _, ok := m.sequencer[*jetID]; !ok {
		m.sequencer[*jetID] = &struct {
			sync.Mutex
			done bool
		}{}
	}
	mu := m.sequencer[*jetID]
	m.seqMutex.Unlock()

	mu.Lock()
	if mu.done {
		mu.Unlock()
		inslogger.FromContext(ctx).Debugf(
			"somebody else updated actuality of jet %s, rechecking our DB",
			jetID.DebugString(),
		)
		return m.fetchJet(ctx, target, pulse)
	}
	defer func() {
		inslogger.FromContext(ctx).Debugf("done fetching jet, cleaning")

		mu.done = true
		mu.Unlock()

		m.seqMutex.Lock()
		inslogger.FromContext(ctx).Debugf("deleting sequencer for jet %s", jetID.DebugString())
		delete(m.sequencer, *jetID)
		m.seqMutex.Unlock()
	}()

	resJet, err := m.handler.fetchActualJetFromOtherNodes(ctx, target, pulse)
	if err != nil {
		return nil, false, err
	}

	err = m.jetStorage.UpdateJetTree(ctx, pulse, true, *resJet)
	if err != nil {
		inslogger.FromContext(ctx).Error(
			errors.Wrapf(err, "couldn't actualize jet %s", resJet.DebugString()),
		)
	}

	return resJet, true, nil
}

func (m *middleware) waitForHotData(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[waitForHotData] for parcel with pulse %v", parcel.Pulse())

		// TODO: 15.01.2019 @egorikas
		// Hack is needed for genesis
		if parcel.Pulse() == core.FirstPulseNumber {
			return handler(ctx, parcel)
		}

		// If the call is a call in redirect-chain
		// skip waiting for the hot records
		if parcel.DelegationToken() != nil {
			logger.Debugf("[waitForHotData] parcel.DelegationToken() != nil")
			return handler(ctx, parcel)
		}

		jetID := jetFromContext(ctx)
		err := m.hotDataWaiter.Wait(ctx, jetID)
		if err != nil {
			return &reply.Error{ErrType: reply.ErrHotDataTimeout}, nil
		}
		return handler(ctx, parcel)
	}
}

func (m *middleware) releaseHotDataWaiters(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)
		logger.Debugf("[releaseHotDataWaiters] pulse %v starts %v", parcel.Pulse(), time.Now())

		hotDataMessage := parcel.Message().(*message.HotData)
		jetID := hotDataMessage.Jet.Record()

		logger.Debugf("[releaseHotDataWaiters] hot data for jet happens - %v, pulse - %v", jetID.DebugString(), parcel.Pulse())
		defer m.hotDataWaiter.Unlock(ctx, *jetID)

		logger.Debugf("[releaseHotDataWaiters] before handler for jet - %v, pulse - %v", jetID.DebugString(), parcel.Pulse())
		return handler(ctx, parcel)
	}
}
