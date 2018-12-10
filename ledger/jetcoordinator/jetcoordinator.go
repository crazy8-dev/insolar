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

package jetcoordinator

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/utils/entropy"
	"github.com/pkg/errors"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	db                         *storage.DB
	roleCounts                 map[core.DynamicRole]int
	NodeNet                    core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(db *storage.DB, conf configuration.JetCoordinator) *JetCoordinator {
	jc := JetCoordinator{db: db}
	jc.loadConfig(conf)

	return &jc
}

func (jc *JetCoordinator) loadConfig(conf configuration.JetCoordinator) {
	jc.roleCounts = map[core.DynamicRole]int{}

	for intRole, count := range conf.RoleCounts {
		role := core.DynamicRole(intRole)
		jc.roleCounts[role] = count
	}
}

// IsAuthorized checks for role on concrete pulse for the address.
func (jc *JetCoordinator) IsAuthorized(
	ctx context.Context,
	role core.DynamicRole,
	obj *core.RecordRef,
	pulse core.PulseNumber,
	node core.RecordRef,
) (bool, error) {
	nodes, err := jc.QueryRole(ctx, role, obj, pulse)
	if err != nil {
		return false, err
	}
	for _, n := range nodes {
		if n == node {
			return true, nil
		}
	}
	return false, nil
}

// QueryRole returns node refs responsible for role bound operations for given object and pulse.
func (jc *JetCoordinator) QueryRole(
	ctx context.Context,
	role core.DynamicRole,
	obj *core.RecordRef,
	pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	pulseData, err := jc.db.GetPulse(ctx, pulse)
	if err != nil {
		return nil, err
	}
	candidates := jc.NodeNet.GetActiveNodesByRole(role)
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no candidates for role %d", role))
	}
	count, ok := jc.roleCounts[role]
	if !ok {
		return nil, errors.New("no candidate count for this role")
	}
	ent := pulseData.Pulse.Entropy[:]

	if obj == nil {
		return getRefs(jc.PlatformCryptographyScheme, ent, candidates, count)
	}

	objHash := obj.Record().Hash()
	if role == core.DynamicRoleLightExecutor {
		jetTree, err := jc.db.GetJetTree(ctx, pulseData.Pulse.PulseNumber)
		if err != nil {
			return nil, err
		}
		_, depth := jetTree.Find(objHash, pulseData.Pulse.PulseNumber)

		// Reset everything except prefix.
		resetBits(objHash, depth+1)
		circleXOR(ent, objHash)
		return getRefs(jc.PlatformCryptographyScheme, ent, candidates, count)
	}

	circleXOR(ent, objHash)
	return getRefs(jc.PlatformCryptographyScheme, ent, candidates, count)
}

// GetActiveNodes return active nodes for specified pulse.
func (jc *JetCoordinator) GetActiveNodes(pulse core.PulseNumber) ([]core.Node, error) {
	return jc.db.GetActiveNodes(pulse)
}

func getRefs(
	scheme core.PlatformCryptographyScheme,
	e []byte,
	values []core.RecordRef,
	count int,
) ([]core.RecordRef, error) {
	// TODO: remove sort when network provides sorted result from GetActiveNodesByRole (INS-890) - @nordicdyno 5.Dec.2018
	sort.SliceStable(values, func(i, j int) bool {
		return bytes.Compare(values[i][:], values[j][:]) < 0
	})
	in := make([]interface{}, 0, len(values))
	for _, value := range values {
		in = append(in, interface{}(value))
	}

	res, err := entropy.SelectByEntropy(scheme, e, in, count)
	if err != nil {
		return nil, err
	}
	out := make([]core.RecordRef, 0, len(res))
	for _, value := range res {
		out = append(out, value.(core.RecordRef))
	}
	return out, nil
}

// CircleXOR performs XOR for 'dst' and 'src'. The result is written into 'dst'.
// If 'src' is smaller than 'dst', it starts from the beginning of 'src'.
func circleXOR(dst, src []byte) {
	srcLen := len(src)
	for i := range dst {
		dst[i] ^= src[i%srcLen]
	}
}

// ResetBits resets all bits in 'value', starting from 'start' number of bit.
func resetBits(value []byte, start int) {
	startByte := start / 8
	startBit := start % 8

	if start > len(value)*8 {
		return
	}

	// Reset bits in starting byte.
	for i := startBit; i < 8; i++ {
		mask := byte(^(1 << byte(i)))
		value[startByte] &= mask
	}

	if startByte >= len(value)-1 {
		return
	}

	// Reset bytes.
	for i := startByte + 1; i < len(value); i++ {
		value[i] = 0
	}
}
