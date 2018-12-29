package convutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpeedFuncs(t *testing.T) {
	assert := require.New(t)

	assert.True(Speed(SpeedMPH).FasterThan(SpeedKPH))
	assert.True(Speed(SpeedKPH).SlowerThan(SpeedMPH))
	assert.True(Speed(3600 * SpeedMPH).Equal(5280 * SpeedFeetPerSecond))
}

func TestParseDistance(t *testing.T) {
	assert := require.New(t)

	assert.Zero(MustParseDistance(``))
	assert.Zero(MustParseDistance(`0m`))
	assert.Zero(MustParseDistance(`0 meter`))
	assert.Zero(MustParseDistance(`0 meters`))
	assert.Zero(MustParseDistance(`0meter`))
	assert.Zero(MustParseDistance(`0meters`))

	assert.EqualValues(1, MustParseDistance(`1m`))
	assert.EqualValues(1, MustParseDistance(`1 meter`))
	assert.EqualValues(1, MustParseDistance(`1 meters`))
	assert.EqualValues(1, MustParseDistance(`1meter`))
	assert.EqualValues(1, MustParseDistance(`1meters`))

	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625km`))
	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625 km`))
	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625 km.`))
	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625 kilometer`))
	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625 kilometers`))
	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625kilometer`))
	assert.EqualValues(Distance(3.141597625)*DistanceKilometer, MustParseDistance(`3.141597625kilometers`))

	assert.EqualValues(Distance(26.2)*DistanceMile, MustParseDistance(`26.2mi`))
	assert.EqualValues(Distance(26.2)*DistanceMile, MustParseDistance(`26.2 miles`))

	assert.EqualValues(Distance(5280)*DistanceFoot, MustParseDistance(`5280ft`))
	assert.EqualValues(Distance(5280)*DistanceFoot, MustParseDistance(`5280 ft`))
	assert.EqualValues(Distance(5280)*DistanceFoot, MustParseDistance(`5280 ft.`))
	assert.EqualValues(Distance(5280)*DistanceFoot, MustParseDistance(`5280 feet`))
	assert.EqualValues(Distance(5280)*DistanceFoot, MustParseDistance(`5280feet`))

	assert.EqualValues(Distance(100)*DistanceYard, MustParseDistance(`100yd.`))
	assert.EqualValues(Distance(100)*DistanceYard, MustParseDistance(`100 yd.`))
	assert.EqualValues(Distance(100)*DistanceYard, MustParseDistance(`100 yd`))
	assert.EqualValues(Distance(100)*DistanceYard, MustParseDistance(`100yd`))
	assert.EqualValues(Distance(100)*DistanceYard, MustParseDistance(`100 yards`))

	assert.EqualValues(Distance(300)*DistanceNauticalMile, MustParseDistance(`300nm`))
	assert.EqualValues(Distance(300)*DistanceNauticalMile, MustParseDistance(`300 nm.`))
	assert.EqualValues(Distance(300)*DistanceNauticalMile, MustParseDistance(`300 nm`))
	assert.EqualValues(Distance(300)*DistanceNauticalMile, MustParseDistance(`300 nautical miles`))
	assert.EqualValues(Distance(300)*DistanceNauticalMile, MustParseDistance(`300 nautical mile`))
}
