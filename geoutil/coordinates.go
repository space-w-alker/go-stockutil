package geoutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/kellydunn/golang-geo"
)

// Describes the error margin (+/-) for each of the named values.
type LocationError struct {
	Latitude  Distance `json:"latitude"`
	Longitude Distance `json:"longitude"`
	Altitude  Distance `json:"altitude"`
	Bearing   float64  `json:"bearing"`
	Speed     Speed    `json:"speed"`
	Timestamp float64  `json:"timestamp"`
}

// Specifies a three-dimensional location within a coordinate reference system.
type Location struct {
	Latitude   float64                `json:"latitude,omitempty"`
	Longitude  float64                `json:"longitude,omitempty"`
	Bearing    float64                `json:"bearing,omitempty"`
	Timestamp  time.Time              `json:"timestamp,omitempty"`
	Altitude   Distance               `json:"altitude,omitempty"`
	Speed      Speed                  `json:"speed,omitempty"`
	Accuracy   float64                `json:"accuracy,omitempty"`
	Error      *LocationError         `json:"error,omitempty"`
	Direction  CardinalDirection      `json:"direction,omitempty"`
	Source     string                 `json:"source,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func NewLocation(latitude float64, longitude float64) *Location {
	return &Location{
		Latitude:   latitude,
		Longitude:  longitude,
		Properties: make(map[string]interface{}),
	}
}

func (self *Location) String() string {
	return strings.TrimSpace(fmt.Sprintf(
		"%f,%f %s",
		self.Latitude,
		self.Longitude,
		maputil.Join(self.Properties, `=`, ` `),
	))
}

func (self *Location) MarshalJSON() ([]byte, error) {
	if self.Bearing < 0 {
		self.Bearing = 360 + self.Bearing
	}

	self.Direction = self.CardinalDirection()

	type Alias Location

	return json.Marshal(&struct {
		*Alias
	}{
		(*Alias)(self),
	})
}

func (self *Location) HasCoordinates() bool {
	if self.Latitude == 0 && self.Longitude == 0 {
		return false
	}

	return true
}

func (self *Location) CardinalDirection() CardinalDirection {
	return GetDirectionFromBearing(self.Bearing)
}

// Return the distance (in meters) between this point and another.  This calulates the
// great-circle distance (shortest distance two points on the surface of a sphere) between
// this Location and another.  Since this (incorrectly) assumes the Earth to be a true
// sphere, this is only reasonably accurate for short-ish distances (is only accurate to
// within ~0.5%).
//
func (self *Location) HaversineDistance(other *Location) Distance {
	if !self.HasCoordinates() || !other.HasCoordinates() {
		panic("Coordinates not specified")
	}

	return Distance((geo.NewPoint(self.Latitude, self.Longitude).GreatCircleDistance(
		geo.NewPoint(other.Latitude, other.Longitude),
	) * Kilometer))
}

func (self *Location) BearingTo(other *Location) float64 {
	if !self.HasCoordinates() || !other.HasCoordinates() {
		panic("Coordinates not specified")
	}

	return geo.NewPoint(self.Latitude, self.Longitude).BearingTo(
		geo.NewPoint(other.Latitude, other.Longitude),
	)
}

func (self *Location) SpeedFrom(other *Location) Speed {
	if self.Timestamp.IsZero() {
		return 0
	}

	if other.Timestamp.IsZero() {
		return 0
	}

	if !other.Timestamp.Before(self.Timestamp) && !other.Timestamp.Equal(self.Timestamp) {
		return 0
	}

	delta := self.Timestamp.Sub(other.Timestamp)
	distance := self.HaversineDistance(other)

	// speed is distance (in meters) / time delta (in seconds); meters/sec.
	return Speed(float64(distance) / delta.Seconds())
}

func NullIsland() *Location {
	return &Location{
		Latitude:  0,
		Longitude: 0,
		Bearing:   0,
		Altitude:  0,
		Accuracy:  1,
		Source:    `null`,
		Timestamp: time.Now(),
	}
}
