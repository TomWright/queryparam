package queryparam_test

import (
	"errors"
	"fmt"
	"github.com/tomwright/queryparam/v4"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

var urlValues = url.Values{}
var urlValuesNameAge = url.Values{
	"name": []string{"tom"},
	"age":  []string{"26"},
}

// ExampleParse creates a dummy http request and parses the data into a struct
func ExampleParse() {
	var err error
	var request = &http.Request{}
	request.URL, err = url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		panic(err)
	}

	requestData := struct {
		Name string `queryparam:"name"`
		Age  int    `queryparam:"age"`
	}{}

	err = queryparam.Parse(request.URL.Query(), &requestData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s is %d", requestData.Name, requestData.Age)

	// Output: Tom is 23
}

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	t.Run("ValueParser", func(t *testing.T) {
		t.Parallel()

		calledCount := 0
		p := &queryparam.Parser{
			Tag:          "queryparam",
			DelimiterTag: "queryparamdelim",
			Delimiter:    ",",
			ValueParsers: map[reflect.Type]queryparam.ValueParser{
				reflect.TypeOf(""): func(value string, delimiter string) (reflect.Value, error) {
					calledCount++
					if value != "Tom" {
						t.Errorf("unexpected value: %s", value)
					}
					if delimiter != "," {
						t.Errorf("unexpected delimiter: %s", delimiter)
					}
					return reflect.ValueOf(""), nil
				},
			},
		}

		urlValues := url.Values{}
		urlValues.Set("name", "Tom")

		if err := p.Parse(urlValues, &struct {
			Name string `queryparam:"name"`
		}{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if calledCount != 1 {
			t.Errorf("unexpected called count: %d", calledCount)
		}
	})

	t.Run("ValueParserCustomDelimiter", func(t *testing.T) {
		t.Parallel()

		calledCount := 0
		p := &queryparam.Parser{
			Tag:          "queryparam",
			DelimiterTag: "queryparamdelim",
			Delimiter:    ",",
			ValueParsers: map[reflect.Type]queryparam.ValueParser{
				reflect.TypeOf(""): func(value string, delimiter string) (reflect.Value, error) {
					calledCount++
					if value != "Tom" {
						t.Errorf("unexpected value: %s", value)
					}
					if delimiter != ":" {
						t.Errorf("unexpected delimiter: %s", delimiter)
					}
					return reflect.ValueOf(""), nil
				},
			},
		}

		urlValues := url.Values{}
		urlValues.Set("name", "Tom")

		if err := p.Parse(urlValues, &struct {
			Name string `queryparam:"name" queryparamdelim:":"`
		}{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if calledCount != 1 {
			t.Errorf("unexpected called count: %d", calledCount)
		}
	})

	t.Run("Present", func(t *testing.T) {
		t.Parallel()

		p := &queryparam.Parser{
			Tag:          "queryparam",
			DelimiterTag: "queryparamdelim",
			Delimiter:    ",",
			ValueParsers: queryparam.DefaultValueParsers(),
		}

		urlValues := url.Values{}
		urlValues.Set("first-name", "Tom")
		urlValues.Set("last-name", "")
		urlValues.Set("age", "26")

		type testData struct {
			FirstName queryparam.Present `queryparam:"first-name"`
			LastName  queryparam.Present `queryparam:"last-name"`
			Age       queryparam.Present `queryparam:"age"`
		}
		var res testData

		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		exp := testData{
			FirstName: true,
			LastName:  false,
			Age:       true,
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})
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
