package model

import (
// "fmt"
// "reflect"
)

///////////////////////////////////////////////////////////////////////////////
// Iterator

// Iteration stops with Next() == nil, check Err() afterwards.
//
// Be careful when using the address of a local variable for resultRef
// in a loop that creates a closure function that will be called later for rendering.
// Go does not make copies of variables that are bound to closure, but references them.
// If a closure is called outside of the loop, it will always reference the
// same variable with the value from the last Next() call in it.
type Iterator interface {
	// Next tries to unmarshalles the next iteration's data into resultRef
	// and returns a bool indicating if it was successful.
	// After Next() has returned false, check the result of Err().
	Next(resultRef interface{}) bool

	// Err returns the error of the last iteration.
	// Err returns nil if the last iteration was successful,
	// or there was no last iteration.
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
