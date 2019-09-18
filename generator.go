package generator

type Generator struct {
	isStarted  bool
	isDoneFlag bool

	doneChan            chan bool
	statusChan          chan *status
	retStatusChan       chan retStatus
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
		retStatusChan:       make(chan retStatus),
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
	g.retStatusChan <- &yieldRetStatus{value}
	g.doneChan <- false
	return (<-g.statusChan).Data()
}

func (g *Generator) Return(value interface{}) (interface{}, bool, error) {
	if g.isDoneFlag {
		return nil, true, nil
	}
	g.retStatusChan <- &returnRetStatus{value}
	g.isDoneFlag = true
	g.doneChan <- true
	return (<-g.statusChan).Data()
}

func (g *Generator) Error(err error) (interface{}, bool, error) {
	if g.isDoneFlag {
		return nil, true, nil
	}
	g.retStatusChan <- &errorRetStatus{err}
	g.doneChan <- false
	return (<-g.statusChan).Data()
}

func (g *Generator) start(generatorFunc GeneratorFunc) {
	controller := &Controller{g: g}

	rs := <-g.retStatusChan
	v, _, e := rs.Data()
	switch rs.Type() {
	case "yield":
	case "error":
		g.unhandledErrorChan <- e
	case "return":
		g.unhandledReturnChan <- v
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
