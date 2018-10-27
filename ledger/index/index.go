/*
 *    Copyright 2018 Insolar
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

package index

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/record"
)

// ClassLifeline represents meta information for record object
type ClassLifeline struct {
	LatestState *record.ID // Amend or activate record
	State       record.State
}

// ObjectLifeline represents meta information for record object.
type ObjectLifeline struct {
	ClassRef            record.Reference
	LatestState         *record.ID // Amend or activate record.
	LatestStateApproved *record.ID // State approved by VM.
	ChildPointer        *record.ID // Meta record about child activation.
	Delegates           map[core.RecordRef]record.Reference
	State               record.State
}
