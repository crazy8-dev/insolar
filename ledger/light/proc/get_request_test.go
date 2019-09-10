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

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
)

func TestGetRequest_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		records     *object.RecordAccessorMock
		sender      *bus.SenderMock
		coordinator *jet.CoordinatorMock
		fetcher     *executor.JetFetcherMock
	)

	setup := func() {
		records = object.NewRecordAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
		fetcher = executor.NewJetFetcherMock(mc)
	}

	t.Run("Not passing, returns error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		meta := payload.Meta{}

		p := proc.NewGetRequest(meta, gen.ID(), gen.ID(), false)
		p.Dep(records, sender, coordinator, fetcher)

		err := p.Proceed(ctx)
		assert.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		assert.True(t, ok)
		assert.Equal(t, uint32(payload.CodeNotFound), insError.GetCode())
	})

	t.Run("Simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		rec := record.Material{
			Virtual:  record.Wrap(&record.IncomingRequest{}),
			ID:       gen.ID(),
			ObjectID: gen.ID(),
		}
		records.ForIDMock.Return(rec, nil)

		reqID := gen.ID()
		expectedMsg, _ := payload.NewMessage(&payload.Request{
			RequestID: reqID,
			Request:   rec.Virtual,
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
		}).Return()

		p := proc.NewGetRequest(payload.Meta{}, gen.ID(), reqID, true)
		p.Dep(records, sender, coordinator, fetcher)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})
}
