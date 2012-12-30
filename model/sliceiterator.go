package model

import (
	"fmt"
	"reflect"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/reflection"
)

func NewSliceIterator(slice interface{}) *SliceIterator {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic(fmt.Errorf("Expected slice or array, got %T", slice))
	}
	return &SliceIterator{slice: v}
}

func NewSliceOrErrorOnlyIterator(slice interface{}, err error) Iterator {
	if err != nil {
		return NewErrorOnlyIterator(err)
	}
	return NewSliceIterator(slice)
}

// SliceIterator
type SliceIterator struct {
	slice reflect.Value
	index int
}

func (self *SliceIterator) Next(resultRef interface{}) bool {
	if self.index >= self.slice.Len() {
		return false
	}
	reflection.SmartCopy(self.slice.Index(self.index).Interface(), resultRef)
	self.index++
	return true
}

func (self *SliceIterator) Err() error {
	return nil
}
