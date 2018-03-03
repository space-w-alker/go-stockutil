package colorutil

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func assertColor(t *testing.T, in string, rIn uint32, gIn uint32, bIn uint32, aIn uint32) {
	assert := require.New(t)

	var color color.Color
	var err error
	var r, g, b, a uint32

	color, err = Parse(in)
	assert.NoError(err)

	r, g, b, a = color.RGBA()
	assert.Equal(rIn, r, fmt.Sprintf("%v: red", in))
	assert.Equal(gIn, g, fmt.Sprintf("%v: green", in))
	assert.Equal(bIn, b, fmt.Sprintf("%v: blue", in))
	assert.Equal(aIn, a, fmt.Sprintf("%v: alpha", in))
}

func TestParse(t *testing.T) {
	assertColor(t, `#FFFFFF`, 255, 255, 255, 255)
	assertColor(t, `#FFFF00`, 255, 255, 0, 255)
	assertColor(t, `#FF00CC`, 255, 0, 204, 255)
	assertColor(t, `#FFFFFF00`, 255, 255, 255, 0)
	assertColor(t, `#FFFF00FF`, 255, 255, 0, 255)
	assertColor(t, `#FF00FFFF`, 255, 0, 255, 255)
	assertColor(t, `#00FFFFFF`, 0, 255, 255, 255)
	assertColor(t, `#FF000000`, 255, 0, 0, 0)
	assertColor(t, `#FF00FF00`, 255, 0, 255, 0)
	assertColor(t, `#00000000`, 0, 0, 0, 0)

	assertColor(t, `rgb(255, 255, 255)`, 255, 255, 255, 255)
	assertColor(t, `rgb(255, 255, 42)`, 255, 255, 42, 255)
	assertColor(t, `rgb(255, 42, 255)`, 255, 42, 255, 255)
	assertColor(t, `rgb(42, 255, 255)`, 42, 255, 255, 255)
	assertColor(t, `rgba(255, 255, 255, 255)`, 255, 255, 255, 255)
	assertColor(t, `rgba(255, 255, 255, 42)`, 255, 255, 255, 42)
	assertColor(t, `rgba(255, 255, 42, 255)`, 255, 255, 42, 255)
	assertColor(t, `rgba(255, 42, 255, 255)`, 255, 42, 255, 255)
	assertColor(t, `rgba(42, 255, 255, 255)`, 42, 255, 255, 255)
}

func TestEquals(t *testing.T) {
	assert := require.New(t)

	assert.True(Equals(`#FFF`, `#FFF`))
	assert.True(Equals(`#FFF`, `#FFFFFF`))
	assert.True(Equals(`#FFF`, `rgb(255,255,255)`))
	assert.True(Equals(`#00AA00`, `rgb(0,170,0)`))
}

func TestLightenDarken(t *testing.T) {
	assert := require.New(t)

	var color Color
	var err error

	// Darken(#800000, 25%) => #200000
	color, err = Darken(`#800000`, 25)
	assert.NoError(err)
	assert.True(color.Equals(`#200000`), fmt.Sprintf("%v", color))

	// Lighten(#800000, 25%) => #E00000
	color, err = Lighten(`#800000`, 25)
	assert.NoError(err)
	assert.True(color.Equals(`#E00000`), fmt.Sprintf("%v", color))
}
