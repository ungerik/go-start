package model

import "github.com/ungerik/go-start/reflection"

func NewSingletonIterator(singleton interface{}) *SingletonIterator {
	return &SingletonIterator{Singleton: singleton}
}

func NewSingletonOrErrorOnlyIterator(singleton interface{}, err error) Iterator {
	if err != nil {
		return NewErrorOnlyIterator(err)
	}
	return NewSingletonIterator(singleton)
}

type SingletonIterator struct {
	Singleton interface{}
	iterated  bool
}

func (self *SingletonIterator) Next(resultRef interface{}) bool {
	if self.iterated || self.Singleton == nil {
		return false
	}
	reflection.SmartCopy(self.Singleton, resultRef)
	self.iterated = true
	return true
}

func (self *SingletonIterator) Err() error {
	return nil
}
