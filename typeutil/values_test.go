package typeutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Thing struct {
	Name  string
	Value interface{}
}

func TestIsZero(t *testing.T) {
	assert := require.New(t)

	var thing Thing
	var things []Thing
	madeThings := make([]Thing, 0)

	assert.True(IsZero(nil))
	assert.True(IsZero(0))
	assert.True(IsZero(0.0))
	assert.True(IsZero(false))
	assert.True(IsZero(``))
	assert.True(IsZero(thing))
	assert.True(IsZero(Thing{}))
	assert.True(IsZero(things))

	things = append(things, Thing{})

	assert.False(IsZero(1))
	assert.False(IsZero(0.1))
	assert.False(IsZero(true))
	assert.False(IsZero(`value`))
	assert.False(IsZero(Thing{`value`, true}))
	assert.False(IsZero(&Thing{}))
	assert.False(IsZero(things))
	assert.False(IsZero(madeThings))
}

func TestIsEmpty(t *testing.T) {
	assert := require.New(t)

	things := make([]Thing, 4)
	thingmap := make(map[string]Thing)
	stringmap := map[int]string{
		1: ``,
		2: `    `,
		3: "\t",
	}

	assert.True(IsEmpty(things))
	assert.True(IsEmpty(` `))
	assert.True(IsEmpty(`     `))
	assert.True(IsEmpty("\t\n  \n\t"))
	assert.True(IsEmpty(thingmap))
	assert.True(IsEmpty(stringmap))
}
