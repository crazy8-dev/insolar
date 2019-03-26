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

package adapter

import (
	"fmt"
	"time"

	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// NewWaitAdapter creates new instance of SimpleWaitAdapter with Waiter as worker
func NewWaitAdapter(id adapterid.ID) TaskSink {
	return NewAdapterWithQueue(NewWaiter(id), id)
}

// WaiterTask is task for adapter for waiting
type WaiterTask struct {
	waitPeriodMilliseconds int
}

// Waiter is worker for adapter for waiting
type Waiter struct{}

// NewWaiter returns new instance of worker which waiting
func NewWaiter(id adapterid.ID) Processor {
	return &Waiter{}
}

// Process implements Processor interface
func (w *Waiter) Process(task AdapterTask, nestedEventHelper NestedEventHelper, cancelInfo CancelInfo) interface{} {
	log.Info("[ Waiter.Process ] Start.")

	payload, ok := task.TaskPayload.(WaiterTask)
	var msg interface{}

	if !ok {
		msg = errors.Errorf("[ Waiter.Process ] Incorrect payload type: %T", task.TaskPayload)
		return msg
	}

	time.Sleep(time.Duration(payload.waitPeriodMilliseconds) * time.Millisecond)
	msg = fmt.Sprintf("Work completed successfully. Waited %d millisecond", payload.waitPeriodMilliseconds)
	log.Info("[ Waiter.Process ] ", msg)

	return msg
}
