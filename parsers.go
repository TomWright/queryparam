package queryparam

import (
	"reflect"
	"strings"
)

// StringValueParser parses a string into a string.
func StringValueParser(value string, _ string, target reflect.Value) error {
	target.SetString(value)
	return nil
}

// StringSliceValueParser parses a string into a []string.
func StringSliceValueParser(value string, delimiter string, target reflect.Value) error {
	target.Set(reflect.AppendSlice(target, reflect.ValueOf(strings.Split(value, delimiter))))
	return nil
}