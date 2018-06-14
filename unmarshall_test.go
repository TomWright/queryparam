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

func TestUnmarshall(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := &TestRequestOne{}

	err = queryparam.Unmarshall(u, req)
	a.NoError(err)

	a.EqualValues("Tom", req.Name)
	a.EqualValues("23", req.Age)
	a.EqualValues("", req.Gender)
}

func TestUnmarshall_UnusedInvalidField(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23&gender=male")
	a.NoError(err)

	req := &TestRequestThree{}

	err = queryparam.Unmarshall(u, req)
	a.NoError(err)

	a.EqualValues("Tom", req.Name)
	a.EqualValues(0, req.Age)
	a.EqualValues("male", req.Gender)
}

func TestUnmarshall_InvalidURL(t *testing.T) {
	a := assert.New(t)

	req := &TestRequestOne{}

	err := queryparam.Unmarshall(nil, req)
	a.Equal(err, queryparam.ErrInvalidURL)
}

func TestUnmarshall_NonPointerTarget(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := TestRequestOne{}

	err = queryparam.Unmarshall(u, req)
	a.Equal(err, queryparam.ErrNonPointerTarget)
}

func TestUnmarshall_InvalidFieldType(t *testing.T) {
	a := assert.New(t)

	u, err := url.Parse("https://example.com/some/path?name=Tom&age=23")
	a.NoError(err)

	req := &TestRequestTwo{}

	err = queryparam.Unmarshall(u, req)
	a.EqualError(err, "invalid field type. `Age` must be a string")
}
