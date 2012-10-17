package model

import (
	"reflect"

	"github.com/ungerik/go-start/errs"
)

func NewIndexedSliceIterator(slice []interface{}, indices []int) *IndexedSliceIterator {
	return &IndexedSliceIterator{slice: slice, indices: indices}
}

type IndexedSliceIterator struct {
	slice   []interface{}
	indices []int
	index   int
	err     error
}

func (self *IndexedSliceIterator) Next(resultPtr interface{}) bool {
	if self.err != nil || self.index >= len(self.indices) {
		return false
	}
	if self.indices[self.index] >= len(self.slice) {
		self.err = errs.Format("Index %d from indices greater or equal than length of slice %d", self.indices[self.index], len(self.slice))
		return false
	}
	v := reflect.ValueOf(self.slice[self.indices[self.index]])
	self.index++
	resultVal := reflect.ValueOf(resultPtr).Elem()
	if resultVal.Type() == v.Type() {
		resultVal.Set(v)
	} else {
		resultVal.Set(v.Elem())
	}
	return false
}

func (self *IndexedSliceIterator) Err() error {
	return self.err
}
