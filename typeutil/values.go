package typeutil

import (
	"fmt"
	"github.com/ghetzel/go-stockutil/stringutil"
	"reflect"
	"strings"
)

// Returns whether the given value represents the underlying type's zero value
func IsZero(value interface{}) bool {
	if value == nil {
		return true
	}

	return reflect.DeepEqual(
		value,
		reflect.Zero(reflect.TypeOf(value)).Interface(),
	)
}

// Returns whether the given value is "empty" in the semantic sense. Zero values
// are considered empty, as are arrays, slices, and maps containing only empty
// values (called recursively). Finally, strings are trimmed of whitespace and
// considered empty if the result is zero-length.
//
func IsEmpty(value interface{}) bool {
	valueV := reflect.ValueOf(value)

	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}

	// short circuit for zero values of certain types
	switch valueV.Kind() {
	case reflect.Struct:
		if IsZero(value) {
			return true
		}
	}

	switch valueV.Kind() {
	case reflect.Array, reflect.Slice:
		if valueV.Len() == 0 {
			return true
		} else {
			for i := 0; i < valueV.Len(); i++ {
				if indexV := valueV.Index(i); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
					return false
				}
			}

			return true
		}

	case reflect.Map:
		if valueV.Len() == 0 {
			return true
		} else {
			for _, keyV := range valueV.MapKeys() {
				if indexV := valueV.MapIndex(keyV); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
					return false
				}
			}

			return true
		}

	case reflect.Chan:
		if valueV.Len() == 0 {
			return true
		}

	case reflect.String:
		if len(strings.TrimSpace(fmt.Sprintf("%v", value))) == 0 {
			return true
		}
	}

	return false
}

func RelaxedEqual(first interface{}, second interface{}) (bool, error) {
	if reflect.DeepEqual(first, second) {
		return true, nil
	} else if stringutil.IsNumeric(first) && stringutil.IsNumeric(second) {
		if fV, err := stringutil.ConvertToFloat(first); err == nil {
			if sV, err := stringutil.ConvertToFloat(second); err == nil {
				return (fV == sV), nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	} else if stringutil.IsBooleanTrue(first) && stringutil.IsBooleanTrue(second) {
		return true, nil
	} else if stringutil.IsBooleanFalse(first) && stringutil.IsBooleanFalse(second) {
		return true, nil
	} else if stringutil.IsTime(first) && stringutil.IsTime(second) {
		if fV, err := stringutil.ConvertToTime(first); err == nil {
			if sV, err := stringutil.ConvertToTime(second); err == nil {
				return fV.Equal(sV), nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	} else {
		if fV, err := stringutil.ToString(first); err == nil {
			if sV, err := stringutil.ToString(second); err == nil {
				return (fV == sV), nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}

	return false, nil
}

func IsKind(in interface{}, kinds ...reflect.Kind) bool {
	inV := reflect.ValueOf(in)

	if inV.IsValid() {
		if inV.Kind() == reflect.Ptr {
			inV = inV.Elem()
		}

		if inV.Kind() == reflect.Interface {
			inV = inV.Elem()
		}

		inT := inV.Type()

		if inT == nil {
			return false
		}

		for _, k := range kinds {
			if inT.Kind() == k {
				return true
			}
		}
	}

	return false
}

// Returns whether the given value is a slice or array.
func IsArray(in interface{}) bool {
	return IsKind(in, reflect.Slice, reflect.Array)
}

// Returns whether the given value is a function of any kind
func IsFunction(in interface{}) bool {
	return IsKind(in, reflect.Func)
}

// Returns whether the given value is a function.  If inParams is not -1, the function must
// accept that number of arguments.  If outParams is not -1, the function must return that
// number of values.
func IsFunctionArity(in interface{}, inParams int, outParams int) bool {
	if IsKind(in, reflect.Func) {
		inT := reflect.TypeOf(in)

		if inParams < 0 || inParams >= 0 && inT.NumIn() == inParams {
			if outParams < 0 || outParams >= 0 && inT.NumOut() == outParams {
				return true
			}
		}
	}

	return false
}
