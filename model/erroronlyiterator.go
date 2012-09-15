package model

func NewErrorOnlyIterator(err error) Iterator {
	return &ErrorOnlyIterator{err}
}

///////////////////////////////////////////////////////////////////////////////
// ErrorOnlyIterator

type ErrorOnlyIterator struct {
	err error
}

func (self *ErrorOnlyIterator) Next() interface{} {
	return nil
}

func (self *ErrorOnlyIterator) Err() error {
	return self.err
}
