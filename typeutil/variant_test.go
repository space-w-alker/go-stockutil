package typeutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariant(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, Variant{}.String())
}
