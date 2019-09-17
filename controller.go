package generator

type Controller struct {
	g *Generator

	wasUsed bool
}

func (c *Controller) Yield(value interface{}) (interface{}, bool, error) {
	if !c.wasUsed {
		c.wasUsed = true
	}

	select {
	case value := <-c.g.unhandledReturnChan:
		return value, true, nil
	case err := <-c.g.unhandledErrorChan:
		c.g.statusChan <- &status{
			value: nil,
			done:  false,
			err:   nil,
		}

		select {
		case <-c.g.yieldChan:
		case <-c.g.errorChan:
		case <-c.g.returnChan:
		}

		c.g.updateAndGetIsDone()

		return nil, false, err
	default:
	}

	c.g.statusChan <- &status{
		value: value,
		done:  false,
		err:   nil,
	}

	select {
	case value := <-c.g.yieldChan:
		if c.g.updateAndGetIsDone() {
			return nil, true, nil
		}
		return value, false, nil
	case err := <-c.g.errorChan:
		c.g.updateAndGetIsDone()
		return nil, false, err
	case <-c.g.returnChan:
		c.g.updateAndGetIsDone()
		return nil, true, nil
	}
}

func (c *Controller) Error(err error) (interface{}, bool, error) {
	if !c.wasUsed {
		c.wasUsed = true
	}

	select {
	case value := <-c.g.unhandledReturnChan:
		return value, true, nil
	case err := <-c.g.unhandledErrorChan:
		c.g.statusChan <- &status{
			value: nil,
			done:  false,
			err:   nil,
		}

		select {
		case <-c.g.yieldChan:
		case <-c.g.errorChan:
		case <-c.g.returnChan:
		}

		c.g.updateAndGetIsDone()

		return nil, false, err
	default:
	}

	c.g.statusChan <- &status{
		value: nil,
		done:  false,
		err:   err,
	}

	select {
	case value := <-c.g.yieldChan:
		if c.g.updateAndGetIsDone() {
			return nil, true, nil
		}
		return value, false, nil
	case err := <-c.g.errorChan:
		c.g.updateAndGetIsDone()
		return nil, false, err
	case <-c.g.returnChan:
		c.g.updateAndGetIsDone()
		return nil, true, nil
	}
}
