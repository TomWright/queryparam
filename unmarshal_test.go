package queryparam_test

import (
	"testing"
	"github.com/tomwright/queryparam"
	"net/url"
	"github.com/stretchr/testify/assert"
	"net/http"
	"fmt"
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

func TestUnmarshal_IntoString(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := &struct {
		Name   string `queryparam:"name"`
		Age    string `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.EqualValues("Tom", req.Name)
	a.EqualValues("23", req.Age)
	a.EqualValues("", req.Gender)
}

func TestUnmarshal_IntoSlice(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	a.NoError(err)

	req := &struct {
		Name []string `queryparam:"name"`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.Len(req.Name, 2)
	a.Equal("Tom", req.Name[0])
	a.Equal("Jim", req.Name[1])
}

func TestUnmarshal_IntoSlice_CustomDelimiter(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom-Jim")
	a.NoError(err)

	req := &struct {
		Name []string `queryparam:"name" queryparamdelim:"-"`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.Len(req.Name, 2)
	a.Equal("Tom", req.Name[0])
	a.Equal("Jim", req.Name[1])
}

func TestUnmarshal_IntoSlice_BlankDelimiterUsesDefaultDelimiter(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	a.NoError(err)

	req := &struct {
		Name []string `queryparam:"name" queryparamdelim:""`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)
	a.EqualValues([]string{"Tom", "Jim"}, req.Name)
}

func TestUnmarshal_IntoSlice_NilSlice(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	a.NoError(err)

	req := &struct {
		Name []string `queryparam:"name"`
	}{}

	a.Nil(req.Name)

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)
	a.EqualValues([]string{"Tom", "Jim"}, req.Name)
}

func TestUnmarshal_IntoSlice_NoParams(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u1, err := url.Parse("https://example.com/some/path")
	a.NoError(err)
	u2, err := url.Parse("https://example.com/some/path?name=")
	a.NoError(err)

	req := &struct {
		Name []string `queryparam:"name"`
	}{}

	a.Nil(req.Name)

	err = queryparam.Unmarshal(u1, req)
	a.NoError(err)
	a.Len(req.Name, 0)

	err = queryparam.Unmarshal(u2, req)
	a.NoError(err)
	a.Len(req.Name, 0)
}

func TestUnmarshal_UnusedInvalidField(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23&gender=male")
	a.NoError(err)

	req := &struct {
		Name   string `queryparam:"name"`
		Age    int    ``
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.EqualValues("Tom", req.Name)
	a.EqualValues(0, req.Age)
	a.EqualValues("male", req.Gender)
}

func TestUnmarshal_InvalidURL(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	req := &struct {
		Name   string `queryparam:"name"`
		Age    string `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err := queryparam.Unmarshal(nil, req)
	a.Equal(err, queryparam.ErrInvalidURL)
}

func TestUnmarshal_NonPointerTarget(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := struct {
		Name   string `queryparam:"name"`
		Age    string `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.Equal(err, queryparam.ErrNonPointerTarget)
}

func TestUnmarshal_InvalidFieldType(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := &struct {
		Name   string `queryparam:"name"`
		Age    int    `queryparam:"age"`
		Gender string `queryparam:"gender"`
	}{}

	err = queryparam.Unmarshal(u, req)
	a.EqualError(err, "invalid field type. `Age` must be `string` or `[]string`")
}
