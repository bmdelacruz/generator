package generator

type status struct {
	value interface{}
	done  bool
	err   error
}

func (s status) Data() (interface{}, bool, error) {
	return s.value, s.done, s.err
}

type retStatus interface {
	Type() string
	Data() (interface{}, bool, error)
}

type yieldRetStatus struct {
	value interface{}
}

func (yieldRetStatus) Type() string {
	return "yield"
}

func (rs yieldRetStatus) Data() (interface{}, bool, error) {
	return rs.value, false, nil
}

type errorRetStatus struct {
	err error
}

func (errorRetStatus) Type() string {
	return "error"
}

func (rs errorRetStatus) Data() (interface{}, bool, error) {
	return nil, false, rs.err
}

type returnRetStatus struct {
	value interface{}
}

func (returnRetStatus) Type() string {
	return "return"
}

func (rs returnRetStatus) Data() (interface{}, bool, error) {
	return rs.value, true, nil
}

type firstCall interface {
	Type() string
	Values() (interface{}, error)
}

type returnFirstCall struct {
	value interface{}
}

func (returnFirstCall) Type() string {
	return "return"
}

func (fc returnFirstCall) Values() (interface{}, error) {
	return fc.value, nil
}

type errorFirstCall struct {
	err error
}

func (errorFirstCall) Type() string {
	return "error"
}

func (fc errorFirstCall) Values() (interface{}, error) {
	return nil, fc.err
}
