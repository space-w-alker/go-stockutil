package rxutil

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestSplitN(t *testing.T) {
	assert.Empty(t, Split(nil, ``))
	assert.Equal(t, []string{`1:2:3`}, Split(nil, `1:2:3`))
	assert.Equal(t, []string{`1`, `2`, `3`}, Split(":", "1:2:3"))
	assert.Equal(t, []string{`1`, `2`, `3`, ``}, Split("\\W+", "1--2::;3    "))
}
