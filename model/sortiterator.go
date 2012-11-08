package model

import (
	"reflect"

	"github.com/ungerik/go-start/utils"
)

type LessFunc func(a, b interface{}) (less bool)

func NewSortIterator(iterator Iterator, lessFunc LessFunc) *SortIterator {
	return &SortIterator{Iterator: iterator, LessFunc: lessFunc}
}

// SortIterator stores all values from Iterator in a slice and
// iterates them sorted by LessFunc.
// LessFunc will always be called with a pointer to the struct if Next
// is called with a pointer to a struct or the address of a
// pointer to a stract.
// For all other types LessFunc will be called with the type that
// that the argument of Next points to.
type SortIterator struct {
	Iterator
	LessFunc      LessFunc
	sliceIterator *SliceIterator
}

func (self *SortIterator) Next(resultRef interface{}) bool {
	if self.Err() != nil {
		return false
	}
	if self.sliceIterator == nil {
		resultType := reflect.ValueOf(resultRef).Elem().Type()
		resultKind := resultType.Kind()
		slice := []interface{}{}
		for self.Iterator.Next(resultRef) {
			resultVal := reflect.ValueOf(resultRef).Elem()
			if resultKind == reflect.Struct {
				resultCopy := reflect.New(resultType)
				resultCopy.Elem().Set(resultVal)
				slice = append(slice, resultCopy.Interface())
			} else {
				slice = append(slice, resultVal.Interface())
			}
		}
		if self.Err() != nil {
			return false
		}
		utils.Sort(slice, self.LessFunc)
		self.sliceIterator = NewSliceIterator(slice...)
	}
	return self.sliceIterator.Next(resultRef)
}
