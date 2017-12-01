package rxutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	assert := require.New(t)

	match := Match(`1234.5678.9`, `(?P<first>\d+)\.(\d+).(?P<second>\d+)`)
	assert.NotNil(match)

	assert.Equal(`1234.5678.9`, match.Group(0))
	assert.Equal(`1234`, match.Group(1))
	assert.Equal(`5678`, match.Group(2))
	assert.Equal(`9`, match.Group(3))
	assert.Empty(match.Group(4))

	assert.Equal(`1234`, match.Group(`first`))
	assert.Equal(`9`, match.Group(`second`))
	assert.Empty(match.Group(`potato`))
}

func TestMatchAndMap(t *testing.T) {
	assert := require.New(t)

	match := Match(`1234.5678.9`, `(?P<first>\d+)\.(\d+).(?P<second>\d+)`)
	assert.NotNil(match)

	assert.Equal(`1234.5678.9`, match.Group(0))
	assert.Equal(map[string]string{
		`first`:  `1234`,
		`second`: `9`,
	}, match.NamedCaptures())
}
