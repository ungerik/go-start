package view

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// ModelIteratorView

type GetModelIteratorFunc func(response *Response) model.Iterator
type GetModelIteratorViewFunc func(model interface{}, response *Response) (view View, err error)

func ModelIterator(iter model.Iterator) GetModelIteratorFunc {
	return func(response *Response) model.Iterator {
		return iter
	}
}

type ModelIteratorView struct {
	ViewBase
	GetModelIterator     GetModelIteratorFunc
	GetModelIteratorView GetModelIteratorViewFunc // nil Views will be ignored
}

func (self *ModelIteratorView) Render(response *Response) (err error) {
	var children Views

	iter := self.GetModelIterator(response)
	for model := iter.Next(); model != nil; model = iter.Next() {
		view, err := self.GetModelIteratorView(model, response)
		if err != nil {
			return err
		}
		if view != nil {
			children = append(children, view)
			view.Init(view)
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}

	return children.Render(response)
}
