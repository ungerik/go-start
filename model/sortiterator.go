package model

import (
	"reflect"

	"github.com/ungerik/go-start/utils"
)

type LessFunc func(a, b interface{}) (less bool)

func NewSortIterator(iterator Iterator, lessFunc LessFunc) *SortIterator {
	return &SortIterator{Iterator: iterator, LessFunc: lessFunc}
}

type SortIterator struct {
	Iterator
	LessFunc      LessFunc
	sliceIterator *SliceIterator
}

func (self *SortIterator) Next(resultPtr interface{}) bool {
	if self.Err() != nil {
		return false
	}
	if self.sliceIterator == nil {
		slice := []interface{}{}
		for self.Iterator.Next(resultPtr) {
			slice = append(slice, reflect.ValueOf(resultPtr).Elem().Interface())
		}
		if self.Err() != nil {
			return false
		}
		utils.Sort(slice, self.LessFunc)
		self.sliceIterator = NewSliceIterator(slice...)
	}
	return self.sliceIterator.Next(resultPtr)
}
