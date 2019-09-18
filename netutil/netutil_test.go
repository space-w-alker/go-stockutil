package netutil

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestFQDN(t *testing.T) {
	assert := require.New(t)
	sys, err := exec.Command(`hostname`, `-f`).Output()
	assert.NoError(err)
	assert.Equal(strings.TrimSpace(string(sys)), FQDN())
}
