package generator_test

import (
	"fmt"
	"testing"

	"github.com/bmdelacruz/generator"
)

func TestGenerator_Next(t *testing.T) {
	t.Run("Call_with_nil_on_empty_generator", func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) interface{} {
				return nil
			},
		)
		value, isDone, err := g.Next(nil)
		if value != nil || !isDone || err != nil {
			t.Fail()
		}
		value, isDone, err = g.Next(nil)
		if value != nil || !isDone || err != nil {
			t.Fail()
		}
	})
	t.Run("Call_with_non-nil_on_empty_generator", func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) interface{} {
				return nil
			},
		)
		value, isDone, err := g.Next("a")
		if value != nil || !isDone || err != nil {
			t.Fail()
		}
		value, isDone, err = g.Next("b")
		if value != nil || !isDone || err != nil {
			t.Fail()
		}
	})
}

func TestGenerator_Return(t *testing.T) {
	t.Run("Call_with_nil_on_empty_generator", func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) interface{} {
				return nil
			},
		)
		value, isDone, err := g.Return(nil)
		if value != nil || !isDone || err != nil {
			t.Fail()
		}
		value, isDone, err = g.Return(nil)
		if value != nil || !isDone || err != nil {
			t.Fail()
		}
	})
	t.Run("Call_with_non-nil_on_empty_generator", func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) interface{} {
				return nil
			},
		)
		value, isDone, err := g.Return("a")
		if value != "a" || !isDone || err != nil {
			t.Fail()
		}
		value, isDone, err = g.Return("b")
		if value != "b" || !isDone || err != nil {
			t.Fail()
		}
	})
}

func ExampleGenerator_Return_without_value() {
	g := generator.New(
		func(gc *generator.Controller) interface{} {
			return nil
		},
	)

	v, r, e := g.Next(nil)
	fmt.Println("next", v, r, e)

	// Output:
	// next <nil> true <nil>
}

func ExampleGenerator_Return_with_value() {
	g := generator.New(
		func(gc *generator.Controller) interface{} {
			return "yay!"
		},
	)

	v, r, e := g.Next(nil)
	fmt.Println("next", v, r, e)

	// Output:
	// next yay! true <nil>
}

func ExampleGenerator_Yield_and_then_return_without_value() {
	g := generator.New(
		func(gc *generator.Controller) interface{} {
			gc.Yield(1)
			gc.Yield(2)
			return nil
		},
	)

	v, r, e := g.Next(nil)
	fmt.Println("next", v, r, e)

	v, r, e = g.Next(nil)
	fmt.Println("next", v, r, e)

	v, r, e = g.Next(nil)
	fmt.Println("next", v, r, e)

	// Output:
	// next 1 false <nil>
	// next 2 false <nil>
	// next <nil> true <nil>
}
