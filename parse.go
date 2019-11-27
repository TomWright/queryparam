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
	// ErrInvalidURL is returned when the given *url.URL is nil
	ErrInvalidURL = errors.New("invalid url provided")
	// ErrUnhandledFieldType is returned when a struct property is tagged but has an unhandled type.
	ErrUnhandledFieldType = errors.New("invalid url provided")
	// ErrInvalidTag is returned when the tag value is invalid
	ErrInvalidTag = errors.New("invalid tag")
)

// DefaultParser is a default parser.
var DefaultParser = &Parser{
	Tag: "queryparam",
	DelimiterTag: "queryparamdelim",
	Delimiter: ",",
	ValueParsers: map[reflect.Type]ValueParser{
		reflect.TypeOf(""): StringValueParser,
		reflect.TypeOf([]string{}): StringSliceValueParser,
		reflect.TypeOf(0): IntValueParser,
		reflect.TypeOf(int32(0)): IntValueParser,
		reflect.TypeOf(int64(0)): IntValueParser,
	},
}

// Parser is used to parse a URL.
type Parser struct {
	Tag          string
	DelimiterTag string
	Delimiter    string
	ValueParsers map[reflect.Type]ValueParser
}

// ValueParser is a func used to parse a value and set it on a target.
type ValueParser func(value string, delimiter string, target reflect.Value) error

// FieldDelimiter returns a delimiter to be used with the given field.
func (p *Parser) FieldDelimiter(field reflect.StructField) string {
	if customDelimiter := field.Tag.Get(p.DelimiterTag); customDelimiter != "" {
		return customDelimiter
	}
	return p.Delimiter
}

// Parse attempts to parse query parameters from the specified URL and store any found values
// into the given target interface.
func (p *Parser) Parse(u *url.URL, target interface{}) error {
	if u == nil {
		return ErrInvalidURL
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return ErrNonPointerTarget
	}

	targetElement := targetValue.Elem()
	targetType := targetElement.Type()

	urlValues := u.Query()
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
		return fmt.Errorf("%w: missing value", ErrInvalidTag)
	}
	queryParameterValue := urlValues.Get(queryParameterName)
	if queryParameterValue == "" {
		return nil
	}

	valueParser, ok := p.ValueParsers[field.Type]
	if !ok {
		return fmt.Errorf("%w: %s: %v", ErrUnhandledFieldType, field.Name, field.Type.String())
	}

	if err := valueParser(queryParameterValue, p.FieldDelimiter(field), value); err != nil {
		return err
	}

	return nil
}

// Parse attempts to parse query parameters from the specified URL and store any found values
// into the given target interface.
func Parse(u *url.URL, target interface{}) error {
	return DefaultParser.Parse(u, target)
}