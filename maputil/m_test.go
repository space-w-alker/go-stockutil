package maputil

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/testify/require"
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

func TestMMarshalXML(t *testing.T) {
	assert := require.New(t)

	m := M(map[string]interface{}{
		`hello`: 1,
		`there`: true,
		`general`: map[string]interface{}{
			`kenobi`: true,
		},
		`xyz`: []string{`a`, `b`, `c`},
		`zzz`: []map[string]interface{}{
			map[string]interface{}{
				`name`:  `a`,
				`value`: 0,
			},
			map[string]interface{}{
				`name`:  `b`,
				`value`: 1,
			},
			map[string]interface{}{
				`name`:  `c`,
				`value`: 2,
			},
		},
	})

	// default marshal
	out, err := xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<data><general><kenobi>true</kenobi></general><hello>1</hello><there>true</there><xyz><element>a</element><element>b</element><element>c</element></xyz><zzz><element><name>a</name><value>0</value></element><element><name>b</name><value>1</value></element><element><name>c</name><value>2</value></element></zzz></data>`), out)

	// custom root tagname
	m.SetRootTagName(`nub_nub`)
	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<nub_nub><general><kenobi>true</kenobi></general><hello>1</hello><there>true</there><xyz><element>a</element><element>b</element><element>c</element></xyz><zzz><element><name>a</name><value>0</value></element><element><name>b</name><value>1</value></element><element><name>c</name><value>2</value></element></zzz></nub_nub>`), out)

	// generic XML structure
	m.SetMarshalXmlGeneric(true)
	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<nub_nub><item type="object" key="general"><item key="kenobi" type="bool">true</item></item><item key="hello" type="int">1</item><item key="there" type="bool">true</item><item type="array" key="xyz"><item key="element" type="str">a</item><item key="element" type="str">b</item><item key="element" type="str">c</item></item><item type="array" key="zzz"><item type="object" key="element"><item key="name" type="str">a</item><item key="value" type="int">0</item></item><item type="object" key="element"><item key="name" type="str">b</item><item key="value" type="int">1</item></item><item type="object" key="element"><item key="name" type="str">c</item><item key="value" type="int">2</item></item></item></nub_nub>`), out)

	// compact structure, custom keyfunc
	m.SetMarshalXmlGeneric(false)
	m.SetMarshalXmlKeyFunc(func(in string) string {
		return stringutil.Camelize(in)
	})

	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<NubNub><General><Kenobi>true</Kenobi></General><Hello>1</Hello><There>true</There><Xyz><Element>a</Element><Element>b</Element><Element>c</Element></Xyz><Zzz><Element><Name>a</Name><Value>0</Value></Element><Element><Name>b</Name><Value>1</Value></Element><Element><Name>c</Name><Value>2</Value></Element></Zzz></NubNub>`), out)

	// generic structure, custom keyfunc
	m.SetMarshalXmlGeneric(true)
	out, err = xml.Marshal(m)
	assert.NoError(err)
	assert.Equal([]byte(`<NubNub><item type="object" key="General"><item key="Kenobi" type="bool">true</item></item><item key="Hello" type="int">1</item><item key="There" type="bool">true</item><item type="array" key="Xyz"><item key="Element" type="str">a</item><item key="Element" type="str">b</item><item key="Element" type="str">c</item></item><item type="array" key="Zzz"><item type="object" key="Element"><item key="Name" type="str">a</item><item key="Value" type="int">0</item></item><item type="object" key="Element"><item key="Name" type="str">b</item><item key="Value" type="int">1</item></item><item type="object" key="Element"><item key="Name" type="str">c</item><item key="Value" type="int">2</item></item></item></NubNub>`), out)
}
