package model

func NewErrorIterator(err error) Iterator {
	return &ErrorIterator{err}
}

///////////////////////////////////////////////////////////////////////////////
// ErrorIterator

type ErrorIterator struct {
	err error
}

func (self *ErrorIterator) Next() interface{} {
	return nil
}

func (self *ErrorIterator) Err() error {
	return self.err
}
