package maputil

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	utilutil "github.com/ghetzel/go-stockutil/utils"
)

var MapXmlRootTagName = `data`

type ItemFunc func(key string, value typeutil.Variant) error

// A Map object (or "M" object) is a utility struct that makes it straightforward to
// work with interface data types that contain map-like data (has a reflect.Kind equal
// to reflect.Map).
type Map struct {
	data              interface{}
	structTagKey      string
	rootTagName       string
	xmlMarshalGeneric bool
}

// Create a new Variant map object from the given value (which should be a map of some kind).
func M(data interface{}) *Map {
	if dataV, ok := data.(typeutil.Variant); ok {
		data = dataV.Value
	} else if dataM, ok := data.(*Map); ok {
		return dataM
	} else if dataM, ok := data.(Map); ok {
		return &dataM
	} else if dataSM, ok := data.(*sync.Map); ok {
		dataM := make(map[string]interface{})
		dataSM.Range(func(key, value interface{}) bool {
			dataM[typeutil.String(key)] = value
			return true
		})

		data = dataM
	} else if uV, ok := data.(url.Values); ok {
		dataM := make(map[string]interface{})

		for k, v := range uV {
			switch len(v) {
			case 0:
				break
			case 1:
				dataM[k] = typeutil.Auto(v[0])
			default:
				dataM[k] = sliceutil.Autotype(v)
			}
		}

		data = dataM
	} else if hV, ok := data.(http.Header); ok {
		dataM := make(map[string]interface{})

		for k, v := range hV {
			switch len(v) {
			case 0:
				break
			case 1:
				dataM[k] = typeutil.Auto(v[0])
			default:
				dataM[k] = sliceutil.Autotype(v)
			}
		}

		data = dataM
	} else if data == nil {
		data = make(map[string]interface{})
	}

	return &Map{
		data:         data,
		structTagKey: UnmarshalStructTag,
	}
}

// Specify which struct tag to honor for generating field names when then
// underlying data is a struct.
func (self *Map) Tag(key string) *Map {
	self.structTagKey = key
	return self
}

// Return the underlying value the M-object was created with.
func (self *Map) Value() interface{} {
	return self.data
}

// Set a value in the Map at the given dot.separated key to a value.
func (self *Map) Set(key string, value interface{}) typeutil.Variant {
	vv := typeutil.V(value)
	self.data = DeepSet(self.data, strings.Split(key, `.`), vv)

	return vv
}

// Set a value in the Map at the given dot.separated key to a value, but only if the
// current value at that key is that type's zero value.
func (self *Map) SetIfZero(key string, value interface{}) (typeutil.Variant, bool) {
	if v := self.Get(key); v.IsZero() {
		return self.Set(key, value), true
	} else {
		return v, false
	}
}

// Set a value in the Map at the given dot.separated key to a value, but only if the
// new value is not a zero value.
func (self *Map) SetValueIfNonZero(key string, value interface{}) (typeutil.Variant, bool) {
	if !typeutil.IsZero(value) {
		return self.Set(key, value), true
	} else {
		return self.Get(key), false
	}
}

