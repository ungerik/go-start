package model

import "reflect"

// ConvertIterator returns an Iterator that calls conversionFunc
// for every non nil from.Next() result and returns the result
// of conversionFunc for every Next().
func ConversionIterator(from Iterator, conversionFunc func(interface{}) interface{}) Iterator {
	return &conversionIterator{from, conversionFunc}
}

type conversionIterator struct {
	Iterator
	conversionFunc func(interface{}) interface{}
}

func (self *conversionIterator) Next(resultPtr interface{}) bool {
	var doc interface{}
	if !self.Iterator.Next(&doc) {
		return false
	}
	doc = self.conversionFunc(doc)
	reflect.ValueOf(resultPtr).Elem().Set(reflect.ValueOf(doc))
	return true
}
