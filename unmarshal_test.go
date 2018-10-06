package queryparam_test

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/tomwright/queryparam"
)

// ExampleUnmarshal creates a dummy http request and unmarshals the data into a struct
func ExampleUnmarshal() {
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

	err = queryparam.Unmarshal(request.URL, &requestData)
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
	Age          string   `queryparam:"age"`
}

func TestUnmarshal(t *testing.T) {
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
				Age:  "23",
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
				Age:          "23",
			},
		},
		{
			TestDesc: "CommaAndDashDelimitedStrings",
			URL:      "https://example.com/some/path?names=x,y-z&age=2-3",
			Data:     testData{},
			OutputData: testData{
				NameList:     []string{"x", "y-z"},
				NameListDash: []string{"x,y", "z"},
				Age:          "2-3",
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.TestDesc, func(st *testing.T) {
			st.Parallel()

			u, err := url.Parse(tc.URL)
			if err != nil {
				st.Errorf("unable to parse url: %s", err)
			}

			err = queryparam.Unmarshal(u, &tc.Data)
			if err != nil {
				st.Errorf("unable to unmarshal url: %s", err)
			}

			if !reflect.DeepEqual(tc.Data, tc.OutputData) {
				st.Errorf("expected `%v`, got `%v`", tc.OutputData, tc.Data)
			}
		})
	}
}

func TestUnmarshal_BlankDelimiterUsesDefaultDelimiter(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name []string `queryparam:"name" queryparamdelim:""`
	}{}

	err = queryparam.Unmarshal(u, req)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if exp, got := []string{"Tom", "Jim"}, req.Name; !reflect.DeepEqual(exp, got) {
		t.Errorf("unexpected result. expected `%v`, got `%v`", exp, got)
	}
}

func TestUnmarshal_UnusedInvalidField(t *testing.T) {
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

	err = queryparam.Unmarshal(u, req)
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

func TestUnmarshal_InvalidURL(t *testing.T) {
	t.Parallel()

	req := &struct {
		Name   string `queryparam:"name"`
		Age    string `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err := queryparam.Unmarshal(nil, req)
	if exp, got := queryparam.ErrInvalidURL, err; exp != got {
		t.Errorf("unexpected error. expected `%v`, got `%v`", exp, got)
	}
}

func TestUnmarshal_NonPointerTarget(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := struct {
		Name   string `queryparam:"name"`
		Age    string `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Unmarshal(u, req)
	if exp, got := queryparam.ErrNonPointerTarget, err; exp != got {
		t.Errorf("unexpected error. expected `%v`, got `%v`", exp, got)
	}
}

func TestUnmarshal_InvalidFieldType(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	req := &struct {
		Name   string `queryparam:"name"`
		Age    int    `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Unmarshal(u, req)
	if exp, got := "invalid field type. `Age` must be `string` or `[]string`", err.Error(); exp != got {
		t.Errorf("unexpected error string. expected `%s`, got `%s`", exp, got)
	}
}

func BenchmarkUnmarshal(b *testing.B) {

	u, err := url.Parse("http://localhost:123?name=abcd&namelist=a,b,c&namelistdash=abc&age=12")
	if err != nil {
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := testData{}
		err = queryparam.Unmarshal(u, &data)
		if err != nil {
			b.FailNow()
		}
	}
	b.StopTimer()
	b.ReportAllocs()
}
