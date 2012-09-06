package mongomedia

import (
	"github.com/ungerik/go-start/model"
)

type imageIterator struct {
	iter model.Iterator
}

func (self imageIterator) Next() interface{} {
	if doc := self.iter.Next(); doc != nil {
		return doc.(*ImageDoc).GetAndInitImage()
	}
	return nil
}

func (self imageIterator) Err() error {
	return self.iter.Err()
}
