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

package functest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMembersWithSameName(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       "NameForTestCreateMembersWithSameName",
	})

	memberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, memberResponse)

	firstMemberRef := memberResponse.Reference
	assert.NotEqual(t, "", firstMemberRef)

	body = getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       "NameForTestCreateMembersWithSameName",
	})

	unmarshalResponse(t, body, memberResponse)

	secondMemberRef := memberResponse.Reference
	assert.NotEqual(t, "", secondMemberRef)

	assert.NotEqual(t, firstMemberRef, secondMemberRef)
}
