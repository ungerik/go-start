package model

type IteratorFunc func() interface{}

func (self IteratorFunc) Next() interface{} {
	return self()
}

func (self IteratorFunc) Err() error {
	return nil
}
