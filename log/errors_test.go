package log

import (
	"errors"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestErrors(t *testing.T) {
	assert := require.New(t)

	e1m := `error 1`
	e1 := errors.New(e1m)

	e2m := `error`
	e2 := errors.New(e2m)

	assert.False(ErrContains(nil, nil))
	assert.False(ErrContains(e1, nil))
	assert.False(ErrContains(nil, e2))
	assert.False(ErrContains(nil, e2m))
	assert.True(ErrContains(e1, e1))
	assert.True(ErrContains(e1, e1m))
	assert.True(ErrContains(e2, e2))
	assert.True(ErrContains(e2, e2m))
	assert.True(ErrContains(e1, e2))
	assert.True(ErrContains(e1, e2m))
	assert.False(ErrContains(e2, e1))
	assert.False(ErrContains(e2, e1m))

	assert.True(ErrHasPrefix(e1, e1))
	assert.True(ErrHasPrefix(e1, e1m))
	assert.True(ErrHasPrefix(e2, e2))
	assert.True(ErrHasPrefix(e2, e2m))
	assert.True(ErrHasPrefix(e1, e2))
	assert.True(ErrHasPrefix(e1, e2m))
	assert.False(ErrHasPrefix(e1, `nope`))
	assert.False(ErrHasPrefix(e2, e1))
	assert.False(ErrHasPrefix(e2, e1m))

	assert.False(ErrHasSuffix(e1, e2))
	assert.False(ErrHasSuffix(e1, e2m))
	assert.True(ErrHasSuffix(e1, `1`))
	assert.False(ErrHasSuffix(e2, e1))
	assert.False(ErrHasSuffix(e2, e1m))
	assert.True(ErrHasSuffix(e1, e1))
	assert.True(ErrHasSuffix(e1, e1m))
	assert.True(ErrHasSuffix(e2, e2))
	assert.True(ErrHasSuffix(e2, e2m))
}
