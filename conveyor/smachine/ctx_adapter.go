///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package smachine

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common/syncrun"
)

var _ ExecutionAdapter = &adapterExecHelper{}

type adapterExecHelper struct {
	adapterID AdapterID
	executor  AdapterExecutor
}

func (p *adapterExecHelper) IsEmpty() bool {
	return p.adapterID.IsEmpty()
}

func (p *adapterExecHelper) GetAdapterID() AdapterID {
	return p.adapterID
}

func (p *adapterExecHelper) PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallContext {
	return &adapterCallContext{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterSyncCallContext}
}

func (p *adapterExecHelper) PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) CallContext {
	return &adapterCallContext{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterAsyncCallContext}
}

const (
	adapterSyncCallContext     = 1
	adapterAsyncCallContext    = 2
	adapterCallContextDisposed = 3
)

type adapterCallContext struct {
	ctx      *executionContext
	fn       AdapterCallFunc
	executor AdapterExecutor
	mode     uint8

	stepBound bool
	cancel    *syncrun.ChainedCancel
}

func (c *adapterCallContext) discard() {
	c.mode = adapterCallContextDisposed
}

func (c *adapterCallContext) ensureMode(mode uint8) {
	if c.mode != mode {
		panic("illegal state")
	}
}

func (c *adapterCallContext) GetCancel(fn *context.CancelFunc) CallContext {
	if c.cancel != nil {
		*fn = c.cancel.Cancel
		return c
	}

	r := *c
	r.cancel = syncrun.NewChainedCancel()
	*fn = r.cancel.Cancel
	return &r
}

func (c *adapterCallContext) CancelOnStep(attach bool) CallContext {
	r := *c
	r.stepBound = attach
	return &r
}

func (c *adapterCallContext) Start() {
	c.ensureMode(adapterAsyncCallContext)
	defer c.discard()

	c._startAsync()
}

func (c *adapterCallContext) Wait() CallConditionalUpdate {
	c.ensureMode(adapterAsyncCallContext)
	defer c.discard()

	return &conditionalUpdate{marker: &c.ctx.marker, kickOff: func(*Slot) {
		c._startAsync()
	}}
}

func (c *adapterCallContext) TryCall() bool {
	c.ensureMode(adapterSyncCallContext)
	defer c.discard()

	return c._startSync()
}

func (c *adapterCallContext) Call() {
	c.ensureMode(adapterSyncCallContext)
	defer c.discard()

	if !c._startSync() {
		panic("call was cancelled")
	}
}

func (c *adapterCallContext) _startAsync() {
	var stepLink StepLink
	if c.stepBound {
		stepLink = c.ctx.s.NewExactStepLink()
	} else {
		stepLink = c.ctx.s.NewAnyStepLink()
	}

	if c.cancel != nil && c.cancel.IsCancelled() {
		return
	}

	c.ctx.countAsyncCalls++

	cancelFn := c.executor.StartCall(stepLink, c.fn, func(fn AsyncResultFunc, recovered interface{}) {
		c.ctx.machine.applyAsyncStateUpdate(stepLink, fn, recovered)
	}, c.cancel != nil)

	if c.cancel != nil {
		c.cancel.SetChain(cancelFn)
	}
}

func (c *adapterCallContext) _startSync() bool {
	resultFn := c._startSyncWithResult()

	if resultFn == nil {
		return false
	}

	rc := asyncResultContext{slot: c.ctx.s}
	rc.executeResult(resultFn)
	return true
}

func (c *adapterCallContext) _startSyncWithResult() AsyncResultFunc {

	if ok, result := c.executor.TrySyncCall(c.fn); ok {
		return result
	}

	stepLink := c.ctx.s.NewExactStepLink()
	ok, wc := c.ctx.worker.GetCond()
	if !ok {
		return nil
	}

	var resultFn AsyncResultFunc
	var resultRecovered interface{}
	var callState int

	cancelFn := c.executor.StartCall(stepLink, c.fn, func(fn AsyncResultFunc, recovered interface{}) {
		wc.L.Lock()
		if callState == 0 {
			resultFn = fn
			resultRecovered = recovered
			callState = 1
			wc.Broadcast()
		}
		wc.L.Unlock()
	}, false)

	wc.L.Lock()
	if callState == 0 {
		wc.Wait()

		if callState == 0 {
			/* Cond can be triggered by Worker for emergent stop */
			callState = 2

			wc.L.Unlock()

			if cancelFn != nil {
				cancelFn()
			}
			return nil
		}
	}
	wc.L.Unlock()

	if resultRecovered != nil {
		panic(resultRecovered)
	}
	return resultFn
}