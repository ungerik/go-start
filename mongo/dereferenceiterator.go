package mongo

// Returns an iterator of dereferenced refs, or an error iterator if there was an error
func NewDereferenceIterator(refs ...Ref) *DereferenceIterator {
	return &DereferenceIterator{Refs: refs}
}

type DereferenceIterator struct {
	Refs  []Ref
	index int
	err   error
}

func (self *DereferenceIterator) Next(resultRef interface{}) bool {
	if self.err != nil || self.index >= len(self.Refs) {
		return false
	}
	self.err = self.Refs[self.index].Get(resultRef)
	self.index++
	return self.err == nil
}

func (self *DereferenceIterator) Err() error {
	return self.err
}
