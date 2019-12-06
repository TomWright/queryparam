package queryparam

import (
	"errors"
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	t.Run("NoPanic", func(t *testing.T) {
		got := func() (err error) {
			defer recoverPanic(&err)()
			return nil
		}()

		if got != nil {
			t.Errorf("unexpected error")
		}
	})
	t.Run("PanicString", func(t *testing.T) {
		got := func() (err error) {
			defer recoverPanic(&err)()
			panic("something bad")
			return
		}()

		var exp error = errors.New("something bad")
		if got == nil || exp.Error() != got.Error() {
			t.Errorf("expected `%v`, got `%v`", exp, got)
		}
	})
	t.Run("PanicError", func(t *testing.T) {
		exp := errors.New("something bad")
		got := func() (err error) {
			defer recoverPanic(&err)()
			panic(exp)
			return
		}()

		if got == nil || !errors.Is(got, exp) {
			t.Errorf("expected `%v`, got `%v`", exp, got)
		}
	})
}
