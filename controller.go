package generator

type Controller struct {
	g *Generator

	wasUsed bool
}

func (c *Controller) Yield(value interface{}) (interface{}, bool, error) {
	if !c.wasUsed {
		c.wasUsed = true
	}

	c.g.statusChan <- &status{
		value: value,
		done:  false,
		err:   nil,
	}

	select {
	case value := <-c.g.yieldChan:
		if c.g.isDone() {
			return nil, true, nil
		}
		return value, false, nil
	case err := <-c.g.unhandledErrorChan:
		c.g.isDone() // don't care
		return nil, false, err
	case err := <-c.g.errorChan:
		c.g.isDone() // don't care
		return nil, false, err
	case <-c.g.returnChan:
		c.g.isDone() // don't care
		return nil, true, nil
	}
}

func (c *Controller) Error(err error) (interface{}, bool, error) {
	if !c.wasUsed {
		c.wasUsed = true
	}

	c.g.statusChan <- &status{
		value: nil,
		done:  false,
		err:   err,
	}

	select {
	case value := <-c.g.yieldChan:
		if c.g.isDone() {
			return nil, true, nil
		}
		return value, false, nil
	case err := <-c.g.unhandledErrorChan:
		c.g.isDone() // don't care
		return nil, false, err
	case err := <-c.g.errorChan:
		c.g.isDone() // don't care
		return nil, false, err
	case <-c.g.returnChan:
		c.g.isDone() // don't care
		return nil, true, nil
	}
}
