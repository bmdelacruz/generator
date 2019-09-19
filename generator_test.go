package generator_test

import (
	"fmt"
	"testing"

	"github.com/bmdelacruz/generator"
)

func TestGenerator_Next(t *testing.T) {
	t.Run(`Next("a"),Next("b"),Next("c")|..`, func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				return 0, nil
			},
		)
		testWith(t).expect(g.Next("a")).toReturn(0, true, nil)
		testWith(t).expect(g.Next("b")).toReturn(nil, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
	t.Run(`Next("a"),Next("b"),Next("c")|Yield(1)`, func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Yield(1)).toReturn("b", false, nil)
				return 0, nil
			},
		)
		testWith(t).expect(g.Next("a")).toReturn(1, false, nil)
		testWith(t).expect(g.Next("b")).toReturn(0, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
	t.Run(`Next("a"),Next("b"),Next("c")|Error(<e1>)`, func(t *testing.T) {
		e1 := fmt.Errorf("e1")
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Error(e1)).toReturn("b", false, nil)
				return 0, nil
			},
		)
		testWith(t).expect(g.Next("a")).toReturn(nil, false, e1)
		testWith(t).expect(g.Next("b")).toReturn(0, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
}

func TestGenerator_Return(t *testing.T) {
	t.Run(`Return("a"),Return("b"),Next("c")|..`, func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				return 0, nil
			},
		)
		testWith(t).expect(g.Return("a")).toReturn("a", true, nil)
		testWith(t).expect(g.Return("b")).toReturn(nil, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
	t.Run(`Return("a"),Return("b"),Next("c")|Yield(1)`, func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Yield(1)).toReturn("a", true, nil)
				return 0, nil
			},
		)
		testWith(t).expect(g.Return("a")).toReturn(0, true, nil)
		testWith(t).expect(g.Return("b")).toReturn(nil, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
	t.Run(`Return("a"),Return("b"),Next("c")|Error(<e1>)`, func(t *testing.T) {
		e1 := fmt.Errorf("e1")
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Error(e1)).toReturn("a", true, nil)
				return 0, nil
			},
		)
		testWith(t).expect(g.Return("a")).toReturn(0, true, nil)
		testWith(t).expect(g.Return("b")).toReturn(nil, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
}

func TestGenerator_Error(t *testing.T) {
	t.Run(`Error("a"),Error("b"),Next("c")|..`, func(t *testing.T) {
		e1 := fmt.Errorf("e1")
		e2 := fmt.Errorf("e2")
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				return 0, nil
			},
		)
		testWith(t).expect(g.Error(e1)).toReturn(0, true, e1)
		testWith(t).expect(g.Error(e2)).toReturn(nil, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
	t.Run(`Error("a"),Error("b"),Next("c")|Yield(1)`, func(t *testing.T) {
		e1 := fmt.Errorf("e1")
		e2 := fmt.Errorf("e2")
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Yield(1)).toReturn(nil, false, e1)
				return 0, nil
			},
		)
		testWith(t).expect(g.Error(e1)).toReturn(nil, false, nil)
		testWith(t).expect(g.Error(e2)).toReturn(0, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
	t.Run(`Error("a"),Error("b"),Next("c")|Error(<e1>)`, func(t *testing.T) {
		e1 := fmt.Errorf("e1")
		e2 := fmt.Errorf("e2")
		e3 := fmt.Errorf("e3")
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Error(e3)).toReturn(nil, false, e1)
				return 0, nil
			},
		)
		testWith(t).expect(g.Error(e1)).toReturn(nil, false, nil)
		testWith(t).expect(g.Error(e2)).toReturn(0, true, nil)
		testWith(t).expect(g.Next("c")).toReturn(nil, true, nil)
	})
}

func TestController_Yield(t *testing.T) {
	t.Run(`Next("a"),Next("b"),Return("c")|Yield(1),Yield(2),Yield(3)`, func(t *testing.T) {
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Yield(1)).toReturn("b", false, nil)
				testWith(t).pexpect(gc.Yield(2)).toReturn("c", true, nil)
				testWith(t).pexpect(gc.Yield(3)).toReturn(nil, true, nil)
				return 0, nil
			},
		)
		testWith(t).expect(g.Next("a")).toReturn(1, false, nil)
		testWith(t).expect(g.Next("b")).toReturn(2, false, nil)
		testWith(t).expect(g.Return("c")).toReturn(0, true, nil)
		testWith(t).expect(g.Next("d")).toReturn(nil, true, nil)
	})
}

func TestController_Error(t *testing.T) {
	t.Run(`Next("a"),Next("b"),Return("c")|Yield(1),Yield(2),Yield(3)`, func(t *testing.T) {
		e1 := fmt.Errorf("e1")
		g := generator.New(
			func(gc *generator.Controller) (interface{}, error) {
				testWith(t).pexpect(gc.Yield(1)).toReturn("b", false, nil)
				testWith(t).pexpect(gc.Error(e1)).toReturn("c", true, nil)
				testWith(t).pexpect(gc.Yield(3)).toReturn(nil, true, nil)
				return 0, nil
			},
		)
		testWith(t).expect(g.Next("a")).toReturn(1, false, nil)
		testWith(t).expect(g.Next("b")).toReturn(nil, false, e1)
		testWith(t).expect(g.Return("c")).toReturn(0, true, nil)
		testWith(t).expect(g.Next("d")).toReturn(nil, true, nil)
	})
}

// utility stuff =====================================================

type tw struct {
	t *testing.T
}

type twe struct {
	tw *tw
	v  interface{}
	r  bool
	e  error
	p  bool
}

func testWith(t *testing.T) *tw {
	return &tw{t}
}

func (tw *tw) expect(v interface{}, r bool, e error) *twe {
	p := false
	return &twe{tw, v, r, e, p}
}

func (tw *tw) pexpect(v interface{}, r bool, e error) *twe {
	p := true
	return &twe{tw, v, r, e, p}
}

func (twe *twe) toReturn(v interface{}, r bool, e error) {
	if v != twe.v || r != twe.r || e != twe.e {
		if twe.p {
			panic(
				fmt.Errorf(
					"got: %v, %v, %v. wanted: %v, %v, %v.",
					twe.v, twe.r, twe.e, v, r, e,
				),
			)
		} else {
			twe.tw.t.Helper()
			twe.tw.t.Fatalf(
				"got: %v, %v, %v. wanted: %v, %v, %v.",
				twe.v, twe.r, twe.e, v, r, e,
			)
		}
	}
}
