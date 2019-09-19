package generator_test

import (
	"fmt"

	"github.com/bmdelacruz/generator"
)

func ExampleGenerator_Next() {
	// Generator func that doesn't call any of the controller
	// functions and returns a nil value
	fmt.Println("Case 1:")

	g1 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return nil, nil
		},
	)
	v, r, e := g1.Next(nil) // will receive (<nil> true <nil>)
	fmt.Println("Generator#Next(nil) returns (", v, r, e, ")")
	v, r, e = g1.Next(nil) // will receive (<nil> true <nil>)
	fmt.Println("Generator#Next(nil) returns (", v, r, e, ")")

	// Generator func that doesn't call any of the controller
	// functions and returns a non-nil value
	fmt.Println("Case 2:")

	g2 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			return "yay!", nil
		},
	)
	v, r, e = g2.Next(nil) // will receive (yay! true <nil>)
	fmt.Println("Generator#Next(nil) returns (", v, r, e, ")")
	v, r, e = g2.Next(nil) // will receive (<nil> true <nil>)
	fmt.Println("Generator#Next(nil) returns (", v, r, e, ")")

	// Generator func that calls yields `1`, receives `"b"`,
	// and then returns a nil value
	fmt.Println("Case 3:")

	g3 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1) // will receive (b false <nil>)
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			return nil, nil
		},
	)
	v, r, e = g3.Next("a") // will receive (1 false <nil>)
	fmt.Println("Generator#Next(\"a\") returns (", v, r, e, ")")
	v, r, e = g3.Next("b") // will receive (<nil> true <nil>)
	fmt.Println("Generator#Next(\"b\") returns (", v, r, e, ")")
	v, r, e = g3.Next("c") // will receive (<nil> true <nil>)
	fmt.Println("Generator#Next(\"c\") returns (", v, r, e, ")")

	// Output:
	// Case 1:
	// Generator#Next(nil) returns ( <nil> true <nil> )
	// Generator#Next(nil) returns ( <nil> true <nil> )
	// Case 2:
	// Generator#Next(nil) returns ( yay! true <nil> )
	// Generator#Next(nil) returns ( <nil> true <nil> )
	// Case 3:
	// Generator#Next("a") returns ( 1 false <nil> )
	// Controller#Yield(1) returns ( b false <nil> )
	// Generator#Next("b") returns ( <nil> true <nil> )
	// Generator#Next("c") returns ( <nil> true <nil> )
}

func ExampleGenerator_Return() {
	fmt.Println("Case 1:")
	// `Generator#Return` will wait until the generator function stops
	// from executing which in turn execute all the `Println`s within
	// the generator `Func` before the `Generator#Return`'s result is
	// printed.
	//
	// Note that this does not have the same behaviour as the generator
	// in JS. In JS, calling iterator's return will immediately return
	// from the current yield statement; the statements that come after
	// that won't be executed.
	g1 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			v, r, e = gc.Yield(2)
			fmt.Println("Controller#Yield(2) returns (", v, r, e, ")")

			return nil, nil
		},
	)
	v, r, e := g1.Return("a")
	fmt.Println("Generator#Return(\"a\") returns (", v, r, e, ")")
	v, r, e = g1.Next("b")
	fmt.Println("Generator#Next(\"b\") returns (", v, r, e, ")")

	fmt.Println("Case 2:")
	g2 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			v, r, e = gc.Yield(2)
			fmt.Println("Controller#Yield(2) returns (", v, r, e, ")")

			return nil, nil
		},
	)

	v, r, e = g2.Next("a")
	fmt.Println("Generator#Next(\"a\") returns (", v, r, e, ")")
	v, r, e = g2.Return("b")
	fmt.Println("Generator#Return(\"b\") returns (", v, r, e, ")")
	v, r, e = g2.Next("c")
	fmt.Println("Generator#Next(\"c\") returns (", v, r, e, ")")

	// Output:
	// Case 1:
	// Controller#Yield(1) returns ( a true <nil> )
	// Controller#Yield(2) returns ( <nil> true <nil> )
	// Generator#Return("a") returns ( <nil> true <nil> )
	// Generator#Next("b") returns ( <nil> true <nil> )
	// Case 2:
	// Generator#Next("a") returns ( 1 false <nil> )
	// Controller#Yield(1) returns ( b true <nil> )
	// Controller#Yield(2) returns ( <nil> true <nil> )
	// Generator#Return("b") returns ( <nil> true <nil> )
	// Generator#Next("c") returns ( <nil> true <nil> )
}

func ExampleGenerator_Error() {
	fmt.Println("Case 1:")
	g1 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			v, r, e = gc.Yield(2)
			fmt.Println("Controller#Yield(2) returns (", v, r, e, ")")

			return nil, nil
		},
	)
	v, r, e := g1.Error(
		fmt.Errorf("some_error"),
	)
	fmt.Println("Generator#Error(fmt.Errorf(\"some_error\")) returns (", v, r, e, ")")
	v, r, e = g1.Next("b")
	fmt.Println("Generator#Next(\"b\") returns (", v, r, e, ")")
	v, r, e = g1.Next("c")
	fmt.Println("Generator#Next(\"c\") returns (", v, r, e, ")")

	fmt.Println("Case 2:")
	g2 := generator.New(
		func(gc *generator.Controller) (interface{}, error) {
			v, r, e := gc.Yield(1)
			fmt.Println("Controller#Yield(1) returns (", v, r, e, ")")

			v, r, e = gc.Yield(2)
			fmt.Println("Controller#Yield(2) returns (", v, r, e, ")")

			return nil, nil
		},
	)
	v, r, e = g2.Next("a")
	fmt.Println("Generator#Next(\"a\") returns (", v, r, e, ")")
	v, r, e = g2.Error(
		fmt.Errorf("some_error"),
	)
	fmt.Println("Generator#Error(fmt.Errorf(\"some_error\")) returns (", v, r, e, ")")
	v, r, e = g2.Next("b")
	fmt.Println("Generator#Next(\"b\") returns (", v, r, e, ")")

	// Output:
	// Case 1:
	// Generator#Error(fmt.Errorf("some_error")) returns ( <nil> false <nil> )
	// Controller#Yield(1) returns ( <nil> false some_error )
	// Generator#Next("b") returns ( 2 false <nil> )
	// Controller#Yield(2) returns ( c false <nil> )
	// Generator#Next("c") returns ( <nil> true <nil> )
	// Case 2:
	// Generator#Next("a") returns ( 1 false <nil> )
	// Controller#Yield(1) returns ( <nil> false some_error )
	// Generator#Error(fmt.Errorf("some_error")) returns ( 2 false <nil> )
	// Controller#Yield(2) returns ( b false <nil> )
	// Generator#Next("b") returns ( <nil> true <nil> )
}
