// Utilities for working with and manipulating structs.
package structutil

import (
	"fmt"
	"reflect"

	"github.com/ghetzel/go-stockutil/typeutil"
)

// Receives a struct field name, the value of that field in the source struct, and the value for that field in the destination struct.
// Returns the value that should be placed in the destination struct fields.  If the returned bool is false, no changes will
// be made.
type StructValueFunc func(field string, sourceValue interface{}, destValue interface{}) (interface{}, bool)

func CopyFunc(dest interface{}, source interface{}, fn StructValueFunc) error {
	if dest == nil || source == nil || fn == nil {
		return nil
	}

	var destV reflect.Value
	var srcV reflect.Value

	if dV, ok := dest.(reflect.Value); ok {
		destV = dV
	} else {
		destV = reflect.ValueOf(dest)
	}

	if sV, ok := source.(reflect.Value); ok {
		srcV = sV
	} else {
		srcV = reflect.ValueOf(source)
	}

	if dV, err := validatePtrToStruct(`destination`, destV); err == nil {
		destV = dV
	} else {
		return err
	}

	if sV, err := validatePtrToStruct(`source`, srcV); err == nil {
		srcV = sV
	} else {
		return err
	}

	destT := destV.Type()
	srcT := srcV.Type()

	for s := 0; s < srcT.NumField(); s++ {
		sFieldT := srcT.Field(s)
		sFieldV := srcV.Field(s)

		// only exported field names leave this empty, so skip if it's not (i.e.: we have an unexported field)
		if sFieldT.PkgPath != `` {
			continue
		}

		if dFieldT, ok := destT.FieldByName(sFieldT.Name); ok {
			dFieldV := destV.FieldByName(dFieldT.Name)

			if dFieldT.Anonymous {
				if err := CopyFunc(dFieldV, sFieldV, fn); err != nil {
					return err
				}
			} else {
				if sFieldV.CanInterface() && dFieldV.CanInterface() {
					if repl, ok := fn(dFieldT.Name, sFieldV.Interface(), dFieldV.Interface()); ok {
						if dFieldV.CanSet() {
							if err := typeutil.SetValue(dFieldV, repl); err != nil {
								return err
							}
						} else {
							return fmt.Errorf("field %q is not settable", dFieldT.Name)
						}
					}
				} else {
					return fmt.Errorf("Cannot retrieve field value %q", dFieldT.Name)
				}
			}
		}
	}

	return nil
}

// Copy all values from the source into the destination, provided the source value for the corresponding
// field is not that type's zero value.
func CopyNonZero(dest interface{}, source interface{}) error {
	return CopyFunc(dest, source, func(name string, s interface{}, d interface{}) (interface{}, bool) {
		if typeutil.IsZero(s) {
			return nil, false
		} else {
			return s, true
		}
	})
}

func validatePtrToStruct(name string, obj reflect.Value) (reflect.Value, error) {
	if obj.Kind() == reflect.Ptr {
		if obj.Elem().Kind() == reflect.Struct {
			return obj.Elem(), nil
		} else {
			return reflect.Value{}, fmt.Errorf("bad %s: expected pointer to struct", name)
		}
	} else {
		return reflect.Value{}, fmt.Errorf("bad %s: expected pointer to struct", name)
	}
}
