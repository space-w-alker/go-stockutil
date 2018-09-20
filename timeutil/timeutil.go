// Utilities for messing with time.
package timeutil

import (
	"time"

	"github.com/ghetzel/go-stockutil/utils"
)

// Provides an API-compatible version of time.ParseDuration that accepts additional
// formats for parsing durations:
//
//   1y, 3year, 5years: Expands to (n*24*365) hours
//   1w, 3weeks, 5wks:  Expands to (n*24*7) hours
//   1d, 1day, 5days:   Expands to (n*24) hours
//
func ParseDuration(in string) (time.Duration, error) {
	return utils.ParseDuration(in)
}
