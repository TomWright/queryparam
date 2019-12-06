package queryparam

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidBoolValue is returned when an unhandled string is parsed.
var ErrInvalidBoolValue = errors.New("unknown bool value")

// DefaultValueParsers returns a set of default value parsers.
func DefaultValueParsers() map[reflect.Type]ValueParser {
	return map[reflect.Type]ValueParser{
		reflect.TypeOf(""):             StringValueParser,
		reflect.TypeOf([]string{}):     StringSliceValueParser,
		reflect.TypeOf(0):              IntValueParser,
		reflect.TypeOf(int32(0)):       Int32ValueParser,
		reflect.TypeOf(int64(0)):       Int64ValueParser,
		reflect.TypeOf(float32(0)):     Float32ValueParser,
		reflect.TypeOf(float64(0)):     Float64ValueParser,
		reflect.TypeOf(time.Time{}):    TimeValueParser,
		reflect.TypeOf(false):          BoolValueParser,
		reflect.TypeOf(Present(false)): PresentValueParser,
	}
}

// StringValueParser parses a string into a string.
func StringValueParser(value string, _ string) (reflect.Value, error) {
	return reflect.ValueOf(value), nil
}

// StringSliceValueParser parses a string into a []string.
func StringSliceValueParser(value string, delimiter string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf([]string{}), nil
	}

	return reflect.ValueOf(strings.Split(value, delimiter)), nil
}

// IntValueParser parses a string into an int64.
func IntValueParser(value string, _ string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf(int(0)), nil
	}

	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return reflect.ValueOf(int(0)), err
	}

	return reflect.ValueOf(int(i64)), nil
}

// Int64ValueParser parses a string into an int64.
func Int64ValueParser(value string, _ string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf(int64(0)), nil
	}

	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return reflect.ValueOf(int64(0)), err
	}

	return reflect.ValueOf(i64), nil
}

// Int32ValueParser parses a string into an int32.
func Int32ValueParser(value string, _ string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf(0), nil
	}

	i64, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return reflect.ValueOf(0), err
	}

	return reflect.ValueOf(i64), nil
}

// TimeValueParser parses a string into an int64.
func TimeValueParser(value string, _ string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf(time.Time{}), nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return reflect.ValueOf(time.Time{}), err
	}
	return reflect.ValueOf(t), nil
}

// BoolValueParser parses a string into a bool.
func BoolValueParser(value string, _ string) (reflect.Value, error) {
	switch strings.ToLower(value) {
	case "true", "1", "y", "yes":
		return reflect.ValueOf(true), nil
	case "", "false", "0", "n", "no":
		return reflect.ValueOf(false), nil
	default:
		return reflect.ValueOf(false), ErrInvalidBoolValue
	}
}

// PresentValueParser sets the target to true.
// This parser will only be executed if the parameter is present.
func PresentValueParser(value string, _ string) (reflect.Value, error) {
	return reflect.ValueOf(Present(value != "")), nil
}

// Float64ValueParser parses a string to a float64.
func Float64ValueParser(value string, _ string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf(float64(0)), nil
	}

	i64, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return reflect.ValueOf(float64(0)), err
	}

	return reflect.ValueOf(i64), nil
}

// Float32ValueParser parses a string to a float32.
func Float32ValueParser(value string, _ string) (reflect.Value, error) {
	if value == "" {
		// ignore blank values.
		return reflect.ValueOf(float32(0)), nil
	}

	f64, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return reflect.ValueOf(float32(0)), err
	}

	return reflect.ValueOf(f64), nil
}
