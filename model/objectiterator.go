package model

func NewObjectIterator(objects ...interface{}) Iterator {
	return &ObjectIterator{objects: objects}
}

func NewObjectOrErrorIterator(object interface{}, err error) Iterator {
	if err != nil {
		return NewErrorIterator(err)
	}
	return NewObjectIterator(object)
}

///////////////////////////////////////////////////////////////////////////////
// ObjectIterator

type ObjectIterator struct {
	objects []interface{}
	index   int
}

func (self *ObjectIterator) Next() interface{} {
	if self.index >= len(self.objects) {
		return nil
	}
	object := self.objects[self.index]
	self.index++
	return object
}

func (self *ObjectIterator) Err() error {
	return nil
}
