package generator_test

import (
	"fmt"

	"github.com/bmdelacruz/generator"
)

func ExampleController_Yield() {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			return nil, nil
		},
	)
	v, r, e := g.Next("a")
	fmt.Println("Generator#Next(\"a\") returns (", v, r, e, ")")
	v, r, e = g.Next("b")
	fmt.Println("Generator#Next(\"b\") returns (", v, r, e, ")")

	// Output:
	// Generator#Next("a") returns ( 1 false <nil> )
	// Controller#Yield(1) returns ( b false <nil> )
	// Generator#Next("b") returns ( <nil> true <nil> )
}

func ExampleController_Error() {
	g := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Error(fmt.Errorf("some_error"))
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			return nil, nil
		},
	)
	v, r, e := g.Next("a")
	fmt.Println("Generator#Next(\"a\") returns (", v, r, e, ")")
	v, r, e = g.Next("b")
	fmt.Println("Generator#Next(\"b\") returns (", v, r, e, ")")

	// Output:
	// Generator#Next("a") returns ( <nil> false some_error )
	// Controller#Yield(1) returns ( b false <nil> )
	// Generator#Next("b") returns ( <nil> true <nil> )
}
