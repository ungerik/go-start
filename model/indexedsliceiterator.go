package model

import (
	"fmt"
	"reflect"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/reflection"
)

func NewIndexedSliceIterator(slice interface{}, indices []int) *IndexedSliceIterator {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic(fmt.Errorf("Expected slice or array, got %T", slice))
	}
	return &IndexedSliceIterator{slice: v, indices: indices}
}

type IndexedSliceIterator struct {
	slice   reflect.Value
	indices []int
	index   int
	err     error
}

func (self *IndexedSliceIterator) Next(resultRef interface{}) bool {
	if self.err != nil || self.index >= len(self.indices) {
		return false
	}
	if self.indices[self.index] >= self.slice.Len() {
		self.err = errs.Format("Index %d from indices greater or equal than length of slice %d", self.indices[self.index], self.slice.Len())
		return false
	}
	reflection.SmartCopy(self.slice.Index(self.indices[self.index]), resultRef)
	self.index++
	return false
}

func (self *IndexedSliceIterator) Err() error {
	return self.err
}
