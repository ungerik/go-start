package model

import "reflect"

// ConvertIterator returns an Iterator that calls conversionFunc
// for every non nil from.Next() result and returns the result
// of conversionFunc for every Next().
func ConversionIterator(from Iterator, sourceResult interface{}, conversionFunc func(interface{}) interface{}) Iterator {
	return &conversionIterator{from, sourceResult, conversionFunc}
}

type conversionIterator struct {
	Iterator
	SourceResultPtr interface{}
	conversionFunc  func(interface{}) interface{}
}

func (self *conversionIterator) Next(resultPtr interface{}) bool {
	if !self.Iterator.Next(self.SourceResultPtr) {
		return false
	}
	destinationResultPtr := self.conversionFunc(self.SourceResultPtr)
	reflect.ValueOf(resultPtr).Elem().Set(reflect.ValueOf(destinationResultPtr).Elem())
	return true
}
