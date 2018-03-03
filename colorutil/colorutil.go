package colorutil

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	colorful "github.com/lucasb-eyer/go-colorful"
)

var rgbaPattern = `rgba?\((?P<r>[\d\.]+%?)\s*,\s*(?P<g>[\d\.]+%?)\s*,\s*(?P<b>[\d\.]+%?)\s*(?:,\s*(?P<a>[\d\.]+%?)?\s*)?\)`
var hslaPattern = `hsla?\((?P<h>[\d\.]+%?)(?:deg)?\s*,\s*(?P<s>[\d\.]+%?)\s*,\s*(?P<l>[\d\.]+%?)\s*(?:,\s*(?P<a>[\d\.]+%?)?\s*)?\)`

type Color struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func (self Color) String() string {
	if self.a < 0xFF {
		return fmt.Sprintf("#%02X%02X%02X%02X", self.r, self.b, self.g, self.a)
	} else {
		return fmt.Sprintf("#%02X%02X%02X", self.r, self.b, self.g)
	}
}

func (self Color) StringRGBA() string {
	if self.a < 0xFF {
		return fmt.Sprintf(
			"rgba(%d, %d, %d, %0.2f)",
			self.r,
			self.g,
			self.b,
			(self.a / 100.0),
		)
	}else{
		return fmt.Sprintf(
			"rgb(%d, %d, %d)",
			self.r,
			self.g,
			self.b,
		)
	}
}

// func (self Color) StringHSLA() string {
// }

func (self Color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(self.r),
		uint32(self.g),
		uint32(self.b),
		uint32(self.a)
}

func (self Color) Equals(other interface{}) bool {
	if color, err := Parse(other); err == nil {
		if self.r == color.r {
			if self.g == color.g {
				if self.b == color.b {
					if self.a == color.a {
						return true
					}
				}
			}
		}
	}

	return false
}

func Parse(colorI interface{}) (Color, error) {
	var colorC colorful.Color
	var alpha float64 = 1.0

	if c, ok := colorI.(colorful.Color); ok {
		colorC = c

	} else if c, ok := colorI.(Color); ok {
		return c, nil
	} else {
		colorS := fmt.Sprintf("%v", colorI)
		colorS = strings.TrimSpace(colorS)

		if rgba := rxutil.Match(rgbaPattern, colorS); rgba != nil {
			// handles rgb(), rgba() patterns
			for v, factorS := range rgba.NamedCaptures() {
				var denom float64 = 0xFF

				if strings.HasSuffix(factorS, `%`) {
					denom = 100
					factorS = factorS[:len(factorS)-1]
				}

				if factorF, err := stringutil.ConvertToFloat(factorS); err == nil {
					if v != `a` && factorF > 1 {
						factorF = (factorF / denom)
					}

					switch v {
					case `r`:
						colorC.R = factorF
					case `g`:
						colorC.G = factorF
					case `b`:
						colorC.B = factorF
					case `a`:
						if factorS != `` {
							if factorF > 1 {
								factorF = (factorF / denom)
							}

							alpha = factorF
						}
					}
				} else {
					return Color{}, fmt.Errorf("Invalid rgba() value in '%v' position: %v", v, err)
				}
			}

		} else if hsla := rxutil.Match(hslaPattern, colorS); hsla != nil {
			// handles hsl(), hsla() patterns
			for v, factorS := range hsla.NamedCaptures() {
				var denom float64 = 0xFF

				if strings.HasSuffix(factorS, `%`) {
					denom = 100
					factorS = factorS[:len(factorS)-1]
				}

				if factorF, err := stringutil.ConvertToFloat(factorS); err == nil {
					var h, s, l float64

					switch v {
					case `h`:
						h = (factorF / 360)
					case `s`:
						s = (factorF / denom)
					case `l`:
						l = (factorF / denom)
					case `a`:
						if factorS != `` {
							if factorF > 1 {
								factorF = (factorF / denom)
							}

							alpha = factorF
						}
					}

					colorC = colorful.Hsl(h, s, l)
				} else {
					return Color{}, fmt.Errorf("Invalid rgba() value in '%v' position: %v", v, err)
				}
			}
		} else {
			colorS = `#` + strings.ToUpper(strings.TrimPrefix(colorS, `#`))

			if len(colorS) == 9 {
				if v, err := hex.DecodeString(colorS[7:]); err == nil {
					if len(v) == 1 {
						alpha = float64(v[0]) / 0xFF
					}
				} else {
					return Color{}, fmt.Errorf("Invalid hex value '%v'", colorS[7:])
				}
			}

			if c, err := colorful.Hex(colorS); err == nil {
				colorC = c
			} else {
				return Color{}, fmt.Errorf("Invalid color '%v': %v", colorS, err)
			}
		}
	}

	returnColor := Color{
		r: uint8(0xFF * colorC.R),
		g: uint8(0xFF * colorC.G),
		b: uint8(0xFF * colorC.B),
		a: uint8(0xFF * alpha),
	}

	return returnColor, nil
}

func Equals(first interface{}, second interface{}) bool {
	if first != nil && second != nil {
		if firstC, err := Parse(first); err == nil {
			return firstC.Equals(second)
		}
	}

	return false
}

func Darken(in interface{}, percent float64) (Color, error) {
	var isLighten bool

	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		isLighten = true
		percent = (percent - 100)
	}

	factor := (percent / 100)

	if sample, err := Parse(in); err == nil {
		h, s, l := colorful.MakeColor(sample).Hsl()

		if isLighten {
			// find the delta between full luminosity and current luminosity, then use the
			// factor to take us that percentage of the way there
			l = l + ((1.0 - l) * factor)
		} else {
			// since we're darkening, we can just apply the factor as-is, because the implicit
			// lower bound is already zero
			l = l * factor
		}

		r, g, b := colorful.Hsl(h, s, l).Clamped().RGB255()

		sample.r = r
		sample.g = g
		sample.b = b

		return sample, nil
	} else {
		return Color{}, err
	}
}

func Lighten(in interface{}, percent float64) (Color, error) {
	return Darken(in, 100+percent)
}
