package model

import "reflect"

func NewSliceIterator(objects ...interface{}) *SliceIterator {
	return &SliceIterator{slice: objects}
}

func NewSliceOrErrorOnlyIterator(slice interface{}, err error) Iterator {
	if err != nil {
		return NewErrorOnlyIterator(err)
	}
	return NewSliceIterator(slice)
}

type SliceIterator struct {
	slice []interface{}
	index int
}

func (self *SliceIterator) Next(resultPtr interface{}) bool {
	if self.index >= len(self.slice) {
		return false
	}
	object := self.slice[self.index]
	self.index++
	reflect.ValueOf(resultPtr).Elem().Set(reflect.ValueOf(object))
	return true
}

func (self *SliceIterator) Err() error {
	return nil
}
