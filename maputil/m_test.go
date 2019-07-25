package maputil

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/stretchr/testify/require"
)

type testMstruct struct {
	ID     string `maputil:"id" json:"_id"`
	Name   string `json:"NAME"`
	Factor float64
}

func TestM(t *testing.T) {
	assert := require.New(t)
	input := M(map[string]interface{}{
		`first`: true,
		`second`: map[string]interface{}{
			`s1`:     `test`,
			`values`: []int{1, 2, 3, 4},
			`truthy`: `True`,
			`strnum`: `42`,
			`then`:   `2006-01-02`,
		},
		`now`:    time.Now(),
		`third`:  3.1415,
		`fourth`: 42,
	})

	assert.Equal(``, M(nil).String(`lol`))
	assert.False(M(nil).Bool(`lol`))
	assert.Equal(int64(0), M(nil).Int(`lol`))
	assert.Equal(float64(0), M(nil).Float(`lol`))
	assert.Len(M(nil).Slice(`lol`), 0)
	assert.Nil(M(nil).Auto(`second.strnum`))
	assert.Zero(M(nil).Time(`now`))

	assert.Equal(`test`, input.String(`second.s1`))
	assert.True(input.Bool(`first`))
	assert.True(input.Bool(`second.truthy`))
	assert.True(input.Bool(`second.s1`))
	assert.Equal(3.1415, input.Float(`third`))
	assert.Equal(int64(3), input.Int(`third`))
	assert.Equal(int64(42), input.Int(`fourth`))
	assert.Equal(int64(3), input.Int(`second.values.2`))
	assert.Equal(int64(0), input.Int(`second.values.99`))
	assert.Equal(float64(42), input.Float(`fourth`))
	assert.Len(input.Slice(`second.values`), 4)
	assert.Equal(int64(42), input.Auto(`second.strnum`))
	assert.Equal(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), input.Time(`second.then`))

	assert.Equal(5, input.Len())
	k := make([]string, 5)
	i := 0

	assert.NoError(input.Each(func(key string, value typeutil.Variant) error {
		k[i] = key
		i++
		return nil
	}))

	assert.ElementsMatch(k, []string{`first`, `second`, `third`, `fourth`, `now`})
}

func TestMSet(t *testing.T) {
	assert := require.New(t)
	input := M(nil)

	assert.Equal(``, input.String(`lol`))

	assert.Equal(`2funny4me`, input.Set(`lol`, `2funny4me`).String())

	assert.Equal(`2funny4me`, input.String(`lol`))
}

func TestMStruct(t *testing.T) {
	assert := require.New(t)
	input := M(&testMstruct{
		ID:     `123`,
		Name:   `tester`,
		Factor: 3.14,
	})

	assert.Equal(`123`, input.String(`id`))
	assert.EqualValues(123, input.Int(`id`))

	assert.Equal(`tester`, input.String(`Name`))
	assert.Equal(3.14, input.Float(`Factor`))

	assert.Equal(map[string]interface{}{
		`id`:     `123`,
		`Name`:   `tester`,
		`Factor`: 3.14,
	}, input.MapNative())

	assert.Equal(map[string]interface{}{
		`_id`:    `123`,
		`NAME`:   `tester`,
		`Factor`: 3.14,
	}, input.MapNative(`json`))
}

func TestMUrlValues(t *testing.T) {
	assert := require.New(t)
	input := M(url.Values{
		`a`: []string{`1`},
		`b`: []string{},
		`c`: []string{`2`, `3`},
	})

	assert.Equal(`1`, input.String(`a`))
	assert.EqualValues(1, input.Int(`a`))

	assert.Equal(``, input.String(`b`))
	assert.Equal(float64(0), input.Float(`b`))
	assert.Nil(input.Auto(`b`))

	assert.Equal([]string{`2`, `3`}, input.Strings(`c`))
}

func TestMHttpHeader(t *testing.T) {
	assert := require.New(t)
	input := M(http.Header{
		`a`: []string{`1`},
		`b`: []string{},
		`c`: []string{`2`, `3`},
	})

	assert.Equal(`1`, input.String(`a`))
	assert.EqualValues(1, input.Int(`a`))

	assert.Equal(``, input.String(`b`))
	assert.Equal(float64(0), input.Float(`b`))
	assert.Nil(input.Auto(`b`))

	assert.Equal([]string{`2`, `3`}, input.Strings(`c`))
}

func TestMStructNested(t *testing.T) {
	assert := require.New(t)

	type msecond struct {
		S1     string
		Values []int
		Truthy interface{}
		Strnum string
		Then   string
	}

	type mtop struct {
		First  bool
		Second msecond
		Now    time.Time
		Third  float64
		Fourth int
	}

	input := M(mtop{
		First: true,
		Second: msecond{
			S1:     `test`,
			Values: []int{1, 2, 3, 4},
			Truthy: `True`,
			Strnum: `42`,
			Then:   `2006-01-02`,
		},
		Now:    time.Now(),
		Third:  3.1415,
		Fourth: 42,
	})

	assert.Equal(`test`, input.String(`Second.S1`))
	assert.True(input.Bool(`First`))
	assert.True(input.Bool(`Second.Truthy`))
	assert.True(input.Bool(`Second.S1`))
	assert.Equal(3.1415, input.Float(`Third`))
	assert.Equal(int64(3), input.Int(`Third`))
	assert.Equal(int64(42), input.Int(`Fourth`))
	assert.Equal(int64(3), input.Int(`Second.Values.2`))
	assert.Equal(int64(0), input.Int(`Second.Values.99`))
	assert.Equal(float64(42), input.Float(`Fourth`))
	assert.Len(input.Slice(`Second.Values`), 4)
	assert.Equal(int64(42), input.Auto(`Second.Strnum`))
	assert.Equal(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), input.Time(`Second.Then`))
}
