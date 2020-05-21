package maputil

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ghetzel/testify/require"
)

type MyTestThing struct {
	Name  string
	Other int
}

func TestMapJoin(t *testing.T) {
	input := map[string]interface{}{
		`key1`: `value1`,
		`key2`: true,
		`key3`: 3,
	}

	output := Join(input, `=`, `&`)

	if output == `` {
		t.Error("Output should not be empty")
	}

	if !strings.Contains(output, `key1=value1`) {
		t.Errorf("Output should contain '%s'", `key1=value1`)
	}

	if !strings.Contains(output, `key2=true`) {
		t.Errorf("Output should contain '%s'", `key2=true`)
	}

	if !strings.Contains(output, `key3=3`) {
		t.Errorf("Output should contain '%s'", `key3=3`)
	}
}

func TestStringKeys(t *testing.T) {
	assert := require.New(t)

	i1 := map[string]interface{}{
		`1`: 1,
		`2`: true,
		`3`: `three`,
	}

	i2 := map[string]bool{
		`1`: true,
		`2`: false,
		`3`: true,
	}

	i3 := map[string]MyTestThing{
		`1`: MyTestThing{},
		`2`: MyTestThing{},
		`3`: MyTestThing{},
	}

	var i4 sync.Map

	i4.Store(`1`, MyTestThing{})
	i4.Store(`2`, 2)
	i4.Store(`3`, 3.14)

	output := []string{`1`, `2`, `3`}

	assert.Empty(StringKeys(nil))

	assert.Equal(output, StringKeys(i1))
	assert.Equal(output, StringKeys(i2))
	assert.Equal(output, StringKeys(i3))
	assert.Equal(output, StringKeys(&i4))

	assert.Empty(StringKeys(true))
	assert.Empty(StringKeys(4))
	assert.Empty(StringKeys([]int{1, 2, 3}))
}

func TestMapSplit(t *testing.T) {
	input := `key1=value1&key2=true&key3=3`

	output := Split(input, `=`, `&`)

	if len(output) == 0 {
		t.Error("Output should not be empty")
	}

	if v, ok := output[`key1`]; !ok || v != `value1` {
		t.Errorf("Output should contain key %s => '%s'", `key1`, `value1`)
	}

	if v, ok := output[`key2`]; !ok || v != `true` {
		t.Errorf("Output should contain key %s => '%s'", `key2`, `true`)
	}

	if v, ok := output[`key3`]; !ok || v != `3` {
		t.Errorf("Output should contain key %s => '%s'", `key3`, `3`)
	}
}

type SubtypeTester struct {
	A int
	B int `maputil:"b"`
}

type MyStructTester struct {
	Name                  string
	Subtype1              SubtypeTester
	Active                bool           `maputil:"active"`
	Subtype2              *SubtypeTester `maputil:"subtype2"`
	TimeTest              time.Duration
	IntTest               int32
	Properties            map[string]interface{}
	StrSliceTest          []string
	InterfaceStrSliceTest []string
	StructSliceTest       []SubtypeTester
	StructSliceTest2      []SubtypeTester
	StructSliceTest3      []SubtypeTester
	nonexported           int
}

func TestStructFromMapEmbedded(t *testing.T) {
	type tPerson struct {
		Name string
		Age  int `potato:"age"`
	}

	type tUser struct {
		tPerson
		Email  string `potato:"email"`
		Active bool   `potato:"ACTIVE"`
	}

	assert := require.New(t)

	var tgt tUser

	assert.NoError(TaggedStructFromMap(map[string]interface{}{
		`Name`:   `Rusty Shackleford`,
		`age`:    420,
		`email`:  `none+of@your.biz`,
		`ACTIVE`: true,
	}, &tgt, `potato`))

	assert.Equal(`Rusty Shackleford`, tgt.Name)
	assert.Equal(420, tgt.Age)
	assert.Equal(`none+of@your.biz`, tgt.Email)
	assert.True(tgt.Active)
}

