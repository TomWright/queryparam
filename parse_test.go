package queryparam_test

import (
	"errors"
	"fmt"
	"github.com/tomwright/queryparam/v3"
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
				reflect.TypeOf(""): func(value string, delimiter string, target reflect.Value) error {
					calledCount++
					if value != "Tom" {
						t.Errorf("unexpected value: %s", value)
					}
					if delimiter != "," {
						t.Errorf("unexpected delimiter: %s", delimiter)
					}
					return nil
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
				reflect.TypeOf(""): func(value string, delimiter string, target reflect.Value) error {
					calledCount++
					if value != "Tom" {
						t.Errorf("unexpected value: %s", value)
					}
					if delimiter != ":" {
						t.Errorf("unexpected delimiter: %s", delimiter)
					}
					return nil
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

	t.Run("SkipBlankQueryParameters", func(t *testing.T) {
		t.Parallel()

		p := &queryparam.Parser{
			Tag:          "queryparam",
			DelimiterTag: "queryparamdelim",
			Delimiter:    ",",
			ValueParsers: map[reflect.Type]queryparam.ValueParser{},
		}

		urlValues := url.Values{}
		urlValues.Set("name", "")

		if err := p.Parse(urlValues, &struct {
			Name string `queryparam:"name"`
		}{}); err != nil {
			t.Errorf("unexpected error: %v", err)
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
			reflect.TypeOf(""): func(value string, delimiter string, target reflect.Value) error {
				return tmpErr
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
