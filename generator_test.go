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

func TestGenerator_Call_next_without_values_on_yielding_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != nil || r || e != nil {
				t.Fatal(v, r, e)
			}
			v, r, e = gc.Yield(2)
			if v != nil || r || e != nil {
				t.Fatal(v, r, e)
			}
			return nil, nil
		},
	)
	v, r, e := g.Next(nil)
	if v != 1 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next(nil)
	if v != 2 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_with_values_on_yielding_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != "a" || r || e != nil {
				t.Fatal(v, r, e)
			}
			v, r, e = gc.Yield(2)
			if v != "b" || r || e != nil {
				t.Fatal(v, r, e)
			}
			return nil, nil
		},
	)
	v, r, e := g.Next(nil)
	if v != 1 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next("a")
	if v != 2 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next("b")
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_return_without_values_on_yielding_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != nil || !r || e != nil {
				t.Fatal(v, r, e)
			}
			return v, nil
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

func TestGenerator_Call_return_with_values_on_yielding_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != "a" || !r || e != nil {
				t.Fatal(v, r, e)
			}
			return v, nil
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

func TestGenerator_Call_error_with_values_on_yielding_generator_function(t *testing.T) {
	e1 := fmt.Errorf("e1")
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != nil || r || e != e1 {
				t.Fatal(v, r, e)
			}
			return nil, e
		},
	)
	v, r, e := g.Error(e1)
	if v != nil || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next(e1)
	if v != nil || !r || e != e1 {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_next_and_return_with_values_on_yielding_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != "a" || r || e != nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			v, r, e = gc.Yield(2)
			if v != "b" || !r || e != nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return v, nil
		},
	)
	v, r, e := g.Next(nil)
	if v != 1 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next("a")
	if v != 2 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Return("b")
	if v != "b" || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_next_and_error_with_values_on_yielding_generator_function(t *testing.T) {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			if v != "a" || r || e != nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			v, r, e = gc.Yield(2)
			if v != nil || r || e == nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return nil, e
		},
	)
	v, r, e := g.Next(nil)
	if v != 1 || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next("a")
	if v != 2 || r || e != nil {
		t.Fatal(v, r, e)
	}
	e1 := fmt.Errorf("e1")
	v, r, e = g.Error(e1)
	if v != nil || !r || e != e1 {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_and_next_with_values_on_erroring_generator_function(t *testing.T) {
	e1 := fmt.Errorf("e1")
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Error(e1)
			if v != nil || r || e != nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return nil, nil
		},
	)
	v, r, e := g.Next(nil)
	if v != nil || r || e != e1 {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_return_without_values_on_erroring_generator_function(t *testing.T) {
	e1 := fmt.Errorf("e1")
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Error(e1)
			if v != nil || !r || e != nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return nil, nil
		},
	)
	v, r, e := g.Return(nil)
	if v != nil || !r || e != nil {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_error_and_next_without_values_on_erroring_generator_function(t *testing.T) {
	e1 := fmt.Errorf("e1")
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Error(e1)
			if v != nil || r || e != e1 {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return nil, e
		},
	)
	v, r, e := g.Error(e1)
	if v != nil || r || e != nil {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Next(nil)
	if v != nil || !r || e != e1 {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_and_error_with_values_on_erroring_generator_function(t *testing.T) {
	e1 := fmt.Errorf("e1")
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Error(e1)
			if v != nil || r || e != e1 {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return nil, e
		},
	)
	v, r, e := g.Next(nil)
	if v != nil || r || e != e1 {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Error(e1)
	if v != nil || !r || e != e1 {
		t.Fatal(v, r, e)
	}
}

func TestGenerator_Call_next_and_return_with_values_on_erroring_generator_function(t *testing.T) {
	e1 := fmt.Errorf("e1")
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Error(e1)
			if v != "a" || !r || e != nil {
				panic(fmt.Errorf("%v %v %v", v, r, e))
			}
			return v, nil
		},
	)
	v, r, e := g.Next(nil)
	if v != nil || r || e != e1 {
		t.Fatal(v, r, e)
	}
	v, r, e = g.Return("a")
	if v != "a" || !r || e != nil {
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
