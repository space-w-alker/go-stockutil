package typeutil

import (
	"fmt"
	"reflect"
)

// Returns the number of input and return arguments a given function has.
func FunctionArity(fn interface{}) (int, int, error) {
	if IsFunction(fn) {
		fnT := reflect.TypeOf(fn)

		return fnT.NumIn(), fnT.NumOut(), nil
	} else {
		return 0, 0, fmt.Errorf("expected function, got %T", fn)
	}
}
