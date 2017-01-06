package sliceutil

import (
	"fmt"
	"reflect"
)

// Returns whether the given string slice contains a given string.
func ContainsString(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}

	return false
}

// Removes all elements from the given interface slice that are "empty", which is
// defined as being nil, a nil or zero-length array, chan, map, slice, or string.
//
// The zero values of any other type are not considered empty and will remain in
// the return value.
//
func Compact(in []interface{}) []interface{} {
	if in == nil {
		return nil
	}

	rv := make([]interface{}, 0)

	for _, v := range in {
		if v != nil {
			vV := reflect.ValueOf(v)

			switch vV.Kind() {
			case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
				if vV.Len() > 0 {
					rv = append(rv, v)
				}
			default:
				rv = append(rv, v)
			}
		}
	}

	return rv
}

// Removes all zero-length strings from the given string slice, returning a new
// slice with the values removed.
func CompactString(in []string) []string {
	if in == nil {
		return nil
	}

	rv := make([]string, 0)

	for _, v := range in {
		if v != `` {
			rv = append(rv, v)
		}
	}

	return rv
}

// Converts all elements of the given interface slice to strings using the "%v"
// format string via the fmt package.
func Stringify(in []interface{}) []string {
	if in == nil {
		return nil
	}

	rv := make([]string, len(in))

	for i, v := range in {
		rv[i] = fmt.Sprintf("%v", v)
	}

	return rv
}

// Returns the first item that is not the zero value for that value's type.
func Or(in ...interface{}) interface{} {
	for _, v := range Compact(in) {
		// if the current value equals the zero value of its type,
		// then skip it, otherwise return it
		if v != reflect.Zero(reflect.TypeOf(v)).Interface() {
			return v
		}
	}

	return nil
}

// Returns the first item that is not a zero-length string.
func OrString(in ...string) string {
	if v := CompactString(in); len(v) > 0 {
		return v[0]
	} else {
		return ``
	}
}
