// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package args

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNil(t *testing.T) {
	require.True(t, IsNil(nil))
	require.False(t, IsNil("test"))

	var v testHolder

	require.False(t, v.Get() == nil)
	require.True(t, IsNil(v.Get()))

	d := 0
	v.value = &d

	require.False(t, v.Get() == nil)
	require.False(t, IsNil(v.Get()))

	v.value = nil

	require.False(t, v.Get() == nil)
	require.True(t, IsNil(v.Get()))
}

type testHolder struct {
	value *int
}

func (v testHolder) Get() interface{} {
	return v.value
}