func TestStructFromMap(t *testing.T) {
	assert := require.New(t)

	input := map[string]interface{}{
		`Name`:           `Foo Bar`,
		`active`:         true,
		`should_not_set`: 4,
		`Subtype1`: map[string]interface{}{
			`A`: 1,
			`b`: 2,
		},
		`subtype2`: map[string]interface{}{
			`A`: 3,
			`b`: 4,
		},
		`TimeTest`: 15000000000,
		`IntTest`:  int64(5),
		`Properties`: map[string]interface{}{
			`first`:  1,
			`second`: true,
			`third`:  `three`,
		},
		`StrSliceTest`:          []string{`one`, `two`, `three`},
		`InterfaceStrSliceTest`: []interface{}{`one`, `two`, `three`},
		`StructSliceTest`:       []SubtypeTester{{10, 11}, {12, 13}, {14, 15}},
		`StructSliceTest2`: []map[string]interface{}{
			{
				`A`: 10,
				`b`: 11,
			},
			{
				`A`: 12,
				`b`: 13,
			},
			{
				`A`: 14,
				`b`: 15,
			},
		},
		`StructSliceTest3`: []interface{}{
			map[string]interface{}{
				`A`: 10,
				`b`: 11,
			},
			map[string]interface{}{
				`A`: 12,
				`b`: 13,
			},
			map[string]interface{}{
				`A`: 14,
				`b`: 15,
			},
		},
	}

	output := MyStructTester{}

	err := StructFromMap(input, &output)
	assert.NoError(err)

	assert.Equal(`Foo Bar`, output.Name)
	assert.True(output.Active)
	assert.Zero(output.nonexported)

	// assert.Equal(1, output.Subtype1.A)
	// assert.Equal(2, output.Subtype1.B)

	assert.NotNil(output.Subtype2)
	assert.Equal(3, output.Subtype2.A)
	assert.Equal(4, output.Subtype2.B)

	assert.Equal(time.Duration(15)*time.Second, output.TimeTest)
	assert.Equal(int32(5), output.IntTest)

	assert.NotNil(output.Properties)
	assert.EqualValues(1, output.Properties[`first`])
	assert.EqualValues(true, output.Properties[`second`])
	assert.Equal(`three`, output.Properties[`third`])

	assert.NotNil(output.StrSliceTest)
	assert.Len(output.StrSliceTest, 3)
	assert.Equal(`one`, output.StrSliceTest[0])
	assert.Equal(`two`, output.StrSliceTest[1])
	assert.Equal(`three`, output.StrSliceTest[2])

	assert.NotNil(output.InterfaceStrSliceTest)
	assert.Len(output.InterfaceStrSliceTest, 3)
	assert.EqualValues(`one`, output.InterfaceStrSliceTest[0])
	assert.EqualValues(`two`, output.InterfaceStrSliceTest[1])
	assert.EqualValues(`three`, output.InterfaceStrSliceTest[2])

	assert.NotNil(output.StructSliceTest)
	assert.Len(output.StructSliceTest, 3)
	assert.EqualValues(10, output.StructSliceTest[0].A)
	assert.EqualValues(11, output.StructSliceTest[0].B)

	assert.EqualValues(12, output.StructSliceTest[1].A)
	assert.EqualValues(13, output.StructSliceTest[1].B)

	assert.EqualValues(14, output.StructSliceTest[2].A)
	assert.EqualValues(15, output.StructSliceTest[2].B)

	assert.NotNil(output.StructSliceTest2)
	assert.Len(output.StructSliceTest2, 3)
	assert.EqualValues(10, output.StructSliceTest2[0].A)
	assert.EqualValues(11, output.StructSliceTest2[0].B)

	assert.EqualValues(12, output.StructSliceTest2[1].A)
	assert.EqualValues(13, output.StructSliceTest2[1].B)

	assert.EqualValues(14, output.StructSliceTest2[2].A)
	assert.EqualValues(15, output.StructSliceTest2[2].B)

	assert.NotNil(output.StructSliceTest3)
	assert.Len(output.StructSliceTest3, 3)
	assert.EqualValues(10, output.StructSliceTest3[0].A)
	assert.EqualValues(11, output.StructSliceTest3[0].B)

	assert.EqualValues(12, output.StructSliceTest3[1].A)
	assert.EqualValues(13, output.StructSliceTest3[1].B)

	assert.EqualValues(14, output.StructSliceTest3[2].A)
	assert.EqualValues(15, output.StructSliceTest3[2].B)
}

func TestMapAppend(t *testing.T) {
	assert := require.New(t)

	assert.Equal(map[string]interface{}{}, Append())

	assert.Equal(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, Append(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}))

	assert.Equal(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, Append(nil, map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, nil))

	assert.Equal(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
		`d`: 4,
		`e`: false,
		`f`: 6.1,
	}, Append(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, map[string]interface{}{
		`d`: 4,
		`e`: false,
		`f`: 6.1,
	}))

	assert.Equal(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Five`,
	}, Append(map[string]interface{}{
		`a`: 1,
		`b`: true,
		`c`: `Three`,
	}, map[string]interface{}{
		`c`: `Five`,
	}))
}

// func TestMapValues(t *testing.T) {
// 	assert := require.New(t)

// 	assert.Equal([]interface{}{
// 		1, 3, 5,
// 	}, MapValues(map[string]int{
// 		`first`:  1,
// 		`second`: 3,
// 		`third`:  5,
// 	}))
// }

func TestApply(t *testing.T) {
	assert := require.New(t)

	assert.Equal(map[string]interface{}{
		`a`: 10,
		`b`: 20,
		`c`: 30,
	}, Apply(map[string]interface{}{
		`a`: 1,
		`b`: 2,
		`c`: 3,
	}, func(_ []string, value interface{}) (interface{}, bool) {
		return value.(int) * 10, true
	}))
}
