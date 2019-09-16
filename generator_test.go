package generator_test

import (
	"fmt"
	"testing"

	"github.com/bmdelacruz/generator"
)

func TestGenerator_Call_next_without_values_on_empty_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)
	v, r, e := g.Next(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_with_values_on_empty_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)
	v, r, e := g.Next("a")
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next("b")
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_return_without_values_on_empty_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)
	v, r, e := g.Return(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Return(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_return_with_values_on_empty_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)
	v, r, e := g.Return("a")
	if v != "a" || !r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Return("b")
	if v != "b" || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_error_with_values_on_empty_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)

	e1 := fmt.Errorf("some error!")
	v, r, e := g.Error(e1)
	if v != nil || !r || e != e1 {
		t.Fatal(v, r, e)
	}

	sampleError2 := fmt.Errorf("some error 2!")
	v, r, e = g.Error(sampleError2)
	if v != nil || !r || e != sampleError2 {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_error_with_values_on_error_returning_empty_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, fmt.Errorf("some generator error")
		},
	)

	e1 := fmt.Errorf("some error!")
	v, r, e := g.Error(e1)
	if v != nil || !r || e == nil {
		t.Fatal(v, r, e)
	}
}

func ExampleGenerator_Return_without_value() {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)

	v, r, e := g.Next(nil)
	fmt.Println("next", v, r, e)

	// Output:
	// next <nil> true <nil>
}

func ExampleGenerator_Return_with_value() {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return "yay!", nil
		},
	)

	v, r, e := g.Next(nil)
	fmt.Println("next", v, r, e)

	// Output:
	// next yay! true <nil>
}

func ExampleGenerator_Yield_and_then_return_without_value() {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			gc.Yield(1)
			gc.Yield(2)
			return nil, nil
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
