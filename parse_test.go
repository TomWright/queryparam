package queryparam_test

import (
	"errors"
	"fmt"
	"github.com/tomwright/queryparam/v4"
	"math"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var urlValues = url.Values{}
var urlValuesNameAge = url.Values{
	"name": []string{"tom"},
	"age":  []string{"26"},
}

// ExampleParse creates a dummy http request and parses the data into a struct
func ExampleParse() {
	urlValues := url.Values{}
	urlValues.Set("name", "Tom")
	urlValues.Set("names", "Tom,Jim,Frank")      // list of names separated by ,
	urlValues.Set("dash-names", "Tom-Jim-Frank") // list of names separated by -
	urlValues.Set("age", "123")
	urlValues.Set("age32", "123")
	urlValues.Set("age64", "123")
	urlValues.Set("float32", "123.45")
	urlValues.Set("float64", "123.45")
	urlValues.Set("created-at", "2019-02-05T13:32:02Z")
	urlValues.Set("bool-false", "false")
	urlValues.Set("bool-true", "true")
	urlValues.Set("bool-empty", "")    // param set but no value means
	urlValues.Set("present-empty", "") // param set but no value
	urlValues.Set("present", "this is here")

	requestData := struct {
		Name         string             `queryparam:"name"`
		Names        []string           `queryparam:"names"`
		DashNames    []string           `queryparam:"dash-names" queryparamdelim:"-"`
		Age          int                `queryparam:"age"`
		Age32        int32              `queryparam:"age32"`
		Age64        int64              `queryparam:"age64"`
		Float32      float32            `queryparam:"float32"`
		Float64      float64            `queryparam:"float64"`
		CreatedAt    time.Time          `queryparam:"created-at"`
		UpdatedAt    time.Time          `queryparam:"updated-at"`
		BoolFalse    bool               `queryparam:"bool-false"`
		BoolTrue     bool               `queryparam:"bool-true"`
		BoolEmpty    bool               `queryparam:"bool-empty"`
		PresentEmpty queryparam.Present `queryparam:"present-empty"`
		Present      queryparam.Present `queryparam:"present"`
		NotPresent   queryparam.Present `queryparam:"not-present"`
	}{}

	if err := queryparam.Parse(urlValues, &requestData); err != nil {
		panic(err)
	}

	fmt.Printf("name: %s\n", requestData.Name)
	fmt.Printf("names: %v\n", requestData.Names)
	fmt.Printf("dash names: %v\n", requestData.DashNames)
	fmt.Printf("age: %d\n", requestData.Age)
	fmt.Printf("age32: %d\n", requestData.Age32)
	fmt.Printf("age64: %d\n", requestData.Age64)
	fmt.Printf("float32: %f\n", math.Round(float64(requestData.Float32))) // rounded to avoid floating point precision issues
	fmt.Printf("float64: %f\n", math.Round(requestData.Float64))          // rounded to avoid floating point precision issues
	fmt.Printf("created at: %s\n", requestData.CreatedAt.Format(time.RFC3339))
	fmt.Printf("updated at: %s\n", requestData.UpdatedAt.Format(time.RFC3339))
	fmt.Printf("bool false: %v\n", requestData.BoolFalse)
	fmt.Printf("bool true: %v\n", requestData.BoolTrue)
	fmt.Printf("bool empty: %v\n", requestData.BoolEmpty)
	fmt.Printf("present empty: %v\n", requestData.PresentEmpty)
	fmt.Printf("present: %v\n", requestData.Present)
	fmt.Printf("not present: %v\n", requestData.NotPresent)

	// Output:
	// name: Tom
	// names: [Tom Jim Frank]
	// dash names: [Tom Jim Frank]
	// age: 123
	// age32: 123
	// age64: 123
	// float32: 123.000000
	// float64: 123.000000
	// created at: 2019-02-05T13:32:02Z
	// updated at: 0001-01-01T00:00:00Z
	// bool false: false
	// bool true: true
	// bool empty: false
	// present empty: false
	// present: true
	// not present: false
}

func TestParse_FieldWithNoTagIsNotUsed(t *testing.T) {
	t.Parallel()

	req := &struct {
		Name string ``
	}{}

	if err := queryparam.Parse(urlValues, req); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if exp, got := "", req.Name; exp != got {
		t.Errorf("unexpected name. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_InvalidURLValues(t *testing.T) {
	t.Parallel()

	req := &struct{}{}

	err := queryparam.Parse(nil, req)
	if exp, got := queryparam.ErrInvalidURLValues, err; exp != got {
		t.Errorf("unexpected error. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_NonPointerTarget(t *testing.T) {
	t.Parallel()

	req := struct{}{}

	err := queryparam.Parse(urlValues, req)
	if exp, got := queryparam.ErrNonPointerTarget, err; exp != got {
		t.Errorf("unexpected error. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_UnhandledFieldType(t *testing.T) {
	t.Parallel()

	req := &struct {
		Age struct{} `queryparam:"age"`
	}{}

	err := queryparam.Parse(urlValuesNameAge, req)
	if !errors.Is(err, queryparam.ErrUnhandledFieldType) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_EmptyTag(t *testing.T) {
	t.Parallel()

	req := &struct {
		Name string `queryparam:""`
	}{}

	err := queryparam.Parse(urlValuesNameAge, req)
	if !errors.Is(err, queryparam.ErrInvalidTag) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_ValueParserErrorReturned(t *testing.T) {
	t.Parallel()

	tmpErr := errors.New("something bad happened")

	p := &queryparam.Parser{
		Tag:          "queryparam",
		DelimiterTag: "queryparamdelim",
		Delimiter:    ",",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf(""): func(value string, delimiter string) (reflect.Value, error) {
				return reflect.ValueOf(""), tmpErr
			},
		},
	}

	req := &struct {
		Name string `queryparam:"name"`
	}{}

	err := p.Parse(urlValuesNameAge, req)
	if !errors.Is(err, tmpErr) {
		t.Errorf("unexpected error: %v", err)
	}
}

func BenchmarkParse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := struct {
			Name         string   `queryparam:"name"`
			NameList     []string `queryparam:"name-list"`
			NameListDash []string `queryparam:"name-list" queryparamdelim:"-"`
			Age          int      `queryparam:"age"`
		}{}
		err := queryparam.Parse(urlValuesNameAge, &data)
		if err != nil {
			b.FailNow()
		}
	}
	b.StopTimer()
	b.ReportAllocs()
}

func TestErrInvalidParameterValue_Unwrap(t *testing.T) {
	e := &queryparam.ErrInvalidParameterValue{
		Err:       errors.New("something bad"),
		Parameter: "Name",
		Field:     "name",
		Value:     "asd",
		Type:      reflect.TypeOf(""),
	}
	exp := "invalid parameter value for field name (string) from parameter Name (asd): something bad"
	if got := e.Error(); exp != got {
		t.Errorf("expected `%s`, got `%s`", exp, got)
	}
}
