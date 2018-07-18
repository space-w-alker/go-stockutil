package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	multierror "github.com/hashicorp/go-multierror"
)

func GenericMarshalJSON(in interface{}, extraData ...map[string]interface{}) ([]byte, error) {
	sval := structs.New(in)
	output := make(map[string]interface{})

FieldLoop:
	for _, field := range sval.Fields() {
		if field.IsExported() {
			key := field.Name()
			value := field.Value()

			// perform type conversions because encoding/json won't Do The Right Thing in all cases.
			if vTime, ok := value.(time.Time); ok {
				if vTime.IsZero() {
					value = nil
				} else {
					value = vTime.Format(time.RFC3339Nano)
				}
			} else if _, ok := value.(fmt.Stringer); ok {
				value = fmt.Sprintf("%v", value)
			}

			// parse "json" tag
			if parts := strings.Split(field.Tag(`json`), `,`); len(parts) > 0 {
				if parts[0] != `` {
					key = parts[0]
				}

				if len(parts) > 1 {
					for _, flag := range parts[1:] {
						switch flag {
						case `omitempty`:
							if field.IsZero() {
								continue FieldLoop
							}
						}
					}
				}
			}

			output[key] = value
		}
	}

	if len(extraData) > 0 && len(extraData[0]) > 0 {
		for k, v := range extraData[0] {
			output[k] = v
		}
	}

	return json.Marshal(output)
}

// Appends on error to another, allowing for operations that return multiple errors
// to remain compatible within a single-valued context.
func AppendError(base error, err error) error {
	if err == nil {
		return base
	} else {
		return multierror.Append(base, err)
	}
}
