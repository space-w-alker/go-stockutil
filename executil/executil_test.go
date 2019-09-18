package executil

import (
	"os/exec"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestWhich(t *testing.T) {
	assert := require.New(t)

	assert.Equal(`/bin/ls`, Which(`/bin/ls`))
	assert.Equal(`/bin/ls`, Which(`ls`))
	assert.Equal(`/usr/bin/tail`, Which(`tail`))
	assert.Empty(Which(`absolutely-not-a-command-@#$%^&*`))
}

func TestJoin(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, Join(nil))
	assert.Equal(``, Join(``))
	assert.Equal(``, Join([]string{}))
	assert.Equal(`ls -l`, Join([]string{`ls`, `-l`}))
	assert.Equal(`ls -l '/this is a folder'`, Join([]string{`ls`, `-l`, `/this is a folder`}))

	assert.Equal(`ls -l '/this is a folder'`, Join(
		exec.Command(`ls`, `-l`, `/this is a folder`),
	))

	assert.Equal(`whoami`, Join(
		exec.Command(`whoami`),
	))

	assert.Equal(``, Join(
		exec.Command(``),
	))
}
