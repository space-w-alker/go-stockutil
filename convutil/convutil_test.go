package convutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	assert := assert.New(t)

	// temperatures
	assert.EqualValues(0, MustConvert(32, Fahrenheit, Celsius))
	assert.EqualValues(0, MustConvert(-459.67, Fahrenheit, Kelvin))
	assert.EqualValues(32, MustConvert(0, Celsius, Fahrenheit))
	assert.EqualValues(273.15, MustConvert(0, Celsius, Kelvin))
	assert.EqualValues(-273.15, MustConvert(0, Kelvin, Celsius))
	assert.EqualValues(-459.67, MustConvert(0, Kelvin, Fahrenheit))

	// lengths
	assert.EqualValues(1609.344, MustConvert(1, Miles, Meters))
	assert.EqualValues(1, MustConvert(1609.344, Meters, Miles))
	assert.EqualValues(4200, MustConvert(2.609759007397, Miles, Meters))
	assert.EqualValues(2.609759000000, MustConvert(4200, Meters, Miles))
	assert.EqualValues(1, MustConvert(9.4607304725808e+15, Meters, Lightyears))
	assert.EqualValues(9.4607304725808e+15, MustConvert(1, Lightyears, Meters))
}
