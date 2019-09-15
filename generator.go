package generator

type status struct {
	value interface{}
	done  bool
	err   error
}

type Generator struct {
	isStarted  bool
	isDoneFlag bool

	doneChan   chan bool
	statusChan chan *status
	yieldChan  chan interface{}
	returnChan chan interface{}
	errorChan  chan error
}

func New(generatorFunc func(controller *Controller) interface{}) *Generator {
	generator := &Generator{
		isStarted:  false,
		isDoneFlag: false,

		doneChan:   make(chan bool),
		statusChan: make(chan *status),
		yieldChan:  make(chan interface{}),
		returnChan: make(chan interface{}),
		errorChan:  make(chan error),
	}

	go generator.start(generatorFunc)

	return generator
}

func (g *Generator) Next(value interface{}) (interface{}, bool, error) {
	if g.isDoneFlag {
		return nil, true, nil
	}

	g.yieldChan <- value
	g.doneChan <- false

	status := <-g.statusChan
	if status.done {
		g.isDoneFlag = true
	}
	return status.value, status.done, status.err
}

func (g *Generator) Return(value interface{}) (interface{}, bool, error) {
	if g.isDoneFlag {
		return value, true, nil
	}

	g.returnChan <- value
	g.doneChan <- true

	return value, true, nil
}

func (g *Generator) Error(err error) (interface{}, bool, error) {
	if g.isDoneFlag {
		return nil, true, err
	}

	g.errorChan <- err
	g.doneChan <- false

	status := <-g.statusChan
	if status.done {
		g.isDoneFlag = true
	}
	return status.value, status.done, status.err
}

func (g *Generator) start(
	generatorFunc func(controller *Controller) interface{},
) {
	controller := &Controller{g}

	select {
	case <-g.yieldChan:
	case <-g.returnChan:
	case <-g.errorChan:
	}

	if !g.isDone() {
		retVal := generatorFunc(controller)
		g.statusChan <- &status{
			value: retVal,
			done:  true,
			err:   nil,
		}
	}
}

func (g *Generator) isDone() bool {
	if g.isDoneFlag {
		return true
	}

	isDone := <-g.doneChan
	if isDone {
		g.isDoneFlag = true
	}

	return isDone
}
