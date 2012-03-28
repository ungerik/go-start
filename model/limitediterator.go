package model

func NewLimitedIterator(iter Iterator, limit int) *LimitedIterator {
	return &LimitedIterator{iter: iter, limit: limit}
}

///////////////////////////////////////////////////////////////////////////////
// LimitedIterator

type LimitedIterator struct {
	iter  Iterator
	limit int
	index int
}

func (self *LimitedIterator) Next() interface{} {
	if self.index >= self.limit {
		return nil
	}
	self.index++
	return self.iter.Next()
}

func (self *LimitedIterator) Err() error {
	return self.iter.Err()
}
