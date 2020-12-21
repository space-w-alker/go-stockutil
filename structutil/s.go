package structutil

import (
	"fmt"
	"reflect"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var StructTag string = `structutil`

type Field struct {
	*reflect.StructField
	Struct *Struct
}

func (self *Field) Kind() reflect.Kind {
	return self.Value().Kind()
}

func (self *Field) Value() reflect.Value {
	if src := self.Struct.srcval; src.IsValid() {
		if src.Kind() == reflect.Ptr {
			src = src.Elem()
		}

		if val := src.FieldByName(self.Name); val.IsValid() {
			return val
		}
	}

	return reflect.Value{}
}

func (self *Field) V() typeutil.Variant {
	if v := self.Value(); v.CanInterface() {
		return typeutil.V(v.Interface())
	}

	return typeutil.V(nil)
}

func (self *Field) Set(value interface{}) error {
	return typeutil.SetValue(self.Value(), value)
}

func (self *Field) MergeValue(in interface{}) error {
	var current = self.V()
	var other = typeutil.V(in)
	var newVal interface{}
	var trySet bool

	if current.IsArray() {
		if err := current.Append(other.Value); err == nil {
			newVal = current.Value
			trySet = true
		} else {
			return err
		}
	} else if current.IsMap() {
		if out, err := maputil.Merge(
			current.Value,
			other.Value,
		); err == nil {
			newVal = out
			trySet = true
		} else {
			return fmt.Errorf("cannot set value for field '%s': %v", self.Name, err)
		}
	} else if !other.IsZero() {
		newVal = other.Value
		trySet = true
	}

	if trySet {
		if err := self.Set(newVal); err != nil {
			return fmt.Errorf("cannot set value for field '%s': %v", self.Name, err)
		}
	}

	return nil
}

// A Struct, or "S-object", can be used to rapidly and safely inspect, iterate over, and modify values of a struct.
type Struct struct {
	Source    interface{}
	fields    []*Field
	fieldmap  map[string]*Field
	populated bool
	srcval    reflect.Value
}

func S(src interface{}) *Struct {
	return &Struct{
		Source:   src,
		fields:   make([]*Field, 0),
		fieldmap: make(map[string]*Field),
	}
}

func (self *Struct) Fields() []*Field {
	if !self.populated {
		self.srcval = reflect.ValueOf(self.Source)

		FieldsFunc(self.Source, func(field *reflect.StructField, value reflect.Value) error {
			var f = &Field{
				StructField: field,
				Struct:      self,
			}

			self.fields = append(self.fields, f)
			self.fieldmap[field.Name] = f

			return nil
		})

		self.populated = true
	}

	return self.fields
}

func (self *Struct) Field(name string) (*Field, bool) {
	self.Fields()

	if f, ok := self.fieldmap[name]; ok {
		return f, true
	}

	return nil, false
}

func (self *Struct) Merge(other interface{}) error {
	for _, otherField := range S(other).Fields() {
		if myField, ok := self.Field(otherField.Name); ok {
			var err error

			switch tag := myField.Tag.Get(StructTag); tag {
			case ``, `merge`:
				err = myField.MergeValue(otherField.V().Value)
			case `replace`:
				err = myField.Set(otherField.V().Value)
			default:
				err = fmt.Errorf("unknown directive %q", tag)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}
