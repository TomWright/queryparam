package queryparam

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidBoolValue = errors.New("invalid bool value")

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

// IntValueParser parses a string into an int64.
func IntValueParser(value string, _ string, target reflect.Value) error {
	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	target.SetInt(i64)
	return nil
}

// TimeValueParser parses a string into an int64.
func TimeValueParser(value string, _ string, target reflect.Value) error {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	target.Set(reflect.ValueOf(t))
	return nil
}

// BoolValueParser parses a string into a bool.
func BoolValueParser(value string, _ string, target reflect.Value) error {
	var x bool
	switch strings.ToLower(value) {
	case "true", "1", "y", "yes":
		x = true
	case "false", "0", "n", "no":
		x = false
	default:
		return ErrInvalidBoolValue
	}
	target.SetBool(x)
	return nil
}
