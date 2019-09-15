package generator

type Controller struct {
	iterator *Iterator
}

func (c *Controller) Yield(value interface{}) (interface{}, bool, error) {
	c.iterator.statusChan <- &status{
		value: value,
		done:  false,
		err:   nil,
	}

	select {
	case value := <-c.iterator.nextChan:
		if c.iterator.isDone() {
			return nil, true, nil
		}
		return value, false, nil
	case err := <-c.iterator.errorChan:
		c.iterator.isDone() // don't care
		return nil, false, err
	case <-c.iterator.returnChan:
		c.iterator.isDone() // don't care
		return nil, true, nil
	}
}

func (c *Controller) Return(value interface{}) (interface{}, bool, error) {
	c.iterator.statusChan <- &status{
		value: value,
		done:  true,
		err:   nil,
	}

	select {
	case err := <-c.iterator.errorChan:
		c.iterator.isDone() // don't care
		return nil, false, err
	case <-c.iterator.nextChan:
	case <-c.iterator.returnChan:
	}

	c.iterator.isDone() // don't care
	return nil, true, nil
}

func (c *Controller) Error(err error) (interface{}, bool, error) {
	c.iterator.statusChan <- &status{
		value: nil,
		done:  false,
		err:   err,
	}

	select {
	case value := <-c.iterator.nextChan:
		if c.iterator.isDone() {
			return nil, true, nil
		}
		return value, false, nil
	case err := <-c.iterator.errorChan:
		c.iterator.isDone() // don't care
		return nil, false, err
	case <-c.iterator.returnChan:
		c.iterator.isDone() // don't care
		return nil, true, nil
	}
}

func (c *Controller) interceptFirstIteration() bool {
	select {
	case <-c.iterator.nextChan:
	case <-c.iterator.returnChan:
	case <-c.iterator.errorChan:
	}
	return !c.iterator.isDone()
}
