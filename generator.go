package generator

// Generator provides functions that will send and receive data to and
// from the `Func` associated with it.
type Generator struct {
	isDone bool

	// isDoneChan is for preventing data race conditions. it is safe to
	// leave this with empty struct type.
	isDoneChan    chan struct{}
	statusChan    chan *status
	retStatusChan chan retStatus
	firstCallChan chan firstCall
}

// Func is the signature of the generator function
type Func func(controller *Controller) (interface{}, error)

// New creates an instance of a generator and spawns a goroutine where
// the generator function will run.
func New(generatorFunc Func) *Generator {
	generator := &Generator{
		isDone: false,

		isDoneChan:    make(chan struct{}),
		statusChan:    make(chan *status),
		retStatusChan: make(chan retStatus),
		firstCallChan: make(chan firstCall, 1),
	}

	go generator.start(generatorFunc)

	return generator
}

// Next provides the value that should be returned by the current
// generator controller function and retrieves the next yielded
// value from the `Func`. The argument is ignored when `Next` is
// the first generator function that was invoked.
//
// Returns ([value], [isDone], [error])
func (g *Generator) Next(value interface{}) (interface{}, bool, error) {
	if g.isDone {
		return nil, true, nil
	}
	g.retStatusChan <- &yieldRetStatus{value}
	g.isDoneChan <- struct{}{}
	return (<-g.statusChan).Data()
}

// Return provides the value the `Func` should return and tells the
// generator to stop the execution of the `Func`. The currently yielding
// generator controller function will return (<nil>, true, <nil>). The
// succeeding generator controller function invocations will return
// (<nil>, true, <nil>).
//
// If the `Func` doesn't call any of the generator controller functions
// and `Return` is called, the value returned by the `Func` will be
// replaced by the value passed as an argument to `Return`.
//
// Returns ([value], [isDone], [error])
func (g *Generator) Return(value interface{}) (interface{}, bool, error) {
	if g.isDone {
		return nil, true, nil
	}
	g.retStatusChan <- &returnRetStatus{value}
	g.isDone = true
	g.isDoneChan <- struct{}{}
	return (<-g.statusChan).Data()
}

// Error provides the error the currently yielding generator controller
// function should receive. Note that it will not stop the `Func`. The
// error should be handled from the `Func`.
//
// If the `Func` doesn't call any of the generator controller functions
// and `Error` is called, the error returned by the `Func` will be
// replaced by the error passed as an argument to `Error`.
//
// Returns ([value], [isDone], [error])
func (g *Generator) Error(err error) (interface{}, bool, error) {
	if g.isDone {
		return nil, true, nil
	}
	g.retStatusChan <- &errorRetStatus{err}
	g.isDoneChan <- struct{}{}
	return (<-g.statusChan).Data()
}

func (g *Generator) start(generatorFunc Func) {
	controller := &Controller{g: g}

	// receive the initial data sent from any of the generator functions
	rs := <-g.retStatusChan

	v, _, e := rs.Data()
	switch rs.Type() {
	case "yield":
		// ignore value from `Next`
	case "error":
		// save the error value from `Error` for later
		g.firstCallChan <- &errorFirstCall{e}
	case "return":
		// save the return value from `Return` for later
		g.firstCallChan <- &returnFirstCall{v}
	}

	// immediately close the first call channel because it's only for
	// first generator function calls
	close(g.firstCallChan)

	// receives like this from this channel would mean that the
	// `isDone` may have already been updated so it's safe to access
	// (to prevent data race)
	<-g.isDoneChan

	value, err := generatorFunc(controller)

	// this condition will be equal to true when any of the generator
	// controller functions has not been called
	if !controller.wasUsed {
		select {
		case fc, ok := <-g.firstCallChan:
			if ok {
				switch fc.Type() {
				case "error":
					_, err = fc.Values()
				case "return":
					value, _ = fc.Values()
				}
			}
		default:
		}
	}

	// don't forget to mark the generator as done. return may not have
	// been called.
	//
	// NOTE:
	// to future bryan, don't forget that `isDone` won't be accessed from
	// any generator functions until after sending to the status chan so
	// it's safe to modify here
	g.isDone = true

	// send the last status to the last proper call to any of the generator
	// functions
	g.statusChan <- &status{
		value: value,
		done:  true,
		err:   err,
	}
}
