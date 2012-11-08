package model

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/reflection"
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

func (self *IndexedSliceIterator) Next(resultRef interface{}) bool {
	if self.err != nil || self.index >= len(self.indices) {
		return false
	}
	if self.indices[self.index] >= len(self.slice) {
		self.err = errs.Format("Index %d from indices greater or equal than length of slice %d", self.indices[self.index], len(self.slice))
		return false
	}
	reflection.SmartCopy(self.slice[self.indices[self.index]], resultRef)
	self.index++
	return false
}

func (self *IndexedSliceIterator) Err() error {
	return self.err
}
