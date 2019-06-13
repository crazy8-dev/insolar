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

package genesisrefs

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/rootdomain"
	"github.com/stretchr/testify/require"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		insolar.GetGenesisNameRootDomain(): {
			got:    ContractRootDomain,
			expect: "1tJEDNVffdf4PQjxhKvQVc3D166RqhmFBS2gkBpGva.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameNodeDomain(): {
			got:    ContractNodeDomain,
			expect: "1tJCaZ7rFeUncXXcYoKFTMza6xypkF8BGGXh48X2Fy.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameNodeRecord(): {
			got:    ContractNodeRecord,
			expect: "1tJCZvWMHXqs4Yk2E1YJFXiETMHzAWzfQu2qK5XYpA.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameRootMember(): {
			got:    ContractRootMember,
			expect: "1tJDPyrf6x41Pe3npMWpjgtge1LmHcb6gFwaJA5oF5.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameRootWallet(): {
			got:    ContractRootWallet,
			expect: "1tJDw8Rqitda9ofpYRQ3DEytsdAFAPH2BBKqUKtzLn.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameDeposit(): {
			got:    ContractDeposit,
			expect: "1tJCxDbykDzGCA83wR9LPeALt4LaM5aNqgrAkiM2Ly.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameTariff(): {
			got:    ContractTariff,
			expect: "1tJDKMVqKhAJufMF2ioE43L57oFC6jeDXVRpre9qH2.11111111111111111111111111111111",
		},
		insolar.GetGenesisNameCostCenter(): {
			got:    ContractCostCenter,
			expect: "1tJDB7kbbc7vg8uX1N7FZSWCgK2YFxbt4U9Xdc5oL6.11111111111111111111111111111111",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}

func TestRootDomain(t *testing.T) {
	ref1 := rootdomain.RootDomain.Ref()
	ref2 := rootdomain.GenesisRef(insolar.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}
