package sliceutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainsString(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}

	assert.True(ContainsString(input, "one"))
	assert.True(ContainsString(input, "three"))
	assert.True(ContainsString(input, "five"))
	assert.False(ContainsString(input, "One"))
	assert.False(ContainsString(input, "two"))
	assert.False(ContainsString(input, "Three"))
	assert.False(ContainsString(input, "four"))
	assert.False(ContainsString(input, "Five"))
	assert.False(ContainsString([]string{}, "one"))
	assert.False(ContainsString([]string{}, "two"))
	assert.False(ContainsString([]string{}, ""))
}

func TestCompact(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Compact(nil))

	assert.Equal([]interface{}{
		0, 1, 2, 3,
	}, Compact([]interface{}{
		0, 1, 2, 3,
	}))

	assert.Equal([]interface{}{
		1, 3, 5,
	}, Compact([]interface{}{
		nil, 1, nil, 3, nil, 5,
	}))

	assert.Equal([]interface{}{
		`one`, `three`, ` `, `five`,
	}, Compact([]interface{}{
		`one`, ``, `three`, ``, ` `, `five`,
	}))

	assert.Equal([]interface{}{
		[]int{1, 2, 3},
	}, Compact([]interface{}{
		nil, []string{}, []int{1, 2, 3}, map[string]bool{},
	}))
}

func TestCompactString(t *testing.T) {
	assert := require.New(t)

	assert.Nil(CompactString(nil))

	assert.Equal([]string{
		`one`, `three`, `five`,
	}, CompactString([]string{
		`one`, `three`, `five`,
	}))

	assert.Equal([]string{
		`one`, `three`, ` `, `five`,
	}, CompactString([]string{
		`one`, ``, `three`, ``, ` `, `five`,
	}))
}

func TestStringify(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Stringify(nil))

	assert.Equal([]string{
		`0`, `1`, `2`,
	}, Stringify([]interface{}{
		0, 1, 2,
	}))

	assert.Equal([]string{
		`0.5`, `0.55`, `0.555`, `0.555001`,
	}, Stringify([]interface{}{
		0.5, 0.55, 0.55500, 0.555001,
	}))

	assert.Equal([]string{
		`true`, `true`, `false`,
	}, Stringify([]interface{}{
		true, true, false,
	}))
}
