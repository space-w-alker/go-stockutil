package convutil

import (
	"encoding/json"
	"fmt"

	"github.com/ghetzel/go-stockutil/typeutil"
)

var SpeedDisplayUnit = MeasurementSystem(Imperial)

type Speed float64

const (
	SpeedMetersPerSecond = 1
	SpeedKPH             = 0.277778
	SpeedFeetPerSecond   = 0.3048
	SpeedMPH             = 0.44704
	SpeedMach            = 340.29
)

func (self Speed) SlowerThan(other Speed) bool {
	return (self < other)
}

func (self Speed) FasterThan(other Speed) bool {
	return (self > other)
}

func (self Speed) Equal(other Speed) bool {
	return (self == other)
}

func (self Speed) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		`value`:   float64(self),
		`display`: self.String(),
	})
}

func (self *Speed) UnmarshalJSON(data []byte) error {
	if typeutil.IsScalar(string(data)) {
		var v float64

		if err := json.Unmarshal(data, &v); err == nil {
			*self = Speed(v)
			return nil
		} else {
			return err
		}
	} else {
		var in map[string]interface{}

		if err := json.Unmarshal(data, &in); err == nil {
			*self = Speed(typeutil.V(in[`value`]).Float())
			return nil
		} else {
			return err
		}
	}
}

func (self Speed) String() string {
	switch SpeedDisplayUnit {
	case Metric:
		switch {
		case self >= SpeedKPH:
			return fmt.Sprintf("%.0f kph", self/SpeedKPH)

		case self >= (SpeedMach * 0.97):
			return fmt.Sprintf("Mach %1.2f", self/SpeedMach)
		}

	case Imperial:
		switch {
		case self < SpeedMPH:
			return fmt.Sprintf("%.0f f/s", self/SpeedFeetPerSecond)

		case self >= (SpeedMach * 0.97):
			return fmt.Sprintf("Mach %1.2f", self/SpeedMach)

		default:
			return fmt.Sprintf("%.0f mph", self/SpeedMPH)
		}
	}

	return fmt.Sprintf("%.0f m/s", self)
}
