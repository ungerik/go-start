package model

import (
	// "fmt"
	// "reflect"
	"math/rand"
	"time"

	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// Iterator

// Iteration stops with Next() == nil, check Err() afterwards
type Iterator interface {
	Next(resultPtr interface{}) bool
	Err() error
}

///////////////////////////////////////////////////////////////////////////////
// Functions

// func Iterate(i Iterator, callback interface{}) error {
// 	v := reflect.ValueOf(callback)
// 	t := v.Type()
// 	if t.Kind() != reflect.Func {
// 		panic(fmt.Errorf("model.Iterate: callback must be a function, got %s ", t))
// 	}
// 	if t.NumIn() != 1 {
// 		panic(fmt.Errorf("model.Iterate: callback must have one argument, got %d arguments", t.NumIn()))
// 	}
// 	if t.NumOut() != 0 {
// 		panic(fmt.Errorf("model.Iterate: callback must not have a result, got %d results", t.NumOut()))
// 	}
// 	for doc := i.Next(); doc != nil; doc = i.Next() {
// 		v.Call([]reflect.Value{reflect.ValueOf(doc)})
// 	}
// 	return i.Err()
// }

// // Channels will be closed after last result
// func IterateToChannel(i Iterator) (docs <-chan interface{}, errs <-chan error) {
// 	docChan := make(chan interface{}, 32)
// 	errChan := make(chan error, 1)
// 	go func() {
// 		defer close(docChan)
// 		defer close(errChan)
// 		for doc := i.Next(); doc != nil; doc = i.Next() {
// 			docChan <- doc
// 		}
// 		if i.Err() != nil {
// 			errChan <- i.Err()
// 		}
// 	}()
// 	return docChan, errChan
// }

func SortIterator(i Iterator, lessFunc func(a, b interface{}) (less bool)) Iterator {
	var slice []interface{}
	var doc interface{}
	for i.Next(&doc) {
		slice = append(slice, doc)
	}
	if i.Err() != nil {
		return i
	}
	utils.Sort(slice, lessFunc)
	return NewSliceIterator(slice...)
}

func RandomIterator(i Iterator) Iterator {
	var slice []interface{}
	var doc interface{}
	for i.Next(&doc) {
		slice = append(slice, doc)
	}
	if i.Err() != nil {
		return i
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	indices := r.Perm(len(slice))
	return NewIndexedSliceIterator(slice, indices)
}
