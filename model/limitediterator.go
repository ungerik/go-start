package model

func NewLimitedIterator(iter Iterator, limit int) *LimitedIterator {
	return &LimitedIterator{Iterator: iter, limit: limit}
}

///////////////////////////////////////////////////////////////////////////////
// LimitedIterator

type LimitedIterator struct {
	Iterator
	limit int
	index int
}

func (self *LimitedIterator) Next(resultRef interface{}) bool {
	if self.index >= self.limit {
		return false
	}
	self.index++
	return self.Iterator.Next(resultRef)
}
