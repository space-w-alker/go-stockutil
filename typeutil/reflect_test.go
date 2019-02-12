package typeutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctionArity(t *testing.T) {
	assert := require.New(t)

	i, o, err := FunctionArity(strings.Compare)
	assert.NoError(err)
	assert.Equal(2, i)
	assert.Equal(1, o)

	f0_0 := func() {}
	i, o, err = FunctionArity(f0_0)
	assert.NoError(err)
	assert.Equal(0, i)
	assert.Equal(0, o)

	f0_1 := func() error { return nil }
	i, o, err = FunctionArity(f0_1)
	assert.NoError(err)
	assert.Equal(0, i)
	assert.Equal(1, o)
}
