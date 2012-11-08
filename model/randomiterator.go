package model

import (
	"math/rand"
	"reflect"
	"time"
)

func NewRandomIterator(iterator Iterator) *RandomIterator {
	return &RandomIterator{Iterator: iterator}
}

// RandomIterator stores all values from Iterator in a slice and
// iterates them in random order.
// LessFunc will always be called with a pointer to the struct if Next
// is called with a pointer to a struct or the address of a
// pointer to a stract.
// For all other types LessFunc will be called with the type that
// that the argument of Next points to.
type RandomIterator struct {
	Iterator
	indexedSliceIterator *IndexedSliceIterator
}

func (self *RandomIterator) Next(resultRef interface{}) bool {
	if self.Err() != nil {
		return false
	}
	if self.indexedSliceIterator == nil {
		resultType := reflect.ValueOf(resultRef).Elem().Type()
		resultKind := resultType.Kind()
		slice := make([]interface{}, 0, 16)
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		indices := r.Perm(len(slice))
		self.indexedSliceIterator = NewIndexedSliceIterator(slice, indices)
	}
	return self.indexedSliceIterator.Next(resultRef)
}
