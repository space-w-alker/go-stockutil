package sliceutil

import (
	"fmt"
	"github.com/ghetzel/go-stockutil/typeutil"
	"reflect"
)

var Stop = fmt.Errorf("stop iterating")

type IterationFunc func(i int, value interface{}) error                  // {}
type CompareFunc func(i int, first interface{}, second interface{}) bool // {}

var StrictEqualityCompare = func(_ int, first interface{}, second interface{}) bool {
	if first == second {
		return true
	}

	return false
}

var RelaxedEqualityCompare = func(_ int, first interface{}, second interface{}) bool {
	if v, err := typeutil.RelaxedEqual(first, second); err == nil && v == true {
		return true
	}

	return false
}

// Return whether the given slice contains the given value.  If a comparator is provided, it will
// be used to compare the elements.
//
func Contains(in interface{}, value interface{}, comparators ...CompareFunc) bool {
	if len(comparators) == 0 {
		comparators = []CompareFunc{StrictEqualityCompare}
	}

	comparator := comparators[0]

	if inV := reflect.ValueOf(in); inV.IsValid() {
		for i := 0; i < inV.Len(); i++ {
			if current := inV.Index(i); current.IsValid() {
				if comparator(i, value, current.Interface()) {
					return true
				}
			}
		}
	}

	return false
}

// Returns whether the given string slice contains a given string.
func ContainsString(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}

	return false
}

// Returns whether the given string slice contains any of the following strings.
func ContainsAnyString(list []string, elems ...string) bool {
	for _, e := range elems {
		if ContainsString(list, e) {
			return true
		}
	}

	return false
}

// Returns whether the given string slice contains all of the following strings.
func ContainsAllStrings(list []string, elems ...string) bool {
	for _, e := range elems {
		if !ContainsString(list, e) {
			return false
		}
	}

	return true
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
func Stringify(in interface{}) []string {
	if !typeutil.IsArray(in) {
		return nil
	}

	inV := reflect.ValueOf(in)

	if inV.IsValid() {
		rv := make([]string, inV.Len())

		for i := 0; i < inV.Len(); i++ {
			if iV := inV.Index(i); iV.IsValid() {
				rv[i] = fmt.Sprintf("%v", iV.Interface())
			}
		}

		return rv
	} else {
		return nil
	}
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

func Each(slice interface{}, iterFn IterationFunc) error {
	if iterFn == nil {
		return nil
	}

	if typeutil.IsArray(slice) {
		sliceV := reflect.ValueOf(slice)

		for i := 0; i < sliceV.Len(); i++ {
			if err := iterFn(i, sliceV.Index(i).Interface()); err != nil {
				if err == Stop {
					return nil
				} else {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("Exected slice or array, got %T", slice)
	}

	return nil
}

// Returns a new slice with only unique elements from the given interface included.
func Unique(in interface{}) []interface{} {
	return unique(in, StrictEqualityCompare)
}

func unique(in interface{}, comparator CompareFunc) []interface{} {
	inV := reflect.ValueOf(in)
	values := make([]interface{}, 0)

	if inV.IsValid() {
	InputLoop:
		for i := 0; i < inV.Len(); i++ {
			if current := inV.Index(i); current.IsValid() {
				for _, existing := range values {
					if comparator(i, existing, current.Interface()) {
						continue InputLoop
					}
				}

				values = append(values, current.Interface())
			}
		}
	} else {
		return nil
	}

	return values
}
