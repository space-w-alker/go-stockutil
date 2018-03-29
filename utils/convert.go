package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ConvertType int

const (
	Invalid ConvertType = iota
	String
	Boolean
	Float
	Integer
	Time
)

var NilStrings = []string{`null`, `NULL`, `<nil>`, `nil`, `Nil`, `None`, `undefined`, ``}
var BooleanTrueValues = []string{`true`, `yes`, `on`}
var BooleanFalseValues = []string{`false`, `no`, `off`}

var TimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC1123,
	time.RFC1123Z,
	time.Kitchen,
	`2006-01-02 15:04:05.000000000`,
	`2006-01-02 15:04:05.000000`,
	`2006-01-02 15:04:05.000`,
	`2006-01-02 15:04:05 -0700 MST`,
	`2006-01-02 15:04:05Z07:00`,
	`2006-01-02 15:04:05`,
	`2006-01-02 15:04`,
	`2006-01-02`,
	`2006-01-02T15:04:05.000000000`,
	`2006-01-02T15:04:05.000000`,
	`2006-01-02T15:04:05.000`,
	`2006-01-02T15:04:05 -0700 MST`,
	`2006-01-02T15:04:05Z07:00`,
	`2006-01-02T15:04:05`,
}

func ToString(in interface{}) (string, error) {
	if in == nil {
		return ``, nil
	}

	if inStr, ok := in.(fmt.Stringer); ok {
		return inStr.String(), nil
	}

	if inT := reflect.TypeOf(in); inT != nil {
		switch inT.Kind() {
		case reflect.Float32:
			return strconv.FormatFloat(float64(in.(float32)), 'f', -1, 32), nil
		case reflect.Float64:
			return strconv.FormatFloat(in.(float64), 'f', -1, 64), nil
		case reflect.Bool:
			return strconv.FormatBool(in.(bool)), nil
		case reflect.String:
			if inStr, ok := in.(string); ok {
				return inStr, nil
			}
		}

		if !IsKind(in, CompoundTypes...) {
			return fmt.Sprintf("%v", in), nil
		}
	}

	return ``, fmt.Errorf("Unable to convert type '%T' to string", in)
}

func ConvertTo(toType ConvertType, inI interface{}) (interface{}, error) {
	if in, err := ToString(inI); err == nil {
		switch toType {
		case Float:
			if in == `` {
				return float64(0), nil
			}

			return strconv.ParseFloat(in, 64)
		case Integer:
			if inTime, ok := inI.(time.Time); ok {
				return int64(inTime.UnixNano()), nil
			} else if in == `` {
				return int64(0), nil
			} else if layout := DetectTimeFormat(in); layout != `` && layout != `epoch` {
				if tm, err := time.Parse(layout, in); err == nil {
					return tm.UnixNano(), nil
				} else {
					return nil, err
				}
			}

			return strconv.ParseInt(in, 10, 64)
		case Boolean:
			if inI == nil {
				return false, nil
			}

			if IsBooleanTrue(in) {
				return true, nil
			} else if IsBooleanFalse(in) {
				return false, nil
			} else {
				return nil, fmt.Errorf("Cannot convert '%s' into a boolean value", in)
			}
		case Time:
			if inTime, ok := inI.(time.Time); ok {
				return inTime, nil
			}

			in := strings.Trim(strings.TrimSpace(in), `"'`)

			if DetectTimeFormat(in) == `epoch` {
				if v, err := strconv.ParseInt(in, 10, 64); err == nil {
					return time.Unix(v, 0), nil
				}
			}

			for _, format := range TimeFormats {
				if tm, err := time.Parse(format, strings.TrimSpace(in)); err == nil {
					return tm, nil
				}
			}

			switch in {
			case `now`:
				return time.Now(), nil
			default:
				// handle time zero values
				tmS := strings.Map(func(r rune) rune {
					switch r {
					case '-', ':', ' ', 'T', 'Z':
						return '0'
					}

					return r
				}, in)

				if v, err := strconv.ParseInt(tmS, 10, 64); err == nil && v == 0 {
					return time.Time{}, nil
				}

				return nil, fmt.Errorf("Cannot convert '%s' into a date/time value", in)
			}

		default:
			return in, nil
		}
	} else {
		return nil, err
	}
}

func ConvertToInteger(in interface{}) (int64, error) {
	if v, err := ConvertTo(Integer, in); err == nil {
		return v.(int64), nil
	} else {
		return int64(0), err
	}
}

func ConvertToFloat(in interface{}) (float64, error) {
	if v, err := ConvertTo(Float, in); err == nil {
		return v.(float64), nil
	} else {
		return float64(0.0), err
	}
}

func ConvertToString(in interface{}) (string, error) {
	if v, err := ConvertTo(String, in); err == nil {
		return v.(string), nil
	} else {
		return ``, err
	}
}

func ConvertToBool(in interface{}) (bool, error) {
	if v, err := ConvertTo(Boolean, in); err == nil {
		return v.(bool), nil
	} else {
		return false, err
	}
}

func ConvertToTime(in interface{}) (time.Time, error) {
	switch in.(type) {
	case time.Time:
		return in.(time.Time), nil
	default:
		if v, err := ConvertTo(Time, in); err == nil {
			return v.(time.Time), nil
		} else {
			return time.Time{}, err
		}
	}
}

func DetectTimeFormat(in string) string {
	if IsInteger(in) {
		return `epoch`
	}

	for _, layout := range TimeFormats {
		if _, err := time.Parse(layout, in); err == nil {
			return layout
		}
	}

	return ``
}
