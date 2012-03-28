package model

///////////////////////////////////////////////////////////////////////////////
// EmptyIterator

type EmptyIterator struct {
}

func (self *EmptyIterator) Next() interface{} {
	return nil
}

func (self *EmptyIterator) Err() error {
	return nil
}
