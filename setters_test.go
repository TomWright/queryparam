package queryparam_test

import (
	"errors"
	"github.com/tomwright/queryparam/v4"
	"reflect"
	"testing"
)

func TestGenericSetter(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		data := struct {
			Value  string
			Target string
		}{}

		err := queryparam.GenericSetter(reflect.ValueOf(&data).Elem().FieldByName("Value"), reflect.ValueOf(&data).Elem().FieldByName("Target"))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if data.Value != data.Target {
			t.Errorf("expected %v, got %v", data.Value, data.Target)
		}
	})
	t.Run("BadValue", func(t *testing.T) {
		data := struct {
			Value  string
			Target int
		}{}

		err := queryparam.GenericSetter(reflect.ValueOf(&data).Elem().FieldByName("Value"), reflect.ValueOf(&data).Elem().FieldByName("Target"))
		expErr := errors.New("reflect.Set: value of type string is not assignable to type int")
		if err == nil || err.Error() != expErr.Error() {
			t.Errorf("expected error `%v`, got `%v`", expErr, err)
			return
		}
	})
}
