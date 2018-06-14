package queryparam

import (
	"reflect"
	"errors"
	"net/url"
	"fmt"
)

var (
	tag        = "queryparam"
	stringType = reflect.TypeOf("")

	// ErrNonPointerTarget is returned when the given interface does not represent a pointer
	ErrNonPointerTarget = errors.New("invalid unmarshall target. must be a pointer")
	// ErrInvalidURL is returned when the given *url.URL is nil
	ErrInvalidURL = errors.New("invalid url provided")
)

// Unmarshall attempts to parse query parameters from the specified URL and store any found values
// into the given interface
func Unmarshall(u *url.URL, i interface{}) error {
	if u == nil {
		return ErrInvalidURL
	}

	iVal := reflect.ValueOf(i)
	if iVal.Kind() != reflect.Ptr || iVal.IsNil() {
		return ErrNonPointerTarget
	}

	v := iVal.Elem()

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tagVal := field.Tag.Get(tag)
		if tagVal != "" {
			if field.Type != stringType {
				return fmt.Errorf("invalid field type. `%s` must be a string", field.Name)
			}

			v.Field(i).SetString(u.Query().Get(tagVal))
		}
	}
	return nil
}
