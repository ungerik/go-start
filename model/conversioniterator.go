package model

import "github.com/ungerik/go-start/reflection"

// ConvertIterator returns an Iterator that calls conversionFunc
// for every from.Next() result and returns the result
// of conversionFunc at every Next().
func ConversionIterator(from Iterator, fromResultPtr interface{}, conversionFunc func(interface{}) interface{}) Iterator {
	return &conversionIterator{from, fromResultPtr, conversionFunc}
}

type conversionIterator struct {
	Iterator
	FromResultPtr  interface{}
	conversionFunc func(interface{}) interface{}
}

func (self *conversionIterator) Next(resultRef interface{}) bool {
	if !self.Iterator.Next(self.FromResultPtr) {
		return false
	}
	reflection.SmartCopy(self.conversionFunc(self.FromResultPtr), resultRef)
	return true
}
