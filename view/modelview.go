package view

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// ModelView

type GetModelIteratorFunc func(response *Response) model.Iterator
type GetModelViewFunc func(model interface{}, response *Response) (view View, err error)

func ModelIterator(iter model.Iterator) GetModelIteratorFunc {
	return func(response *Response) model.Iterator {
		return iter
	}
}

type ModelView struct {
	ViewBase
	GetModelIterator GetModelIteratorFunc
	GetModelView     GetModelViewFunc // nil Views will be ignored
}

func (self *ModelView) Render(response *Response) (err error) {
	var children Views

	iter := self.GetModelIterator(response)
	for model := iter.Next(); model != nil; model = iter.Next() {
		view, err := self.GetModelView(model, response)
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
