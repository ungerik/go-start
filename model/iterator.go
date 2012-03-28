package model

import (
	"github.com/ungerik/go-start/utils"
	"math/rand"
	"time"
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
