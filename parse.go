package queryparam

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

var (
	// ErrNonPointerTarget is returned when the given interface does not represent a pointer
	ErrNonPointerTarget = errors.New("invalid target. must be a non nil pointer")
	// ErrInvalidURLValues is returned when the given *url.URL is nil
	ErrInvalidURLValues = errors.New("invalid url provided")
	// ErrUnhandledFieldType is returned when a struct property is tagged but has an unhandled type.
	ErrUnhandledFieldType = errors.New("unhandled field type")
	// ErrInvalidTag is returned when the tag value is invalid
	ErrInvalidTag = errors.New("invalid tag")
)

// ErrInvalidParameterValue is an error adds extra context to a parser error.
type ErrInvalidParameterValue struct {
	Err       error
	Parameter string
	Field     string
	Value     string
	Type      reflect.Type
}

// Error returns the full error message.
func (e *ErrInvalidParameterValue) Error() string {
	return fmt.Sprintf("invalid parameter value for field %s (%s) from parameter %s (%s): %s", e.Field, e.Type, e.Parameter, e.Value, e.Err.Error())
}

// Unwrap returns the wrapped error.
func (e *ErrInvalidParameterValue) Unwrap() error {
	return e.Err
}

// Present allows you to determine whether or not a query parameter was present in a request.
type Present bool

// DefaultParser is a default parser.
var DefaultParser = &Parser{
	// Tag is the name of the struct tag where the query parameter name is set.
	Tag: "queryparam",
	// Delimiter is the name of the struct tag where a string delimiter override is set.
	DelimiterTag: "queryparamdelim",
	// Delimiter is the default string delimiter.
	Delimiter: ",",
	// ValueParsers is a map[reflect.Type]ValueParser that defines how we parse query
	// parameters based on the destination variable type.
	ValueParsers: DefaultValueParsers(),
}

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

// Parser is used to parse a URL.
type Parser struct {
	Tag          string
	DelimiterTag string
	Delimiter    string
	ValueParsers map[reflect.Type]ValueParser
}

// ValueParser is a func used to parse a value.
type ValueParser func(value string, delimiter string) (reflect.Value, error)

// FieldDelimiter returns a delimiter to be used with the given field.
func (p *Parser) FieldDelimiter(field reflect.StructField) string {
	if customDelimiter := field.Tag.Get(p.DelimiterTag); customDelimiter != "" {
		return customDelimiter
	}
	return p.Delimiter
}

// Parse attempts to parse query parameters from the specified URL and store any found values
// into the given target interface.
func (p *Parser) Parse(urlValues url.Values, target interface{}) error {
	if urlValues == nil {
		return ErrInvalidURLValues
	}
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return ErrNonPointerTarget
	}

	targetElement := targetValue.Elem()
	targetType := targetElement.Type()

	for i := 0; i < targetType.NumField(); i++ {
		if err := p.ParseField(targetType.Field(i), targetElement.Field(i), urlValues); err != nil {
			return err
		}
	}
	return nil
}

// ParseField parses the given field and sets the given value on the target.
func (p *Parser) ParseField(field reflect.StructField, value reflect.Value, urlValues url.Values) error {
	queryParameterName, ok := field.Tag.Lookup(p.Tag)
	if !ok {
		return nil
	}
	if queryParameterName == "" {
		return fmt.Errorf("missing tag value for field: %s: %w", field.Name, ErrInvalidTag)
	}
	queryParameterValue := urlValues.Get(queryParameterName)

	valueParser, ok := p.ValueParsers[field.Type]
	if !ok {
		return fmt.Errorf("%w: %s: %v", ErrUnhandledFieldType, field.Name, field.Type.String())
	}

	parsedValue, err := valueParser(queryParameterValue, p.FieldDelimiter(field))
	if err != nil {
		err = &ErrInvalidParameterValue{
			Err:       err,
			Value:     queryParameterValue,
			Parameter: queryParameterName,
			Type:      field.Type,
			Field:     field.Name,
		}
		return err
	}

	// handle edge case value types
	switch field.Type {
	case reflect.TypeOf(int32(0)):
		value.SetInt(parsedValue.Int())
	case reflect.TypeOf(float32(0)):
		value.SetFloat(parsedValue.Float())
	default:
		value.Set(parsedValue)
	}

	return nil
}

// Parse attempts to parse query parameters from the specified URL and store any found values
// into the given target interface.
func Parse(urlValues url.Values, target interface{}) error {
	return DefaultParser.Parse(urlValues, target)
}
