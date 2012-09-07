package model

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// Iterator

// Iteration stops with Next() == nil, check Err() afterwards
type Iterator interface {
	Next() interface{}
	Err() error
}

///////////////////////////////////////////////////////////////////////////////
// Functions

func Iterate(i Iterator, callback interface{}) error {
	v := reflect.ValueOf(callback)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("model.Iterate: callback must be a function, got %s ", t))
	}
	if t.NumIn() != 1 {
		panic(fmt.Errorf("model.Iterate: callback must have one argument, got %d arguments", t.NumIn()))
	}
	if t.NumOut() != 0 {
		panic(fmt.Errorf("model.Iterate: callback must not have a result, got %d results", t.NumOut()))
	}
	for doc := i.Next(); doc != nil; doc = i.Next() {
		v.Call([]reflect.Value{reflect.ValueOf(doc)})
	}
	return i.Err()
}

// ConvertIterator returns an Iterator that calls convertFunc
// for every non nil from.Next() result and returns the result
// of convertFunc for every Next().
func ConvertIterator(from Iterator, convertFunc func(interface{}) interface{}) Iterator {
	return &convertIterator{from, convertFunc}
}

type convertIterator struct {
	Iterator
	convertFunc func(interface{}) interface{}
}

func (self *convertIterator) Next() interface{} {
	doc := self.Iterator.Next()
	if doc == nil {
		return nil
	}
	return self.convertFunc(doc)
}

// Channels will be closed after last result
func IterateToChannel(i Iterator) (docs <-chan interface{}, errs <-chan error) {
	docChan := make(chan interface{}, 32)
	errChan := make(chan error, 1)
	go func() {
		defer close(docChan)
		defer close(errChan)
		for doc := i.Next(); doc != nil; doc = i.Next() {
			docChan <- doc
		}
		if i.Err() != nil {
			errChan <- i.Err()
		}
	}()
	return docChan, errChan
}

func SortIterator(i Iterator, lessFunc func(a, b interface{}) (less bool)) Iterator {
	var slice []interface{}
	for doc := i.Next(); doc != nil; doc = i.Next() {
		slice = append(slice, doc)
	}
	if i.Err() != nil {
		return i
	}
	utils.Sort(slice, lessFunc)
	return NewObjectIterator(slice...)
}

func RandomIterator(i Iterator) Iterator {
	var slice []interface{}
	for doc := i.Next(); doc != nil; doc = i.Next() {
		slice = append(slice, doc)
	}
	if i.Err() != nil {
		return i
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	indices := r.Perm(len(slice))
	return NewIndexedObjectIterator(slice, indices)
}
