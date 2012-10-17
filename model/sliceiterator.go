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
	v := reflect.ValueOf(self.slice[self.index])
	self.index++
	resultVal := reflect.ValueOf(resultPtr).Elem()
	if resultVal.Type() == v.Type() {
		resultVal.Set(v)
	} else {
		resultVal.Set(v.Elem())
	}
	return true
}

func (self *SliceIterator) Err() error {
	return nil
}
