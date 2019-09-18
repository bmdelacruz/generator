package generator

type Controller struct {
	g *Generator

	wasUsed bool
}

func (c *Controller) Yield(value interface{}) (interface{}, bool, error) {
	return c.sendAndReceive(
		&status{
			value: value,
			done:  false,
			err:   nil,
		},
	)
}

func (c *Controller) Error(err error) (interface{}, bool, error) {
	return c.sendAndReceive(
		&status{
			value: nil,
			done:  false,
			err:   err,
		},
	)
}

func (c *Controller) sendAndReceive(statusToSend *status) (interface{}, bool, error) {
	if !c.wasUsed {
		c.wasUsed = true
	}

	select {
	case fc, ok := <-c.g.firstCallChan:
		if ok {
			switch fc.Type() {
			case "return":
				value, _ := fc.Values()
				return value, true, nil
			case "error":
				_, err := fc.Values()
				c.g.statusChan <- &status{
					value: nil,
					done:  false,
					err:   nil,
				}
				<-c.g.retStatusChan
				<-c.g.isDoneChan
				return nil, false, err
			}
		}
	default:
	}

	if c.g.isDone {
		return nil, true, nil
	}

	c.g.statusChan <- statusToSend
	rs := <-c.g.retStatusChan
	<-c.g.isDoneChan
	return rs.Data()
}
