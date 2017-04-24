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
	assert.True(IsEmpty(``))
	assert.True(IsEmpty(` `))
	assert.True(IsEmpty(`     `))
	assert.True(IsEmpty("\t\n  \n\t"))
	assert.True(IsEmpty(thingmap))
	assert.True(IsEmpty(stringmap))
}

func TestRelaxedEqual(t *testing.T) {
	assert := require.New(t)

	eq, err := RelaxedEqual(nil, nil)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(1, 1)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(int(1), int64(1))
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(float64(1), byte(1))
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(float64(1.00), `1`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(true, true)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(false, false)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(`true`, `on`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(`true`, `yes`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(`boo`, `boo`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(1, true)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(true, false)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(false, true)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(`true`, `no`)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(`false`, `yes`)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(`boo`, `Boo`)
	assert.NoError(err)
	assert.False(eq)
}

func TestIsArray(t *testing.T) {
	assert := require.New(t)

	assert.False(IsArray(nil))

	var a []string
	assert.True(IsArray(a))
	assert.True(IsArray([]string{`1`}))
	assert.True(IsArray(&[]string{`1`}))

	var b interface{}
	b = []string{`1`}

	assert.True(IsArray(b))
	assert.True(IsArray(&b))

	assert.False(IsArray(``))
	assert.False(IsArray(`123`))
	assert.False(IsArray(123))
	assert.False(IsArray(true))
}
