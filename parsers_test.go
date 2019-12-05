package queryparam_test

import (
	"errors"
	"github.com/tomwright/queryparam/v4"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestStringValueParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res, err := queryparam.StringValueParser("", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := "", res.String(); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		res, err := queryparam.StringValueParser("hello", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := "hello", res.String(); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
}

func TestStringSliceValueParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res, err := queryparam.StringSliceValueParser("", ",")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		exp := make([]string, 0)
		got, ok := res.Interface().([]string)
		if !ok || !reflect.DeepEqual(exp, got) {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("SingleValue", func(t *testing.T) {
		res, err := queryparam.StringSliceValueParser("hello", ",")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		exp := []string{"hello"}
		got, ok := res.Interface().([]string)
		if !ok || !reflect.DeepEqual(exp, got) {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("MultipleValueComma", func(t *testing.T) {
		res, err := queryparam.StringSliceValueParser("hello,hi,hey", ",")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		exp := []string{"hello", "hi", "hey"}
		got, ok := res.Interface().([]string)
		if !ok || !reflect.DeepEqual(exp, got) {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("MultipleValueDash", func(t *testing.T) {
		res, err := queryparam.StringSliceValueParser("hello-hi-hey", "-")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		exp := []string{"hello", "hi", "hey"}
		got, ok := res.Interface().([]string)
		if !ok || !reflect.DeepEqual(exp, got) {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
}

func TestTimeValueParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res, err := queryparam.TimeValueParser("", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		got, ok := res.Interface().(time.Time)
		if !ok || !got.IsZero() {
			t.Errorf("expected res to be zero, got `%v`", got)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		res, err := queryparam.TimeValueParser("2019-02-05T13:32:02Z", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		exp := time.Date(2019, 2, 5, 13, 32, 2, 0, time.UTC)
		got, ok := res.Interface().(time.Time)
		if !ok || !exp.Equal(got) {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		res, err := queryparam.TimeValueParser("not-a-date", "")
		var timeErr *time.ParseError
		if err == nil || !errors.As(err, &timeErr) {
			t.Errorf("unexpected error: %v", err)
		}
		got, ok := res.Interface().(time.Time)
		if !ok || !got.IsZero() {
			t.Errorf("expected res to be zero, got `%v`", got)
		}
	})
}

func TestIntValueParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res, err := queryparam.IntValueParser("", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := 0, int(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		res, err := queryparam.IntValueParser("21323", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := 21323, int(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		res, err := queryparam.IntValueParser("asd", "")
		var numErr *strconv.NumError
		if err == nil || !errors.As(err, &numErr) || !errors.Is(numErr.Err, strconv.ErrSyntax) {
			t.Errorf("unexpected error: %v", err)
		}
		if exp, got := int(0), int(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
}

func TestInt32ValueParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res, err := queryparam.Int32ValueParser("", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := int32(0), int32(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		res, err := queryparam.Int32ValueParser("21323", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := int32(21323), int32(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
		t.Run("Invalid", func(t *testing.T) {
			res, err := queryparam.Int32ValueParser("asd", "")
			var numErr *strconv.NumError
			if err == nil || !errors.As(err, &numErr) || !errors.Is(numErr.Err, strconv.ErrSyntax) {
				t.Errorf("unexpected error: %v", err)
			}
			if exp, got := int32(0), int32(res.Int()); exp != got {
				t.Errorf("expected res `%v`, got `%v`", exp, got)
			}
		})
	})
}

func TestInt64ValueParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res, err := queryparam.Int64ValueParser("", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := int64(0), int64(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Valid", func(t *testing.T) {
		res, err := queryparam.Int64ValueParser("21643", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if exp, got := int64(21643), int64(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		res, err := queryparam.Int64ValueParser("asd", "")
		var numErr *strconv.NumError
		if err == nil || !errors.As(err, &numErr) || !errors.Is(numErr.Err, strconv.ErrSyntax) {
			t.Errorf("unexpected error: %v", err)
		}
		if exp, got := int64(0), int64(res.Int()); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})
}

func TestBoolValueParser(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		res, err := queryparam.BoolValueParser("non-bool-string", "")
		expErr := queryparam.ErrInvalidBoolValue
		if err == nil || !errors.Is(err, expErr) {
			t.Errorf("expected error `%v` but got `%v`", expErr, err)
			return
		}
		if exp, got := false, res.Bool(); exp != got {
			t.Errorf("expected res `%v`, got `%v`", exp, got)
		}
	})

	checkBoolFn := func(input []string, exp bool) func(*testing.T) {
		return func(t *testing.T) {
			for _, in := range input {
				i := in
				t.Run(i, func(t *testing.T) {
					res, err := queryparam.BoolValueParser(i, "")
					if err != nil {
						t.Errorf("unexpected error: %s", err.Error())
						return
					}
					if exp, got := exp, res.Bool(); exp != got {
						t.Errorf("expected res `%v`, got `%v`", exp, got)
					}
				})
			}
		}
	}
	t.Run("True", checkBoolFn([]string{"true", "TRUE", "1", "yes", "YES", "y", "Y"}, true))
	t.Run("False", checkBoolFn([]string{"", "false", "FALSE", "0", "no", "NO", "n", "N"}, false))
}
