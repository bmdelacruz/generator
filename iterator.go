package generator

type Iterator struct {
	isStarted  bool
	isDoneFlag bool

	continueChan chan struct{}
	doneChan     chan struct{}
	statusChan   chan *status
	nextChan     chan interface{}
	returnChan   chan interface{}
	errorChan    chan error
}

func (i *Iterator) Next(value interface{}) (interface{}, bool, error) {
	if i.isDoneFlag {
		return nil, true, nil
	}

	i.nextChan <- value
	i.continueChan <- struct{}{}

	status := <-i.statusChan
	return status.value, status.done, status.err
}

func (i *Iterator) Return(value interface{}) (interface{}, bool, error) {
	if i.isDoneFlag {
		return value, true, nil
	}

	i.returnChan <- value
	i.doneChan <- struct{}{}

	return value, true, nil
}

func (i *Iterator) Error(err error) (interface{}, bool, error) {
	if i.isDoneFlag {
		return nil, true, err
	}

	i.errorChan <- err
	i.continueChan <- struct{}{}

	status := <-i.statusChan
	return status.value, status.done, status.err
}

func (i *Iterator) flush() {
	i.statusChan <- &status{
		value: nil,
		done:  true,
		err:   nil,
	}
}

func (i *Iterator) isDone() bool {
	if i.isDoneFlag {
		return true
	}
	select {
	case <-i.continueChan:
		return false
	case <-i.doneChan:
		i.isDoneFlag = true
		return true
	}
}
