package sliceutil

import (
	"fmt"
	"reflect"
)

func ContainsString(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}

	return false
}

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
