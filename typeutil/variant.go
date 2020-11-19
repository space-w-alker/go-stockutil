package typeutil

import (
	"encoding/json"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/ghetzel/go-stockutil/utils"
)

// Represents an interface type with helper functions for making it easy to do
// type conversions.
type Variant struct {
	Value interface{}
}

// Shortcut for creating a Variant.
func V(value interface{}) Variant {
	return Variant{
		Value: value,
	}
}

// Returns whether the underlying value is nil.
func (self Variant) IsNil() bool {
	return (self.Value == nil)
}

// Returns whether the underlying value is a zero value.
func (self Variant) IsZero() bool {
	return IsZero(self.Value)
}

// Return the value as a string, or an empty string if the value could not be converted.
func (self Variant) String() string {
	if v, err := utils.ConvertToString(self.Value); err == nil {
		return v
	} else {
		return ``
	}
}

// Return true if the value can be interpreted as a boolean true value, or false otherwise.
func (self Variant) Bool() bool {
	if v, err := utils.ConvertToBool(self.Value); err == nil {
		return v
	} else if vF, err := utils.ConvertToFloat(self.Value); err == nil && vF == 0 {
		return false
	} else {
		// use a more relaxed set of values for determining "true" because
		// the user has very explicitly asked us to try
		truthy, _ := template.IsTrue(self.Value)
		return truthy
	}
}

// Return the value as a float if it can be interpreted as such, or 0 otherwise.
func (self Variant) Float() float64 {
	if v, err := utils.ConvertToFloat(self.Value); err == nil {
		return v
	} else {
		return 0
	}
}

// Return the value as an integer if it can be interpreted as such, or 0 otherwise. Float values
// will be truncated to integers.
func (self Variant) Int() int64 {
	if v, err := utils.ConvertToFloat(self.Value); err == nil {
		return int64(v)
	} else {
		return 0
	}
}

// Return the value as a native integer if it can be interpreted as such, or 0 otherwise.
// Float values will be truncated to integers.
func (self Variant) NInt() int {
	return int(self.Int())
}

// Return the value as a slice of Variants. Scalar types will return a slice containing
// a single Variant element representing the value.
func (self Variant) Slice() []Variant {
	values := make([]Variant, 0)

	utils.SliceEach(self.Value, func(_ int, v interface{}) error {
		values = append(values, Variant{
			Value: v,
		})
		return nil
	})

	return values
}

// Same as Slice(), but returns a []string.
func (self Variant) Strings() []string {
	values := self.Slice()
	out := make([]string, len(values))

	for i, value := range values {
		out[i] = value.String()
	}

	return out
}

// Converts the value to a string, then splits on the given delimiter.
func (self Variant) Split(on string) []string {
	if s := self.String(); s == `` {
		return nil
	} else {
		return strings.Split(s, on)
	}
}

// Return the value converted to an error, or nil if it is not an error.
func (self Variant) Err() error {
	if err, ok := self.Value.(error); ok {
		return err
	} else {
		return nil
	}
}

// Return the value automatically converted to the appropriate type.
func (self Variant) Auto() interface{} {
	return utils.Autotype(self.Value)
}

// Return the value as a time.Time if it can be interpreted as such, or zero time otherwise.
func (self Variant) Time() time.Time {
	if v, err := utils.ConvertToTime(self.Value); err == nil {
		return v
	} else {
		return time.Time{}
	}
}

// Return the value as a time.Duration if it can be interpreted as such, or zero otherwise.
func (self Variant) Duration() time.Duration {
	if v, err := utils.ParseDuration(self.String()); err == nil {
		return v
	} else {
		return 0
	}
}

// Return the value at key as a byte slice.
func (self Variant) Bytes() []byte {
	if v, err := utils.ConvertToBytes(self.Value); err == nil {
		return v
	} else {
		return []byte{}
	}
}

func (self Variant) Interface() interface{} {
	out := self.Value

	for {
		if subvariant, ok := out.(Variant); ok {
			out = subvariant.Value
		} else if subvariant, ok := out.(*Variant); ok {
			out = subvariant.Value
		} else {
			break
		}
	}

	return out
}

// Return the value as a map[Variant]Variant if it can be interpreted as such, or an empty map otherwise.
// If the variant contains a struct, and a tagName is specified, the key names of the output map will be taken
// from the struct field's tag value, consistent with the rules used in encoding/json.
func (self Variant) Map(tagName ...string) map[Variant]Variant {
	output := make(map[Variant]Variant)
	var tn string

	if len(tagName) > 0 && tagName[0] != `` {
		tn = tagName[0]
	}

	if IsMap(self.Value) {
		mapV := reflect.ValueOf(self.Value)

		for _, key := range mapV.MapKeys() {
			if key.CanInterface() {
				if value := mapV.MapIndex(key); value.CanInterface() {
					output[V(key.Interface())] = V(value.Interface())
				}
			}
		}
	} else if elem := ResolveValue(self.Value); IsStruct(elem) {
		structV := reflect.ValueOf(elem)
		structT := structV.Type()

	FieldLoop:
		for i := 0; i < structT.NumField(); i++ {
			if structF := structT.Field(i); !structF.Anonymous {
				if value := structV.Field(i); value.IsValid() {
					key := structF.Name

					if tn != `` {
						parts := strings.Split(structF.Tag.Get(tn), `,`)

						if tag := parts[0]; tag != `` {
							key = tag
						}

						if len(parts) > 1 {
							for _, p := range parts {
								switch p {
								case `omitempty`:
									if IsZero(value) {
										continue FieldLoop
									}
								}
							}
						}
					}

					if value.CanInterface() {
						output[V(key)] = V(value.Interface())
					}
				}
			}
		}
	}

	return output
}

// Return the value as a map[string]interface{} if it can be interpreted as such, or an empty map otherwise.
func (self Variant) MapNative(tagName ...string) map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range self.Map(tagName...) {
		out[k.String()] = v.Interface()
	}

	return out
}

// Satisfy the json.Marshaler interface
func (self Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Auto())
}

// Package-level string converter
func String(in interface{}) string {
	return V(in).String()
}

// Package-level bool converter
func Bool(in interface{}) bool {
	return V(in).Bool()
}

// Package-level float converter
func Float(in interface{}) float64 {
	return V(in).Float()
}

// Package-level int64 converter
func Int(in interface{}) int64 {
	return V(in).Int()
}

// Package-level native int converter
func NInt(in interface{}) int {
	return V(in).NInt()
}

// Package-level slice converter
func Slice(in interface{}) []Variant {
	return V(in).Slice()
}

// Package-level string slice converter
func Strings(in interface{}) []string {
	return V(in).Strings()
}

// Package-level string splitter.
func Split(in interface{}, on string) []string {
	return V(in).Split(on)
}

// Package-level error converter
func Err(in interface{}) error {
	return V(in).Err()
}

// Package-level auto converter
func Auto(in interface{}) interface{} {
	return V(in).Auto()
}

// Package-level time converter
func Time(in interface{}) time.Time {
	return V(in).Time()
}

// Package-level duration converter
func Duration(in interface{}) time.Duration {
	return V(in).Duration()
}

// Package-level bytes converter
func Bytes(in interface{}) []byte {
	return V(in).Bytes()
}

// Package-level map converter
func Map(in interface{}, tagName ...string) map[Variant]Variant {
	return V(in).Map(tagName...)
}

// Package-level map[string]interface{} converter
func MapNative(in interface{}, tagName ...string) map[string]interface{} {
	return V(in).MapNative(tagName...)
}
