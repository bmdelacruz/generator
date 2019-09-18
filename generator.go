package generator

type Generator struct {
	isStarted bool
	isDone    bool

	isDoneChan          chan struct{}
	statusChan          chan *status
	retStatusChan       chan retStatus
	unhandledReturnChan chan interface{}
	unhandledErrorChan  chan error
}

type GeneratorFunc func(controller *Controller) (interface{}, error)

func New(generatorFunc GeneratorFunc) *Generator {
	generator := &Generator{
		isStarted: false,
		isDone:    false,

		isDoneChan:          make(chan struct{}),
		statusChan:          make(chan *status),
		retStatusChan:       make(chan retStatus),
		unhandledReturnChan: make(chan interface{}, 1),
		unhandledErrorChan:  make(chan error, 1),
	}

	go generator.start(generatorFunc)

	return generator
}

func (g *Generator) Next(value interface{}) (interface{}, bool, error) {
	if g.isDone {
		return nil, true, nil
	}
	g.retStatusChan <- &yieldRetStatus{value}
	g.isDoneChan <- struct{}{}
	return (<-g.statusChan).Data()
}

func (g *Generator) Return(value interface{}) (interface{}, bool, error) {
	if g.isDone {
		return nil, true, nil
	}
	g.retStatusChan <- &returnRetStatus{value}
	g.isDone = true
	g.isDoneChan <- struct{}{}
	return (<-g.statusChan).Data()
}

func (g *Generator) Error(err error) (interface{}, bool, error) {
	if g.isDone {
		return nil, true, nil
	}
	g.retStatusChan <- &errorRetStatus{err}
	g.isDoneChan <- struct{}{}
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
	<-g.isDoneChan

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
		g.isDone = true
	}

	g.statusChan <- &status{
		value: value,
		done:  true,
		err:   err,
	}
}
