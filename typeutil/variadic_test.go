package typeutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariadic(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, Variadic{}.String())
}
