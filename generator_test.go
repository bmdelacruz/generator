package generator_test

import (
	"fmt"

	"github.com/bmdelacruz/generator"
)

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
