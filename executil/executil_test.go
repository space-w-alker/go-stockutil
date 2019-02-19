package executil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhich(t *testing.T) {
	assert := require.New(t)

	assert.Equal(`/bin/ls`, Which(`/bin/ls`))
	assert.Equal(`/bin/ls`, Which(`ls`))
	assert.Equal(`/usr/bin/tail`, Which(`tail`))
	assert.Empty(Which(`absolutely-not-a-command-@#$%^&*`))
}
