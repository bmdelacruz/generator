package generator

import "sync"

type GeneratorFunc func(controller *Controller)

type Generator struct {
	generatorFunc GeneratorFunc
	controller    *Controller

	startOnce     *sync.Once
	stopWaitGroup *sync.WaitGroup

	Iterator *Iterator
}

func New(generatorFunc GeneratorFunc) *Generator {
	iterator := &Iterator{
		isStarted:  false,
		isDoneFlag: false,

		continueChan: make(chan struct{}),
		doneChan:     make(chan struct{}),
		statusChan:   make(chan *status),
		nextChan:     make(chan interface{}),
		returnChan:   make(chan interface{}),
		errorChan:    make(chan error),
	}
	return &Generator{
		generatorFunc: generatorFunc,
		controller:    &Controller{iterator},

		startOnce:     &sync.Once{},
		stopWaitGroup: &sync.WaitGroup{},

		Iterator: iterator,
	}
}

func (g *Generator) Start() {
	g.startOnce.Do(func() {
		go g.startGenerator()
	})
}

func (g *Generator) WaitToFinish() {
	g.stopWaitGroup.Wait()
}

func (g *Generator) startGenerator() {
	g.stopWaitGroup.Add(1)
	defer g.stopWaitGroup.Done()

	shouldContinue := g.controller.interceptFirstIteration()
	if shouldContinue {
		g.generatorFunc(g.controller)
		if !g.Iterator.isDoneFlag {
			g.controller.iterator.flush()
		}
	}
}
