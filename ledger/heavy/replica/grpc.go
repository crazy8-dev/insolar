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

package replica

import (
	context "context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"

	"github.com/insolar/insolar/insolar"
)

type grpcTransport struct {
	lis        net.Listener
	grpcServer *grpc.Server
	handlers   map[string]Handle
}

func NewGRPCTransport() Transport {
	return &grpcTransport{
		handlers: make(map[string]Handle),
	}
}

func (t *grpcTransport) Init(ctx context.Context) error {
	t.grpcServer = grpc.NewServer()
	RegisterReplicaServer(t.grpcServer, t)
	return nil
}

func (t *grpcTransport) Start(ctx context.Context) error {
	const port = 20111
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "failed to open replication port %d", port)
	}
	t.lis = lis
	return t.grpcServer.Serve(t.lis)
}

func (t *grpcTransport) Call(ctx context.Context, request *Request) (*Response, error) {
	method := request.Method
	data := request.Data
	if _, ok := t.handlers[method]; !ok {
		return nil, errors.Errorf("handle function: %v not defined", method)
	}
	result, err := t.handlers[method](data)
	resError, err := insolar.Serialize(err)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to serialize error value")
	}
	return &Response{Data: result, Error: resError}, nil
}

func (t *grpcTransport) Send(ctx context.Context, receiver, method string, data []byte) ([]byte, error) {
	req := Request{Method: method, Data: data}
	conn, err := grpc.Dial(receiver, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to receiver %s", receiver)
	}
	client := NewReplicaClient(conn)

	res, err := client.Call(ctx, &req)
	if err != nil || res == nil {
		return nil, errors.Wrapf(err, "failed to call RPC method %v", method)
	}
	resError := error(nil)
	err = insolar.Deserialize(res.Error, &resError)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deserialize error value")
	}
	return res.Data, resError
}

func (t *grpcTransport) Register(method string, handle Handle) {
	t.handlers[method] = handle
}

func (t *grpcTransport) Me() string {
	return t.lis.Addr().String()
}
