package convutil

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var rxDistanceExtract = regexp.MustCompile(`^(?P<number>\d+(?:\.\d+)?)\s*(?P<unit>[\w\s]+)\W*$`)

type Distance float64

const (
	DistanceMeter        = 1
	DistanceKilometer    = 1000
	DistanceFoot         = 0.3048
	DistanceYard         = 0.9144
	DistanceMile         = 1609.344
	DistanceNauticalMile = 1852
)

var DistanceDisplayUnit = MeasurementSystem(Imperial)

func MustParseDistance(in interface{}) Distance {
	if distance, err := ParseDistance(in); err == nil {
		return distance
	} else {
		panic(`invalid distance: ` + err.Error())
	}
}

func ParseDistance(in interface{}) (Distance, error) {
	if typeutil.IsZero(in) {
		return 0, nil
	}

	if match := rxutil.Match(rxDistanceExtract, strings.TrimSpace(fmt.Sprintf("%v", in))); match != nil {
		if v := typeutil.V(match.Group(`number`)).Float(); v >= 0 {
			unit := match.Group(`unit`)
			unit = strings.TrimSpace(unit)
			unit = strings.ToLower(unit)
			unit = strings.TrimSuffix(unit, `s`)

			switch unit {
			case `meter`, `m`:
				return Distance(v), nil
			case `kilometer`, `km`:
				return Distance(v) * DistanceKilometer, nil
			case `mile`, `mi`:
				return Distance(v) * DistanceMile, nil
			case `feet`, `foot`, `ft`:
				return Distance(v) * DistanceFoot, nil
			case `yard`, `yd`:
				return Distance(v) * DistanceYard, nil
			case `nm`, `nautical mile`:
				return Distance(v) * DistanceNauticalMile, nil
			default:
				return 0, fmt.Errorf("Unrecognized distance unit %q", unit)
			}
		} else {
			return 0, fmt.Errorf("Unable to extract number from distance value")
		}
	} else if v := typeutil.V(in).Float(); v >= 0 {
		return Distance(v), nil
	} else {
		return 0, fmt.Errorf("unable to parse distance value")
	}
}

func (self Distance) Within(other Distance) bool {
	return (self <= other)
}

func (self Distance) Beyond(other Distance) bool {
	return (self > other)
}

func (self Distance) Equal(other Distance) bool {
	return (self == other)
}

func (self Distance) MetricString() string {
	switch {
	case self >= 10*DistanceKilometer:
		return fmt.Sprintf("%.0f kilometers", self/DistanceKilometer)

	case self > DistanceKilometer:
		return fmt.Sprintf("%.1f kilometers", self/DistanceKilometer)

	case self == DistanceKilometer:
		return fmt.Sprintf("%.0f kilometer", self/DistanceKilometer)

	default:
		return fmt.Sprintf("%.0f meters", self)

	}
}

func (self Distance) ImperialString() string {
	switch {
	case self >= 5*DistanceMile:
		return fmt.Sprintf("%.0f miles", self/DistanceMile)

	case self >= (0.9*DistanceMile) && self <= (1.1*DistanceMile):
		return fmt.Sprintf("%.0f mile", self/DistanceMile)

	case self >= (0.1 * DistanceMile):
		return fmt.Sprintf("%.1f miles", self/DistanceMile)

	default:
		return fmt.Sprintf("%.0f feet", self/DistanceFoot)
	}
}

func (self Distance) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		`value`: float64(self),
		`units`: map[string]interface{}{
			`default`:  self.String(),
			`imperial`: self.ImperialString(),
			`metric`:   self.MetricString(),
		},
	})
}

func (self *Distance) UnmarshalJSON(data []byte) error {
	if typeutil.IsScalar(string(data)) {
		var v float64

		if err := json.Unmarshal(data, &v); err == nil {
			*self = Distance(v)
			return nil
		} else {
			return err
		}
	} else {
		var in map[string]interface{}

		if err := json.Unmarshal(data, &in); err == nil {
			*self = Distance(typeutil.V(in[`value`]).Float())
			return nil
		} else {
			return err
		}
	}
}

func (self Distance) String() string {
	switch DistanceDisplayUnit {
	case Metric:
		return self.MetricString()

	case Imperial:
		return self.ImperialString()
	}

	return fmt.Sprintf("%f meters", self)
}
