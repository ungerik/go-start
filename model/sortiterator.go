package model

import (
	"github.com/ungerik/go-start/reflection"
	"reflect"
)

func NewSortIterator(iterator Iterator, compareFunc interface{}) *SortIterator {
	return &SortIterator{Iterator: iterator, CompareFunc: compareFunc}
}

type SortIterator struct {
	Iterator      Iterator
	CompareFunc   interface{}
	sliceIterator *SliceIterator
	err           error
}

func (self *SortIterator) Next(resultRef interface{}) bool {
	if self.err != nil {
		return false
	}

	if self.sliceIterator == nil {
		f, err := reflection.NewSortCompareFunc(self.CompareFunc)
		if err != nil {
			self.err = err
			return false
		}

		slice := make([]interface{}, 0, 16)
		result := reflect.New(f.ArgType).Interface()
		for self.Iterator.Next(result) {
			slice = append(slice, result)
			result = reflect.New(f.ArgType).Interface()
		}
		self.err = self.Iterator.Err()
		if self.err != nil {
			return false
		}

		self.err = f.Sort(slice)
		if self.err != nil {
			return false
		}

		self.sliceIterator = NewSliceIterator(slice)
	}

	return self.sliceIterator.Next(resultRef)
}

func (self *SortIterator) Err() error {
	return self.err
}
