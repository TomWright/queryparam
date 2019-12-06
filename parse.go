package queryparam

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
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

// ErrCannotSetValue is an error adds extra context to a setter error.
type ErrCannotSetValue struct {
	Err         error
	Parameter   string
	Field       string
	Value       string
	Type        reflect.Type
	ParsedValue reflect.Value
}

// Error returns the full error message.
func (e *ErrCannotSetValue) Error() string {
	return fmt.Sprintf("cannot set value for field %s (%s) from parameter %s (%s - %v): %s", e.Field, e.Type, e.Parameter, e.Value, e.ParsedValue, e.Err.Error())
}

// Unwrap returns the wrapped error.
func (e *ErrCannotSetValue) Unwrap() error {
	return e.Err
}

// Present allows you to determine whether or not a query parameter was present in a request.
type Present bool

// DefaultParser is a default parser.
var DefaultParser = &Parser{
	Tag:          "queryparam",
	DelimiterTag: "queryparamdelim",
	Delimiter:    ",",
	ValueParsers: DefaultValueParsers(),
	ValueSetters: DefaultValueSetters(),
}

// Parser is used to parse a URL.
type Parser struct {
	// Tag is the name of the struct tag where the query parameter name is set.
	Tag string
	// Delimiter is the name of the struct tag where a string delimiter override is set.
	DelimiterTag string
	// Delimiter is the default string delimiter.
	Delimiter string
	// ValueParsers is a map[reflect.Type]ValueParser that defines how we parse query
	// parameters based on the destination variable type.
	ValueParsers map[reflect.Type]ValueParser
	// ValueSetters is a map[reflect.Type]ValueSetter that defines how we set values
	// onto target variables.
	ValueSetters map[reflect.Type]ValueSetter
}

// ValueParser is a func used to parse a value.
type ValueParser func(value string, delimiter string) (reflect.Value, error)

// ValueSetter is a func used to set a value on a target variable.
type ValueSetter func(value reflect.Value, target reflect.Value) error

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
		return &ErrInvalidParameterValue{
			Err:       err,
			Value:     queryParameterValue,
			Parameter: queryParameterName,
			Type:      field.Type,
			Field:     field.Name,
		}
	}

	valueSetter, ok := p.ValueSetters[field.Type]
	if !ok {
		valueSetter, ok = p.ValueSetters[GenericType]
	}
	if !ok {
		return &ErrCannotSetValue{
			Err:         ErrUnhandledFieldType,
			Value:       queryParameterValue,
			ParsedValue: parsedValue,
			Parameter:   queryParameterName,
			Type:        field.Type,
			Field:       field.Name,
		}
	}

	if err := valueSetter(parsedValue, value); err != nil {
		return &ErrCannotSetValue{
			Err:         err,
			Value:       queryParameterValue,
			ParsedValue: parsedValue,
			Parameter:   queryParameterName,
			Type:        field.Type,
			Field:       field.Name,
		}
	}

	return nil
}

// Parse attempts to parse query parameters from the specified URL and store any found values
// into the given target interface.
func Parse(urlValues url.Values, target interface{}) error {
	return DefaultParser.Parse(urlValues, target)
}
