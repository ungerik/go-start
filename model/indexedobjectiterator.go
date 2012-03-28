package model

import "github.com/ungerik/go-start/errs"

func NewIndexedObjectIterator(objects []interface{}, indices []int) Iterator {
	return &IndexedObjectIterator{objects: objects, indices: indices}
}

///////////////////////////////////////////////////////////////////////////////
// IndexedObjectIterator

type IndexedObjectIterator struct {
	objects []interface{}
	indices []int
	index   int
	err     error
}

func (self *IndexedObjectIterator) Next() interface{} {
	if self.err != nil || self.index >= len(self.indices) {
		return nil
	}
	if self.indices[self.index] >= len(self.objects) {
		self.err = errs.Format("Index %d from indices slice greater or equal than length of objects slice %d", self.indices[self.index], len(self.objects))
		return nil
	}
	object := self.objects[self.indices[self.index]]
	self.index++
	return object
}

func (self *IndexedObjectIterator) Err() error {
	return self.err
}
