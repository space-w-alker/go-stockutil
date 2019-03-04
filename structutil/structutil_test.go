package structutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TBase struct {
	Name    string
	Age     int
	Enabled bool
}

type tPerson struct {
	*TBase
	Address string
}

func TestCopyFunc(t *testing.T) {
	assert := require.New(t)

	dest := tPerson{
		TBase: &TBase{
			Enabled: true,
		},
	}

	src := tPerson{
		TBase: &TBase{
			Name: `Bob Johnson`,
			Age:  42,
		},
		Address: `123 Fake St.`,
	}

	assert.NoError(CopyNonZero(&dest, &src))
	assert.Equal(tPerson{
		Address: `123 Fake St.`,
		TBase: &TBase{
			Age:     42,
			Enabled: true,
			Name:    `Bob Johnson`,
		},
	}, dest)
}
