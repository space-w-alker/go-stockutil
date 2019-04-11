package executil

import (
	"os"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
)

func env(name string, fallback ...interface{}) typeutil.Variant {
	if v := os.Getenv(name); v != `` {
		return typeutil.V(v)
	} else if len(fallback) > 0 {
		return typeutil.V(fallback[0])
	} else {
		return typeutil.V(nil)
	}
}

func Env(name string, fallback ...interface{}) string {
	return env(name, fallback...).String()
}

func EnvInt(name string, fallback ...interface{}) int64 {
	return env(name, fallback...).Int()
}

func EnvFloat(name string, fallback ...interface{}) float64 {
	return env(name, fallback...).Float()
}

func EnvBool(name string, fallback ...interface{}) bool {
	return env(name, fallback...).Bool()
}

func EnvTime(name string, fallback ...interface{}) time.Time {
	return env(name, fallback...).Time()
}

func EnvDuration(name string, fallback ...interface{}) time.Duration {
	return env(name, fallback...).Duration()
}
