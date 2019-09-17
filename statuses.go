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

func (rs yieldRetStatus) Type() string {
	return "yield"
}

func (rs yieldRetStatus) Data() (interface{}, bool, error) {
	return rs.value, false, nil
}

type errorRetStatus struct {
	err error
}

func (rs errorRetStatus) Type() string {
	return "error"
}

func (rs errorRetStatus) Data() (interface{}, bool, error) {
	return nil, false, rs.err
}

type returnRetStatus struct {
	value interface{}
}

func (rs returnRetStatus) Type() string {
	return "return"
}

func (rs returnRetStatus) Data() (interface{}, bool, error) {
	return rs.value, true, nil
}
