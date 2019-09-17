package generator

type status struct {
	value interface{}
	done  bool
	err   error
}

type Generator struct {
	isStarted  bool
	isDoneFlag bool

	doneChan            chan bool
	statusChan          chan *status
	yieldChan           chan interface{}
	returnChan          chan interface{}
	errorChan           chan error
	unhandledReturnChan chan interface{}
	unhandledErrorChan  chan error
}

type GeneratorFunc func(controller *Controller) (interface{}, error)

func New(generatorFunc GeneratorFunc) *Generator {
	generator := &Generator{
		isStarted:  false,
		isDoneFlag: false,

		doneChan:            make(chan bool),
		statusChan:          make(chan *status),
		yieldChan:           make(chan interface{}),
		returnChan:          make(chan interface{}),
		errorChan:           make(chan error),
		unhandledReturnChan: make(chan interface{}, 1),
		unhandledErrorChan:  make(chan error, 1),
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
	return status.value, status.done, status.err
}

func (g *Generator) Return(value interface{}) (interface{}, bool, error) {
	if g.isDoneFlag {
		return nil, true, nil
	}

	g.returnChan <- value
	g.isDoneFlag = true
	g.doneChan <- true

	status := <-g.statusChan
	return status.value, status.done, status.err
}

func (g *Generator) Error(err error) (interface{}, bool, error) {
	if g.isDoneFlag {
		return nil, true, nil
	}

	g.errorChan <- err
	g.doneChan <- false

	status := <-g.statusChan
	return status.value, status.done, status.err
}

func (g *Generator) start(generatorFunc GeneratorFunc) {
	controller := &Controller{g: g}

	select {
	case <-g.yieldChan:
	case err := <-g.errorChan:
		g.unhandledErrorChan <- err
		close(g.unhandledErrorChan)
	case value := <-g.returnChan:
		g.unhandledReturnChan <- value
		close(g.unhandledReturnChan)
	}

	<-g.doneChan

	value, err := generatorFunc(controller)

	if !controller.wasUsed {
		select {
		case unhandledReturn := <-g.unhandledReturnChan:
			value = unhandledReturn
		default:
		}
		select {
		case unhandledErr := <-g.unhandledErrorChan:
			err = unhandledErr
		default:
		}

		g.isDoneFlag = true
	}

	g.statusChan <- &status{
		value: value,
		done:  true,
		err:   err,
	}
}
