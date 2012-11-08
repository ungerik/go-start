package model

func NewErrorOnlyIterator(err error) Iterator {
	return &ErrorOnlyIterator{err}
}

type ErrorOnlyIterator struct {
	err error
}

func (self *ErrorOnlyIterator) Next(resultRef interface{}) bool {
	return false
}

func (self *ErrorOnlyIterator) Err() error {
	return self.err
}
