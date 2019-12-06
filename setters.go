package queryparam

import (
	"fmt"
	"reflect"
)

// GenericType is a reflect.Type we can use to identity generic value setters.
var GenericType = reflect.TypeOf(struct{}{})

// DefaultValueSetters returns a set of default value setters.
func DefaultValueSetters() map[reflect.Type]ValueSetter {
	return map[reflect.Type]ValueSetter{
		GenericType:                GenericSetter,
		reflect.TypeOf(int32(0)):   Int32ValueSetter,
		reflect.TypeOf(float32(0)): Float32ValueSetter,
	}
}

// recoverPanic recovers from a panic and sets the panic value into the given error.
func recoverPanic(err *error) func() {
	return func() {
		rec := recover()
		if rec != nil {
			if e, ok := rec.(error); ok {
				*err = e
			} else {
				*err = fmt.Errorf("%s", rec)
			}
		}
	}
}

// GenericSetter sets the targets value.
func GenericSetter(value reflect.Value, target reflect.Value) (err error) {
	defer recoverPanic(&err)()
	target.Set(value)
	return
}

// Int32ValueSetter sets the targets value to an int32.
func Int32ValueSetter(value reflect.Value, target reflect.Value) (err error) {
	defer recoverPanic(&err)()
	target.SetInt(value.Int())
	return nil
}

// Float32ValueSetter sets the targets value to a float32.
func Float32ValueSetter(value reflect.Value, target reflect.Value) (err error) {
	defer recoverPanic(&err)()
	target.SetFloat(value.Float())
	return nil
}
