package convutil

import (
	"fmt"
	"strings"
	"unicode"
)

type UnitFamily int

const (
	GeneralUnits UnitFamily = iota
	TemperatureUnits
	LengthUnits
	SpeedUnits
)

type Unit int

const (
	Invalid Unit = iota
	AU
	Celsius
	Fahrenheit
	Feet
	Kelvin
	KilometersPerHour
	Lightseconds
	Lightminutes
	Lightyears
	Meters
	MetersPerSecond
	Miles
	MilesPerHour
	NauticalMiles
)

func (self Unit) Family() UnitFamily {
	switch self {
	case Celsius, Fahrenheit, Kelvin:
		return TemperatureUnits
	case AU, Feet, Lightseconds, Lightminutes, Lightyears, Meters, Miles, NauticalMiles:
		return LengthUnits
	case KilometersPerHour, MetersPerSecond, MilesPerHour:
		return SpeedUnits
	default:
		return GeneralUnits
	}
}

func (self Unit) IsValid() bool {
	return (self != Invalid)
}

func (self Unit) String() string {
	switch self {
	case AU:
		return `AU`
	case Celsius:
		return `°C`
	case Fahrenheit:
		return `°F`
	case Feet:
		return `ft`
	case Kelvin:
		return `°K`
	case KilometersPerHour:
		return `kph`
	case Lightseconds:
		return `Ls`
	case Lightminutes:
		return `Lm`
	case Lightyears:
		return `Ly`
	case Meters:
		return `m`
	case MetersPerSecond:
		return `m/s`
	case Miles:
		return `mi`
	case MilesPerHour:
		return `mph`
	case NauticalMiles:
		return `nm`
	}

	return ``
}

func (self Unit) LongString() string {
	switch self {
	case AU:
		return `au`
	case Celsius:
		return `celsius`
	case Fahrenheit:
		return `fahrenhiet`
	case Feet:
		return `feet`
	case Kelvin:
		return `kelvin`
	case KilometersPerHour:
		return `kilometers per hour`
	case Lightseconds:
		return `lightseconds`
	case Lightminutes:
		return `lightminutes`
	case Lightyears:
		return `lightyears`
	case Meters:
		return `meters`
	case MetersPerSecond:
		return `meters per second`
	case Miles:
		return `miles`
	case MilesPerHour:
		return `miles per hour`
	case NauticalMiles:
		return `nautical miles`
	}

	return ``
}

func MustParseUnit(unit string) Unit {
	if unit, err := ParseUnit(unit); err == nil {
		return unit
	} else {
		panic(err.Error())
	}
}

func ParseUnit(unit string) (Unit, error) {
	unit = strings.TrimSpace(unit)

	switch unit {
	case `C`:
		return Celsius, nil
	case `F`:
		return Fahrenheit, nil
	case `K`:
		return Kelvin, nil
	case `KPH`:
		return KilometersPerHour, nil
	case `MPH`:
		return MilesPerHour, nil
	case `Ly`:
		return Lightyears, nil
	case `Lm`:
		return Lightminutes, nil
	case `Ls`:
		return Lightseconds, nil
	case `M`, `NM`:
		return NauticalMiles, nil
	}

	unit = strings.ToLower(unit)
	unit = strings.Map(func(r rune) rune {
		switch r {
		case '-', '/', '°':
			return r
		default:
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return r
			}
		}

		return -1
	}, unit)

	switch unit {
	case `au`, `astronomical unit`, `astronomical units`:
		return AU, nil
	case `celsius`, `°c`, `c°`:
		return Celsius, nil
	case `fahrenheit`, `°f`, `f°`:
		return Fahrenheit, nil
	case `foot`, `feet`:
		return Feet, nil
	case `kelvin`, `°k`, `k°`:
		return Kelvin, nil
	case `kilometer/hour`, `kilometer per hour`, `kilometers/hour`, `kilometers per hour`, `kph`, `km/h`:
		return KilometersPerHour, nil
	case `lightminute`, `lightminutes`, `lm`:
		return Lightminutes, nil
	case `lightsecond`, `lightseconds`, `ls`:
		return Lightseconds, nil
	case `lightyear`, `lightyears`, `ly`:
		return Lightyears, nil
	case `meter`, `meters`, `m`:
		return Meters, nil
	case `mile`, `miles`, `mi`:
		return Miles, nil
	case `mile/hour`, `mile per hour`, `miles/hour`, `miles per hour`, `mph`:
		return MilesPerHour, nil
	case `meter per second`, `meter/second`, `meters per second`, `meters/second`, `mps`, `m/s`:
		return MetersPerSecond, nil
	case `nautical mile`, `nautical miles`, `nm`, `nmi`:
		return NauticalMiles, nil
	}

	return Invalid, fmt.Errorf("Cannot parse unit '%s'", unit)
}
