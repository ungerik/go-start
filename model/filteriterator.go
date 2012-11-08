package model

type FilterFunc func(resultRef interface{}) bool

type FilterIterator struct {
	Iterator
	PassFilter FilterFunc
}

func (self *FilterIterator) Next(resultRef interface{}) bool {
	for self.Iterator.Next(resultRef) {
		if self.PassFilter(resultRef) {
			return true
		}
	}
	return false
}
