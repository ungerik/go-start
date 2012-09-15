package model

///////////////////////////////////////////////////////////////////////////////
// FilterIterator

type FilterFunc func(doc interface{}) (ok bool)

type FilterIterator struct {
	Iterator
	PassFilter FilterFunc
}

func (self *FilterIterator) Next() interface{} {
	for doc := self.Iterator.Next(); doc != nil; doc = self.Iterator.Next() {
		if self.PassFilter(doc) {
			return doc
		}
	}
	return nil
}
