package model

import (
	"math/rand"
	"reflect"
	"time"
)

func NewRandomIterator(iterator Iterator) *RandomIterator {
	return &RandomIterator{Iterator: iterator}
}

type RandomIterator struct {
	Iterator
	indexedSliceIterator *IndexedSliceIterator
}

func (self *RandomIterator) Next(resultPtr interface{}) bool {
	if self.Err() != nil {
		return false
	}
	if self.indexedSliceIterator == nil {
		slice := []interface{}{}
		for self.Iterator.Next(resultPtr) {
			slice = append(slice, reflect.ValueOf(resultPtr).Elem().Interface())
		}
		if self.Err() != nil {
			return false
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		indices := r.Perm(len(slice))
		self.indexedSliceIterator = NewIndexedSliceIterator(slice, indices)
	}
	return self.indexedSliceIterator.Next(resultPtr)
}
