package queryparam

import (
	"reflect"
	"strconv"
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

// StringValueParser parses a string into an int64.
func IntValueParser(value string, _ string, target reflect.Value) error {
	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	target.SetInt(i64)
	return nil
}