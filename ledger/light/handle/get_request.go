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

package handle

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetRequest struct {
	dep     *proc.Dependencies
	message payload.Meta
	passed  bool
}

func NewGetRequest(dep *proc.Dependencies, msg payload.Meta, passed bool) *GetRequest {
	return &GetRequest{
		dep:     dep,
		message: msg,
		passed:  passed,
	}
}

func (s *GetRequest) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetRequest message")
	}
	msg, ok := pl.(*payload.GetRequest)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(msg.ObjectID, msg.RequestID.Pulse(), s.message, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}

	passIfNotFound := !s.passed
	req := proc.NewGetRequest(s.message, msg.ObjectID, msg.RequestID, passIfNotFound)
	s.dep.GetRequest(req)
	return f.Procedure(ctx, req, false)
}
