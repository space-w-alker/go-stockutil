package structutil

import (
	"reflect"
	"testing"

	"github.com/ghetzel/testify/require"
)

type TBase struct {
	Name    string
	Age     int
	Enabled bool
	hidden  bool
}

type tPerson struct {
	*TBase
	Address string
}

func TestFieldsFunc(t *testing.T) {
	assert := require.New(t)

	out := make(map[string]interface{})

	FieldsFunc(&TBase{
		Name:    `hello`,
		Age:     32,
		Enabled: true,
		hidden:  true,
	}, func(field *reflect.StructField, value reflect.Value) error {
		if value.CanInterface() {
			out[field.Name] = value.Interface()
		}

		return nil
	})

	assert.Equal(map[string]interface{}{
		`Name`:    `hello`,
		`Age`:     32,
		`Enabled`: true,
	}, out)
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
