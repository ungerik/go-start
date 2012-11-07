package model

import "github.com/ungerik/go-start/reflection"

func NewSliceIterator(objects ...interface{}) *SliceIterator {
	return &SliceIterator{slice: objects}
}

func NewSliceOrErrorOnlyIterator(slice interface{}, err error) Iterator {
	if err != nil {
		return NewErrorOnlyIterator(err)
	}
	return NewSliceIterator(slice)
}

// SliceIterator
// When calling Next, resultPtr must be a pointer to the slice element type
type SliceIterator struct {
	slice []interface{}
	index int
}

func (self *SliceIterator) Next(resultPtr interface{}) bool {
	if self.index >= len(self.slice) {
		return false
	}
	reflection.AssignToResultPtr(self.slice[self.index], resultPtr)
	self.index++
	return true
}

func (self *SliceIterator) Err() error {
	return nil
}
