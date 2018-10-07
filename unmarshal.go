package queryparam

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	tag              = "queryparam"
	delimiterTag     = "queryparamdelim"
	defaultDelimiter = ","
)

var (
	// ErrNonPointerTarget is returned when the given interface does not represent a pointer
	ErrNonPointerTarget = errors.New("invalid Unmarshal target. must be a pointer")
	// ErrInvalidURL is returned when the given *url.URL is nil
	ErrInvalidURL = errors.New("invalid url provided")
)

func delimiterFromField(field reflect.StructField) string {
	if customDelimiter := field.Tag.Get(delimiterTag); customDelimiter != "" {
		return customDelimiter
	}
	return defaultDelimiter
}

func unmarshalField(v reflect.Value, t reflect.Type, i int, qs url.Values) error {
	var (
		field    reflect.StructField
		paramVal string
		tagVal   string
	)

	field = t.Field(i)

	tagVal = field.Tag.Get(tag)
	if tagVal == "" {
		return nil
	}

	paramVal = qs.Get(tagVal)
	if len(paramVal) == 0 {
		return nil
	}

	switch {
	case field.Type.Kind() == reflect.String:
		v.Field(i).SetString(paramVal)
	case field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String:
		vField := v.Field(i)
		vField.Set(reflect.AppendSlice(vField, reflect.ValueOf(strings.Split(paramVal, delimiterFromField(field)))))
	default:
		return fmt.Errorf("invalid field type. `%s` must be `string` or `[]string`", field.Name)
	}

	return nil
}

// Unmarshal attempts to parse query parameters from the specified URL and store any found values
// into the given interface
func Unmarshal(u *url.URL, i interface{}) error {
	if u == nil {
		return ErrInvalidURL
	}

	iVal := reflect.ValueOf(i)
	if iVal.Kind() != reflect.Ptr || iVal.IsNil() {
		return ErrNonPointerTarget
	}

	v := iVal.Elem()
	t := v.Type()

	qs := u.Query()
	for i := 0; i < t.NumField(); i++ {
		if err := unmarshalField(v, t, i, qs); err != nil {
			return err
		}
	}
	return nil
}
