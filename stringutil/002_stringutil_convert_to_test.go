package stringutil

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConvertToFloat(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Float, nil)
	assert.NoError(err)
	assert.Equal(float64(0), v)

	v, err = ConvertTo(Float, "1.5")
	assert.NoError(err)
	assert.Equal(float64(1.5), v)

	v, err = ConvertTo(Float, "1")
	assert.NoError(err)
	assert.Equal(float64(1.0), v)

	v, err = ConvertToFloat("1.5")
	assert.NoError(err)
	assert.Equal(float64(1.5), v)

	v, err = ConvertToFloat("1.0")
	assert.NoError(err)
	assert.Equal(float64(1.0), v)

	for _, fail := range []string{`potato`, `true`, `2015-05-01 00:15:16`} {
		_, err := ConvertTo(Float, fail)
		assert.Error(err)

		_, err = ConvertToFloat(fail)
		assert.Error(err)
	}
}

func TestConvertToInteger(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Integer, nil)
	assert.NoError(err)
	assert.Equal(int64(0), v)

	v, err = ConvertTo(Integer, "7")
	assert.NoError(err)
	assert.Equal(int64(7), v)

	v, err = ConvertToInteger("7")
	assert.NoError(err)
	assert.Equal(int64(7), v)

	tm := time.Date(2010, 2, 21, 15, 14, 13, 0, time.UTC)

	v, err = ConvertTo(Integer, tm)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	v, err = ConvertToInteger(tm)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	v, err = ConvertTo(Integer, `2010-02-21 15:14:13`)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	for _, fail := range []string{`0.0`, `1.5`, `potato`, `true`} {
		_, err := ConvertTo(Integer, fail)
		assert.Error(err)

		_, err = ConvertToInteger(fail)
		assert.Error(err)
	}
}

func TestConvertToBytes(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Bytes, nil)
	assert.NoError(err)
	assert.Equal([]byte{}, v)

	v, err = ConvertTo(Bytes, []byte{})
	assert.NoError(err)
	assert.Equal([]byte{}, v)

	v, err = ConvertTo(Bytes, []byte{1, 2, 3})
	assert.NoError(err)
	assert.Equal([]byte{1, 2, 3}, v)

	v, err = ConvertTo(Bytes, `test`)
	assert.NoError(err)
	assert.Equal([]byte{0x74, 0x65, 0x73, 0x74}, v)

	v, err = ConvertTo(Bytes, []int{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, string(v.([]byte)))
}

func TestConvertToBoolean(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Boolean, nil)
	assert.Equal(false, v)

	v, err = ConvertTo(Boolean, `true`)
	assert.NoError(err)
	assert.Equal(true, v)

	v, err = ConvertTo(Boolean, `false`)
	assert.NoError(err)
	assert.Equal(false, v)

	v, err = ConvertToBool(`true`)
	assert.NoError(err)
	assert.Equal(true, v)

	v, err = ConvertToBool(`false`)
	assert.NoError(err)
	assert.Equal(false, v)

	for _, fail := range []string{`1.5`, `potato`, `01`, `2015-05-01 00:15:16`} {
		_, err := ConvertTo(Boolean, fail)
		assert.Error(err)

		_, err = ConvertToBool(fail)
		assert.Error(err)
	}
}

func TestConvertToTime(t *testing.T) {
	assert := require.New(t)

	atLeastNow := time.Now()

	values := map[string]time.Time{
		`2015-05-01 00:15:16`:         time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		`Fri May 1 00:15:16 UTC 2015`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `Fri May 01 00:15:16 +0000 2015`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `01 May 15 00:15 UTC`:            time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `01 May 15 00:15 +0000`:          time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `Friday, 01-May-15 00:15:16 UTC`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		`1136239445`: time.Date(2006, 1, 2, 17, 4, 5, 0, time.Now().Location()),
	}

	v, err := ConvertToTime(`now`)
	assert.Nil(err)
	assert.True(v.After(atLeastNow))

	v, err = ConvertToTime(time.Now())
	assert.Nil(err)
	assert.True(v.After(atLeastNow))

	v, err = ConvertToTime(`0000-00-00 00:00:00`)
	assert.Nil(err)
	assert.Zero(v)

	for in, out := range values {
		v, err := ConvertTo(Time, in)
		assert.NoError(err)
		assert.IsType(time.Now(), v)
		assert.Equal(out, v.(time.Time))

		v, err = ConvertToTime(in)
		assert.NoError(err)
		assert.Equal(out, v)
	}

	for _, fail := range []string{`1.5`, `potato`, `false`} {
		_, err := ConvertTo(Time, fail)
		assert.Error(err)

		_, err = ConvertToTime(fail)
		assert.Error(err)
	}
}

func TestAutotypeNil(t *testing.T) {
	assert := require.New(t)

	for _, testValue := range []string{
		``,
		`nil`,
		`null`,
		`Nil`,
		`NULL`,
		`None`,
		`undefined`,
	} {
		assert.Nil(Autotype(testValue), fmt.Sprintf("%q was not autotyped to nil", testValue))
	}
}

func TestAutotypeFloat(t *testing.T) {
	assert := require.New(t)

	for _, testValue := range []string{
		`-0.00000000001`,
		`0.00000000001`,
		`1.5`,
		`-1.5`,
		fmt.Sprintf("%f", math.MaxFloat64),
		fmt.Sprintf("%f", -1*math.MaxFloat64),
	} {
		assert.IsType(float64(0), Autotype(testValue), testValue)
	}
}

func TestAutotypeInt(t *testing.T) {
	assert := require.New(t)

	for _, testValue := range []string{
		`-1`,
		`0`,
		`1`,
		`12345`,
		`-12345`,
		fmt.Sprintf("%d", math.MaxInt64),
		fmt.Sprintf("%d", math.MinInt64),
	} {
		assert.IsType(int64(0), Autotype(testValue))
	}
}

func TestAutotypePreserveLeadingZeroes(t *testing.T) {
	assert := require.New(t)

	for _, testValue := range []string{
		`00`,
		`01`,
		`07753`,
		`06094`,
		`0000000010000000`,
	} {
		assert.IsType(``, Autotype(testValue))
	}
}

func TestAutotypeDate(t *testing.T) {
	assert := require.New(t)

	for _, testValue := range TimeFormats {
		tvS := strings.Replace(string(testValue), `_`, ``, -1)
		tvS = strings.TrimSuffix(tvS, `07:00`)
		assert.IsType(time.Now(), Autotype(tvS))
	}
}

func TestAutotypeBool(t *testing.T) {
	assert := require.New(t)

	for _, testValue := range []string{
		`true`,
		`True`,
		`false`,
		`False`,
	} {
		assert.IsType(true, Autotype(testValue))
	}

	for _, testValue := range []string{
		`trues`,
		`Falses`,
		`potato`,
	} {
		assert.IsType(``, Autotype(testValue))
	}
}
