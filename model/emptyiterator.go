package model

type EmptyIterator struct {
}

func (self *EmptyIterator) Next(resultPtr interface{}) bool {
	return false
}

func (self *EmptyIterator) Err() error {
	return nil
}
