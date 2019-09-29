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

package logicrunner

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/metrics"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.RequestFetcher -o ./ -s _mock.go -g

type RequestFetcher interface {
	FetchPendings(ctx context.Context, trs chan<- *common.Transcript)
	Abort(ctx context.Context)
}

type requestFetcher struct {
	object insolar.Reference

	stopFetching func()

	broker           ExecutionBrokerI
	artifactsManager artifacts.Client
	outgoingsSender  OutgoingRequestSender
}

func NewRequestsFetcher(
	obj insolar.Reference, am artifacts.Client, br ExecutionBrokerI, os OutgoingRequestSender,
) RequestFetcher {
	return &requestFetcher{
		object:           obj,
		broker:           br,
		artifactsManager: am,
		outgoingsSender:  os,
		stopFetching:     func() {},
	}
}

func (rf *requestFetcher) Abort(ctx context.Context) {
	rf.stopFetching()
}

func (rf *requestFetcher) FetchPendings(ctx context.Context, trs chan<- *common.Transcript) {
	defer close(trs)

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	rf.stopFetching = cancelFunc

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"object": rf.object.String(),
	})

	logger.Debug("request fetcher starting")

	err := rf.fetch(ctx, trs)
	if err != nil {
		logger.Error("couldn't make fetch round: ", err.Error())
	}
}

func (rf *requestFetcher) fetch(ctx context.Context, trs chan<- *common.Transcript) error {
	logger := inslogger.FromContext(ctx)

	for {
		stats.Record(ctx, metrics.RequestFetcherFetchCall.M(1))
		reqRefs, err := rf.artifactsManager.GetPendings(ctx, rf.object)
		if err != nil {
			if err == insolar.ErrNoPendingRequest {
				logger.Debug("no more pendings on ledger")
				rf.broker.NoMoreRequestsOnLedger(ctx)
				return nil
			}
			return err
		}

		for _, reqRef := range reqRefs {
			if !reqRef.IsRecordScope() {
				logger.Errorf("skipping request with bad reference, ref=%s", reqRef.String())
				continue
			}

			logger.Debug("getting request from ledger")
			stats.Record(ctx, metrics.RequestFetcherFetchUnique.M(1))
			request, err := rf.artifactsManager.GetAbandonedRequest(ctx, rf.object, reqRef)
			if err != nil {
				return errors.Wrap(err, "couldn't get request")
			}

			select {
			case <-ctx.Done():
				logger.Debug("request fetcher stopping")
				return nil
			default:
			}

			switch v := request.(type) {
			case *record.IncomingRequest:
				if err := checkIncomingRequest(ctx, v); err != nil {
					err = errors.Wrap(err, "failed to check incoming request")
					logger.Error(err.Error())

					continue
				}
				tr := common.NewTranscriptCloneContext(ctx, reqRef, *v)
				trs <- tr
			case *record.OutgoingRequest:
				if err := checkOutgoingRequest(ctx, v); err != nil {
					err = errors.Wrap(err, "failed to check outgoing request")
					logger.Error(err.Error())

					continue
				}
				// FIXME: limit there may slow down things, placing "go" here is not good too
				rf.outgoingsSender.SendAbandonedOutgoingRequest(ctx, reqRef, v)
			default:
				logger.Error("requestFetcher fetched unknown request")
			}
		}
	}
}
