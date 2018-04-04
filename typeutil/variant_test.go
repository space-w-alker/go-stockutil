package typeutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariant(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, Variant{}.String())
	assert.Equal(map[Variant]Variant{
		V(`test`):  V(1),
		V(`other`): V(2.4),
	}, V(map[string]interface{}{
		`test`:  1,
		`other`: 2.4,
	}).Map())
}
