package queryparam_test

import (
	"errors"
	"github.com/tomwright/queryparam/v3"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestTimeValueParser(t *testing.T) {
	t.Parallel()

	type testData struct {
		CreatedAt time.Time `queryparam:"created-at"`
	}

	p := &queryparam.Parser{
		Tag: "queryparam",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf(time.Time{}): queryparam.TimeValueParser,
		},
	}

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("created-at", "2019-02-05T13:32:02Z")
		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		exp := testData{
			CreatedAt: time.Date(2019, 2, 5, 13, 32, 2, 0, time.UTC),
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("created-at", "asd")
		err := p.Parse(urlValues, &res)
		var timeErr *time.ParseError
		if err == nil || !errors.As(err, &timeErr) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}

func TestStringValueParser(t *testing.T) {
	t.Parallel()

	type testData struct {
		Name string `queryparam:"name"`
	}

	p := &queryparam.Parser{
		Tag:          "queryparam",
		DelimiterTag: "queryparamdelim",
		Delimiter:    ",",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf(""): queryparam.StringValueParser,
		},
	}

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("name", "Tom")
		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		exp := testData{
			Name: "Tom",
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})
}

func TestStringSliceValueParser(t *testing.T) {
	t.Parallel()

	type testData struct {
		Names []string `queryparam:"name"`
	}

	p := &queryparam.Parser{
		Tag:          "queryparam",
		DelimiterTag: "queryparamdelim",
		Delimiter:    ",",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf([]string{}): queryparam.StringSliceValueParser,
		},
	}

	t.Run("SingleValue", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("name", "Tom")
		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		exp := testData{
			Names: []string{"Tom"},
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})

	t.Run("MultipleValues", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("name", "Tom,Jim,Amelia")
		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		exp := testData{
			Names: []string{"Tom", "Jim", "Amelia"},
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})

	t.Run("MultipleValuesCustomDelimiter", func(t *testing.T) {
		t.Parallel()

		type testData struct {
			Names []string `queryparam:"name" queryparamdelim:"-"`
		}

		var res testData
		urlValues := url.Values{}
		urlValues.Set("name", "Tom-Jim-Amelia")
		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		exp := testData{
			Names: []string{"Tom", "Jim", "Amelia"},
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})
}

func TestIntValueParser(t *testing.T) {
	t.Parallel()

	type testData struct {
		Age int `queryparam:"age"`
	}

	p := &queryparam.Parser{
		Tag:          "queryparam",
		DelimiterTag: "queryparamdelim",
		Delimiter:    ",",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf(0): queryparam.IntValueParser,
		},
	}

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("age", "26")

		if err := p.Parse(urlValues, &res); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		exp := testData{
			Age: 26,
		}
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("age", "asd")

		err := p.Parse(urlValues, &res)
		var numErr *strconv.NumError
		if err == nil || !errors.As(err, &numErr) || !errors.Is(numErr.Err, strconv.ErrSyntax) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}

func TestBoolValueParser(t *testing.T) {
	t.Parallel()

	type testData struct {
		Accept bool `queryparam:"accept"`
	}

	p := &queryparam.Parser{
		Tag:          "queryparam",
		DelimiterTag: "queryparamdelim",
		Delimiter:    ",",
		ValueParsers: map[reflect.Type]queryparam.ValueParser{
			reflect.TypeOf(false): queryparam.BoolValueParser,
		},
	}

	t.Run("ValidTrue", func(t *testing.T) {
		t.Parallel()

		vals := []string{"true", "TRUE", "1", "yes", "YES", "y", "Y"}
		for _, v := range vals {
			t.Run(v, func(t *testing.T) {
				urlValues := url.Values{
					"accept": []string{v},
				}
				var res testData

				if err := p.Parse(urlValues, &res); err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				exp := testData{
					Accept: true,
				}
				if !reflect.DeepEqual(exp, res) {
					t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
				}
			})
		}
	})

	t.Run("ValidFalse", func(t *testing.T) {
		t.Parallel()

		vals := []string{"false", "FALSE", "0", "no", "NO", "n", "N"}
		for _, v := range vals {
			t.Run(v, func(t *testing.T) {
				urlValues := url.Values{
					"accept": []string{v},
				}
				var res testData

				if err := p.Parse(urlValues, &res); err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				exp := testData{
					Accept: false,
				}
				if !reflect.DeepEqual(exp, res) {
					t.Errorf("expected result:\n%v\ngot result:\n%v\n", exp, res)
				}
			})
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		var res testData
		urlValues := url.Values{}
		urlValues.Set("accept", "asd")

		err := p.Parse(urlValues, &res)
		if err == nil || !errors.Is(err, queryparam.ErrInvalidBoolValue) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}