// Retrieve a value from the Map by the given dot.separated key, or return a fallback
// value.  Return values are a typeutil.Variant, which can be easily coerced into
// various types.
func (self *Map) Get(key string, fallbacks ...interface{}) typeutil.Variant {
	native := self.MapNative(self.structTagKey)

	if v, ok := native[key]; ok && v != nil {
		return typeutil.Variant{
			Value: v,
		}
	} else if v := DeepGet(self.data, strings.Split(key, `.`), fallbacks...); v != nil {
		return typeutil.Variant{
			Value: v,
		}
	} else {
		return typeutil.Variant{}
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

// Return the value at key interpreted as a Duration.
func (self *Map) Duration(key string, fallbacks ...interface{}) time.Duration {
	return self.Get(key, fallbacks...).Duration()
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

// Return the value at key as a byte slice.
func (self *Map) Bytes(key string) []byte {
	return self.Get(key).Bytes()
}

// Return the value at key as a slice.  Scalar values will be returned as a slice containing
// only that value.
func (self *Map) Slice(key string) []typeutil.Variant {
	return self.Get(key).Slice()
}

// Same as Slice(), but returns a []string
func (self *Map) Strings(key string) []string {
	return self.Get(key).Strings()
}

// Return the value at key as an error, or nil if the value is not an error.
func (self *Map) Err(key string) error {
	return self.Get(key).Err()
}

// Return the value at key as a Map.  If the resulting value is nil or not a
// map type, a null Map will be returned.  All values retrieved from a null
// Map will return that type's zero value.
func (self *Map) Map(key string, tagName ...string) map[typeutil.Variant]typeutil.Variant {
	if len(tagName) == 0 {
		tagName = []string{self.structTagKey}
	}

	return self.Get(key).Map(tagName...)
}

// Return the value as a map[string]interface{} {
func (self *Map) MapNative(tagName ...string) map[string]interface{} {
	if len(tagName) == 0 {
		tagName = []string{self.structTagKey}
	}

	return typeutil.MapNative(self.data, tagName...)
}

func (self *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.data)
}

func xn(generic bool, ifGeneric string, otherwise string) xml.Name {
	if generic {
		return _xn(ifGeneric)
	} else {
		return _xn(otherwise)
	}
}

func _xn(tagName string) xml.Name {
	return xml.Name{
		Local: tagName,
	}
}

func xt(value interface{}) string {
	if typeutil.IsMap(value) {
		return `object`
	} else if typeutil.IsArray(value) {
		return `array`
	} else {
		return strings.TrimPrefix(stringutil.Hyphenate(fmt.Sprintf("%T", value)), `-`)
	}
}

func (self *Map) valueToXmlTokens(parent *xml.StartElement, value interface{}, key string) (tokens []xml.Token, ferr error) {
	g := self.xmlMarshalGeneric

	if typeutil.IsScalar(value) {
		open := xml.StartElement{
			Name: xn(g, `item`, key),
		}

		if g {
			open.Attr = []xml.Attr{
				{
					Name:  _xn(`key`),
					Value: key,
				}, {
					Name:  _xn(`type`),
					Value: utilutil.DetectConvertType(value).String(),
				},
			}
		}

		tokens = append(tokens, open, xml.CharData(typeutil.String(value)), xml.EndElement{
			Name: xn(g, `item`, key),
		})

	} else if typeutil.IsArray(value) {
		start := xml.StartElement{
			Name: xn(g, `item`, key),
		}

		if g {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  _xn(`type`),
				Value: xt(value),
			}, xml.Attr{
				Name:  _xn(`key`),
				Value: key,
			})
		}

		tokens = append(tokens, start)

		for i, v := range sliceutil.Sliceify(value) {
			if ts, err := self.valueToXmlTokens(&start, v, `element`); err == nil {
				tokens = append(tokens, ts...)
			} else {
				ferr = fmt.Errorf("[%d]: %v", i, err)
				return
			}
		}

		tokens = append(tokens, xml.EndElement{
			Name: start.Name,
		})
	} else {
		children := M(value).MapNative()
		ckeys := StringKeys(children)
		sort.Strings(ckeys)

		start := xml.StartElement{
			Name: xn(g, `item`, key),
		}

		if g {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  _xn(`type`),
				Value: xt(value),
			}, xml.Attr{
				Name:  _xn(`key`),
				Value: key,
			})
		}

		tokens = append(tokens, start)

		for _, k := range ckeys {
			v := children[k]

			if ts, err := self.valueToXmlTokens(&start, v, k); err == nil {
				tokens = append(tokens, ts...)
			} else {
				ferr = fmt.Errorf("%s: %v", k, err)
				return
			}
		}

		tokens = append(tokens, xml.EndElement{
			Name: start.Name,
		})
	}

	return
}

func (self *Map) SetMarshalXmlGeneric(yes bool) {
	self.xmlMarshalGeneric = yes
}

// set the name of the root XML tag, used by MarshalXML.
func (self *Map) SetRootTagName(root string) {
	self.rootTagName = root
}

// Marshals the current data into XML.  Nested maps are output as nested elements.  Map values that
// are scalars (strings, numbers, bools, dates/times) will appear as attributes on the parent element.
func (self *Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	root := MapXmlRootTagName

	if self.rootTagName != `` {
		root = self.rootTagName
	}

	start.Name = _xn(root)
	tokens := []xml.Token{start}

	children := self.MapNative()
	ckeys := StringKeys(children)
	sort.Strings(ckeys)

	for _, k := range ckeys {
		v := children[k]

		if ts, err := self.valueToXmlTokens(&start, v, k); err == nil {
			tokens = append(tokens, ts...)
		} else {
			return err
		}
	}

	tokens = append(tokens, xml.EndElement{
		Name: start.Name,
	})

	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}

	return e.Flush()
}

// Return whether the value at the given key is that type's zero value.
func (self *Map) IsZero(key string) bool {
	return self.Get(key).IsZero()
}

// Return the keys in this Map object.  You may specify the name of a struct tag on the underlying
// object to use for generating key names.
func (self *Map) Keys(tagName ...string) []interface{} {
	return Keys(self.MapNative(tagName...))
}

// A string slice version of Keys()
func (self *Map) StringKeys(tagName ...string) []string {
	return sliceutil.Stringify(self.Keys(tagName...))
}

// Return the length of the Map.
func (self *Map) Len() int {
	return len(self.MapNative())
}

// Iterate through each item in the map.
func (self *Map) Each(fn ItemFunc, tagName ...string) error {
	if fn != nil {
		for _, key := range self.StringKeys(tagName...) {
			if err := fn(key, self.Get(key)); err != nil {
				return err
			}
		}
	}

	return nil
}

// A recursive walk form of Each()
func (self *Map) Walk(fn WalkFunc) error {
	return WalkStruct(self.data, fn)
}
