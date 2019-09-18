package generator

// Controller provides functions that will control the generator
// associated with its instance.
type Controller struct {
	g *Generator

	// wasUsed is equal to true if any of the functions of this
	// controller was used
	wasUsed bool
}

// Yield sends the value to the consumer of the generator and then
// waits for the next generator function invocation that will get
// the data that will be returned by this function.
//
// When the previous call already returned shouldReturn equal to true,
// the current and the succeeding calls will return (<nil>, true, <nil>).
//
// Returns ([value], [shouldReturn], [error])
func (c *Controller) Yield(value interface{}) (interface{}, bool, error) {
	return c.sendAndReceive(
		&status{
			value: value,
			done:  false,
			err:   nil,
		},
	)
}

// Error sends an error to the consumer of the generator and then
// waits for the next generator function invocation that will get
// the data that will be returned by this function.
//
// When the previous call already returned shouldReturn equal to true,
// the current and the succeeding calls will return (<nil>, true, <nil>).
//
// Returns ([value], [shouldReturn], [error])
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
		// mark that any of the controller function has been used
		c.wasUsed = true
	}

	select {
	// if there is a saved error or return value earlier, receive it
	case fc, ok := <-c.g.firstCallChan:
		if ok {
			switch fc.Type() {
			case "return":
				// there's no need to send and receive here since there will
				// be no more succeeding generator function calls that will
				// be sending values to them

				// just return the saved value
				value, _ := fc.Values()
				return value, true, nil
			case "error":
				_, err := fc.Values()

				// the generator controller function needs to be overridden
				// since there is a pending error that was sent by the consumer
				// of the generator
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
