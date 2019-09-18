package generator

type Generator struct {
	isDone bool

	isDoneChan    chan struct{}
	statusChan    chan *status
	retStatusChan chan retStatus
	firstCallChan chan firstCall
}

type Func func(controller *Controller) (interface{}, error)

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

func (g *Generator) start(generatorFunc Func) {
	controller := &Controller{g: g}

	rs := <-g.retStatusChan
	v, _, e := rs.Data()
	switch rs.Type() {
	case "yield":
	case "error":
		g.firstCallChan <- &errorFirstCall{e}
	case "return":
		g.firstCallChan <- &returnFirstCall{v}
	}
	close(g.firstCallChan)
	<-g.isDoneChan

	value, err := generatorFunc(controller)

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
		g.isDone = true
	}

	g.statusChan <- &status{
		value: value,
		done:  true,
		err:   err,
	}
}
