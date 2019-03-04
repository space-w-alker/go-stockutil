// Utilities for working with and manipulating structs.
package structutil

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
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

	destS := structs.New(dest)
	srcS := structs.New(source)

	for _, sField := range srcS.Fields() {
		if sField.IsExported() {
			if dField, ok := destS.FieldOk(sField.Name()); ok {
				sValue := sField.Value()
				dValue := dField.Value()

				if typeutil.IsStruct(sValue) {
					if err := CopyFunc(dValue, sValue, fn); err != nil {
						return err
					}
				} else if repl, ok := fn(sField.Name(), sValue, dValue); ok {
					// set the destination field value to whatever came back from the function
					if err := dField.Set(repl); err != nil {
						if strings.HasSuffix(err.Error(), `is not settable`) {
							return fmt.Errorf("field %q is not settable", dField.Name())
						} else {
							return err
						}
					}
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
