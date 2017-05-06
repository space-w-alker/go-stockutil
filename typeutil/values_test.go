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

func TestIsFunction(t *testing.T) {
	assert := require.New(t)

	assert.False(IsFunction(nil))
	assert.False(IsFunction(1))
	assert.False(IsFunction(true))
	assert.False(IsFunction(`three`))

	assert.True(IsFunction(func() {}))
	assert.True(IsFunctionArity(func() {}, 0, 0))
	assert.True(IsFunctionArity(func() {}, 0, -1))
	assert.True(IsFunctionArity(func() {}, -1, 0))
	assert.True(IsFunctionArity(func() {}, -1, -1))
	assert.False(IsFunctionArity(func() {}, 99, 0))
	assert.False(IsFunctionArity(func() {}, 0, 99))
	assert.False(IsFunctionArity(func() {}, 99, 99))

	assert.True(IsFunction(func(interface{}) {}))
	assert.True(IsFunctionArity(func(interface{}) {}, 1, 0))
	assert.True(IsFunctionArity(func(interface{}) {}, 1, -1))
	assert.True(IsFunctionArity(func(interface{}) {}, -1, 0))
	assert.True(IsFunctionArity(func(interface{}) {}, -1, -1))
	assert.False(IsFunctionArity(func(interface{}) {}, 99, 0))
	assert.False(IsFunctionArity(func(interface{}) {}, 0, 99))
	assert.False(IsFunctionArity(func(interface{}) {}, 99, 99))

	assert.True(IsFunction(func(interface{}) error { return nil }))
	assert.True(IsFunctionArity(func(interface{}) error { return nil }, 1, 1))
	assert.True(IsFunctionArity(func(interface{}) error { return nil }, 1, -1))
	assert.True(IsFunctionArity(func(interface{}) error { return nil }, -1, 1))
	assert.True(IsFunctionArity(func(interface{}) error { return nil }, -1, -1))
	assert.False(IsFunctionArity(func(interface{}) error { return nil }, 99, 1))
	assert.False(IsFunctionArity(func(interface{}) error { return nil }, 0, 99))
	assert.False(IsFunctionArity(func(interface{}) error { return nil }, 99, 99))

	assert.True(IsFunction(func() error { return nil }))
	assert.True(IsFunctionArity(func() error { return nil }, 0, 1))
	assert.True(IsFunctionArity(func() error { return nil }, 0, -1))
	assert.True(IsFunctionArity(func() error { return nil }, -1, 1))
	assert.True(IsFunctionArity(func() error { return nil }, -1, -1))
	assert.False(IsFunctionArity(func() error { return nil }, 99, 0))
	assert.False(IsFunctionArity(func() error { return nil }, 0, 99))
	assert.False(IsFunctionArity(func() error { return nil }, 99, 99))
}
