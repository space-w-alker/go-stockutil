package maputil

import (
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
)

type Map struct {
	data interface{}
}

// Create a new variadic map object from the given value (which should be a map of some kind).
func M(data interface{}) *Map {
	return &Map{
		data: data,
	}
}

// Retrieve a value from the Map by the given dot.separated key, or return a fallback
// value.  Return values are a typeutil.Variadic, which can be easily coerced into
// various types.
func (self *Map) Get(key string, fallbacks ...interface{}) typeutil.Variadic {
	if v := DeepGet(self.data, strings.Split(key, `.`), fallbacks...); v != nil {
		return typeutil.Variadic{
			Value: v,
		}
	} else {
		return typeutil.Variadic{}
	}
}

// Return the value at key as an automatically converted value.
func (self *Map) Auto(key string, fallbacks ...interface{}) interface{} {
	return self.Get(key, fallbacks...).Auto()
}

// Return the value at key as a string.
func (self *Map) String(key string, fallbacks ...interface{}) string {
	return self.Get(key, fallbacks...).String()
}

// Return the value at key interpreted as a Time.
func (self *Map) Time(key string, fallbacks ...interface{}) time.Time {
	return self.Get(key, fallbacks...).Time()
}

// Return the value at key as a bool.
func (self *Map) Bool(key string) bool {
	return self.Get(key).Bool()
}

// Return the value at key as an integer.
func (self *Map) Int(key string, fallbacks ...interface{}) int64 {
	return self.Get(key, fallbacks...).Int()
}

// Return the value at key as a float.
func (self *Map) Float(key string, fallbacks ...interface{}) float64 {
	return self.Get(key, fallbacks...).Float()
}

// Return the value at key as a slice.  Scalar values will be returned as a slice containing
// only that value.
func (self *Map) Slice(key string) []typeutil.Variadic {
	return self.Get(key).Slice()
}

// Return the value at key as a Map.  If the resulting value is nil or not a
// map type, a null Map will be returned.  All values retrieved from a null
// Map will return that type's zero value.
func (self *Map) Map(key string) *Map {
	if v := self.Get(key); v.Value != nil {
		if typeutil.IsMap(v) {
			return M(v)
		}
	}

	return M(nil)
}
