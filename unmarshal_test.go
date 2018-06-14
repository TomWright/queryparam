package queryparam_test

import (
	"testing"
	"github.com/tomwright/queryparam"
	"net/url"
	"github.com/stretchr/testify/assert"
)

type TestRequestOne struct {
	Name   string `queryparam:"name"`
	Age    string `queryparam:"age"`
	Gender string `queryparam:"gender"`
}

type TestRequestTwo struct {
	Name   string `queryparam:"name"`
	Age    int    `queryparam:"age"`
	Gender string `queryparam:"gender"`
}

type TestRequestThree struct {
	Name   string `queryparam:"name"`
	Age    int    ``
	Gender string `queryparam:"gender"`
}

type TestRequestFour struct {
	Name []string `queryparam:"name"`
}

func TestUnmarshal_IntoString(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := &TestRequestOne{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.EqualValues("Tom", req.Name)
	a.EqualValues("23", req.Age)
	a.EqualValues("", req.Gender)
}

func TestUnmarshal_IntoSlice(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	a.NoError(err)

	req := &TestRequestFour{}
	req.Name = make([]string, 0)

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.Len(req.Name, 2)
	a.Equal("Tom", req.Name[0])
	a.Equal("Jim", req.Name[1])
}

func TestUnmarshal_IntoSlice_CustomDelimiter(t *testing.T) {
	defer func() {
		queryparam.Delimiter = ","
	}()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom-Jim")
	a.NoError(err)

	req := &TestRequestFour{}
	req.Name = make([]string, 0)

	queryparam.Delimiter = "-"

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.Len(req.Name, 2)
	a.Equal("Tom", req.Name[0])
	a.Equal("Jim", req.Name[1])
}

func TestUnmarshal_IntoSlice_InvalidDelimiter(t *testing.T) {
	defer func() {
		queryparam.Delimiter = ","
	}()

	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom-Jim")
	a.NoError(err)

	req := &TestRequestFour{}
	req.Name = make([]string, 0)

	queryparam.Delimiter = ""

	err = queryparam.Unmarshal(u, req)
	a.Equal(queryparam.ErrInvalidDelimiter, err)
}

func TestUnmarshal_IntoSlice_NilSlice(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom,Jim")
	a.NoError(err)

	req := &TestRequestFour{}

	err = queryparam.Unmarshal(u, req)
	a.Equal(queryparam.ErrNilSliceField, err)
}

func TestUnmarshal_UnusedInvalidField(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23&gender=male")
	a.NoError(err)

	req := &TestRequestThree{}

	err = queryparam.Unmarshal(u, req)
	a.NoError(err)

	a.EqualValues("Tom", req.Name)
	a.EqualValues(0, req.Age)
	a.EqualValues("male", req.Gender)
}

func TestUnmarshal_InvalidURL(t *testing.T) {
	a := assert.New(t)

	req := &TestRequestOne{}

	err := queryparam.Unmarshal(nil, req)
	a.Equal(err, queryparam.ErrInvalidURL)
}

func TestUnmarshal_NonPointerTarget(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := TestRequestOne{}

	err = queryparam.Unmarshal(u, req)
	a.Equal(err, queryparam.ErrNonPointerTarget)
}

func TestUnmarshal_InvalidFieldType(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := &TestRequestTwo{}

	err = queryparam.Unmarshal(u, req)
	a.EqualError(err, "invalid field type. `Age` must be `string` or `[]string`")
}
