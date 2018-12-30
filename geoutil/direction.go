package geoutil

type CardinalDirection string

const (
	North     CardinalDirection = `N`
	NorthEast                   = `NE`
	East                        = `E`
	SouthEast                   = `SE`
	South                       = `S`
	SouthWest                   = `SW`
	West                        = `W`
	NorthWest                   = `NW`
)

func GetDirectionFromBearing(bearing float64) CardinalDirection {
	switch {
	case (bearing >= 0 && bearing <= 22.5) || (bearing > 337.5 && bearing <= 360):
		return North
	case bearing > 22.5 && bearing <= 67.5:
		return NorthEast
	case bearing > 67.5 && bearing <= 112.5:
		return East
	case bearing > 112.5 && bearing <= 157.5:
		return SouthEast
	case bearing > 157.5 && bearing <= 202.5:
		return South
	case bearing > 202.5 && bearing <= 247.5:
		return SouthWest
	case bearing > 247.5 && bearing <= 292.5:
		return West
	case bearing > 292.5 && bearing <= 337.5:
		return NorthWest
	}

	return ``
}
