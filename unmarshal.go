package queryparam

import (
	"reflect"
	"errors"
	"net/url"
	"fmt"
	"strings"
)

const (
	tag              = "queryparam"
	delimiterTag     = "queryparamdelim"
	defaultDelimiter = ","
)

var (
	stringType      = reflect.TypeOf("")
	stringSliceType = reflect.TypeOf(make([]string, 0))

	// ErrNonPointerTarget is returned when the given interface does not represent a pointer
	ErrNonPointerTarget = errors.New("invalid Unmarshal target. must be a pointer")
	// ErrInvalidURL is returned when the given *url.URL is nil
	ErrInvalidURL = errors.New("invalid url provided")
	// ErrInvalidDelimiter is returned when trying to split a query param into a slice with an invalid separator
	ErrInvalidDelimiter = errors.New("invalid query param separator")
)

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
	var paramVal, tagVal string
	var field reflect.StructField
	var vField reflect.Value

	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)

		tagVal = field.Tag.Get(tag)
		if tagVal != "" {
			paramVal = u.Query().Get(tagVal)

			switch field.Type {
			case stringType:
				v.Field(i).SetString(paramVal)
			case stringSliceType:
				delimiter := defaultDelimiter
				if customDelimiter := field.Tag.Get(delimiterTag); customDelimiter != "" {
					delimiter = customDelimiter
				}
				if len(delimiter) == 0 {
					return ErrInvalidDelimiter
				}
				if len(paramVal) == 0 {
					continue
				}
				vField = v.Field(i)
				vField.Set(reflect.AppendSlice(vField, reflect.ValueOf(strings.Split(paramVal, delimiter))))
			default:
				return fmt.Errorf("invalid field type. `%s` must be `string` or `[]string`", field.Name)
			}
		}
	}
	return nil
}
