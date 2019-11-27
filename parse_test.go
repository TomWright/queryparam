package queryparam_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/tomwright/queryparam/v3"
)

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
		Age  string `queryparam:"age"`
	}{}

	err = queryparam.Parse(request.URL, &requestData)
	if err != nil {
		panic(err)
	}

	fmt.Println(requestData.Name + " is " + requestData.Age)

	// Output: Tom is 23
}

type testData struct {
	Name         string   `queryparam:"name"`
	NameList     []string `queryparam:"names"`
	NameListDash []string `queryparam:"names" queryparamdelim:"-"`
	Age          int   `queryparam:"age"`
}

func TestParser_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		TestDesc   string
		URL        string
		Data       testData
		OutputData testData
	}{
		{
			TestDesc: "SingleStrings",
			URL:      "https://example.com/some/path?name=Tom&age=23",
			Data:     testData{},
			OutputData: testData{
				Name: "Tom",
				Age:  23,
			},
		},
		{
			TestDesc: "CommaDelimitedStrings",
			URL:      "https://example.com/some/path?name=Tom&names=x,y,z&age=23",
			Data:     testData{},
			OutputData: testData{
				Name:         "Tom",
				NameList:     []string{"x", "y", "z"},
				NameListDash: []string{"x,y,z"},
				Age:          23,
			},
		},
		{
			TestDesc: "CommaAndDashDelimitedStrings",
			URL:      "https://example.com/some/path?names=x,y-z&age=23",
			Data:     testData{},
			OutputData: testData{
				NameList:     []string{"x", "y-z"},
				NameListDash: []string{"x,y", "z"},
				Age:          23,
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.TestDesc, func(t *testing.T) {
			t.Parallel()

			u, err := url.Parse(tc.URL)
			if err != nil {
				t.Errorf("unable to parse url: %s", err)
				return
			}

			err = queryparam.Parse(u, &tc.Data)
			if err != nil {
				t.Errorf("unable to parse url: %s", err)
				return
			}

			if !reflect.DeepEqual(tc.Data, tc.OutputData) {
				t.Errorf("expected `%v`, got `%v`", tc.OutputData, tc.Data)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name []string `queryparam:"name" queryparamdelim:""`
	}{}

	err = queryparam.Parse(u, req)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if exp, got := []string{"Tom", "Jim"}, req.Name; !reflect.DeepEqual(exp, got) {
		t.Errorf("unexpected result. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_UnusedInvalidField(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23&gender=male")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name   string `queryparam:"name"`
		Age    int    ``
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Parse(u, req)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if exp, got := "Tom", req.Name; exp != got {
		t.Errorf("unexpected name. expected `%v`, got `%v`", exp, got)
	}
	if exp, got := 0, req.Age; exp != got {
		t.Errorf("unexpected age. expected `%v`, got `%v`", exp, got)
	}
	if exp, got := "male", req.Gender; exp != got {
		t.Errorf("unexpected gender. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_InvalidURL(t *testing.T) {
	t.Parallel()

	req := &struct {}{}

	err := queryparam.Parse(nil, req)
	if exp, got := queryparam.ErrInvalidURL, err; exp != got {
		t.Errorf("unexpected error. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_NonPointerTarget(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := struct {}{}

	err = queryparam.Parse(u, req)
	if exp, got := queryparam.ErrNonPointerTarget, err; exp != got {
		t.Errorf("unexpected error. expected `%v`, got `%v`", exp, got)
	}
}

func TestParse_UnhandledFieldType(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name   string `queryparam:"name"`
		Age    struct{}    `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Parse(u, req)
	if !errors.Is(err, queryparam.ErrUnhandledFieldType) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_EmptyTag(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name   string `queryparam:""`
	}{}

	err = queryparam.Parse(u, req)
	if !errors.Is(err, queryparam.ErrInvalidTag) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestIntValueParser_InvalidInt(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?age=asd")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Age   int `queryparam:"age"`
	}{}

	err = queryparam.Parse(u, req)
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if numErr.Err != strconv.ErrSyntax {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_ValueParserErrorReturned(t *testing.T) {
	t.Parallel()

	tmpErr := errors.New("something bad happened")

	p := &queryparam.Parser{
		Tag: "queryparam",
		DelimiterTag: "queryparamdelim",
		Delimiter: ",",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf(""): func(value string, delimiter string, target reflect.Value) error {
				return tmpErr
			},
		},
	}

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name   string `queryparam:"name"`
	}{}

	err = p.Parse(u, req)
	if !errors.Is(err, tmpErr) {
		t.Errorf("unexpected error: %v", err)
	}
}

func BenchmarkParse(b *testing.B) {
	u, err := url.Parse("http://localhost:123?name=abcd&namelist=a,b,c&namelistdash=abc&age=12")
	if err != nil {
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := testData{}
		err = queryparam.Parse(u, &data)
		if err != nil {
			b.FailNow()
		}
	}
	b.StopTimer()
	b.ReportAllocs()
}
