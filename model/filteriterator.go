package model

import "reflect"

type FilterFunc func(doc interface{}) (ok bool)

type FilterIterator struct {
	Iterator
	PassFilter FilterFunc
}

func (self *FilterIterator) Next(resultPtr interface{}) bool {
	var doc interface{}
	for self.Iterator.Next(&doc) {
		if self.PassFilter(doc) {
			reflect.ValueOf(resultPtr).Elem().Set(reflect.ValueOf(doc))
			return true
		}
	}
	return false
}
