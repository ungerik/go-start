package model

type FilterFunc func(resultPtr interface{}) bool

type FilterIterator struct {
	Iterator
	PassFilter FilterFunc
}

func (self *FilterIterator) Next(resultPtr interface{}) bool {
	for self.Iterator.Next(resultPtr) {
		if self.PassFilter(resultPtr) {
			return true
		}
	}
	return false
}
