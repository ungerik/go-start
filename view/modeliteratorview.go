package view

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// ModelIteratorView

type GetModelIteratorFunc func(ctx *Context) model.Iterator
type GetModelIteratorViewFunc func(model interface{}, ctx *Context) (view View, err error)

func ModelIterator(iter model.Iterator) GetModelIteratorFunc {
	return func(ctx *Context) model.Iterator {
		return iter
	}
}

type ModelIteratorView struct {
	ViewBase
	GetModelIterator     GetModelIteratorFunc
	GetModelIteratorView GetModelIteratorViewFunc // nil Views will be ignored
}

func (self *ModelIteratorView) Render(ctx *Context) (err error) {
	var children Views

	iter := self.GetModelIterator(ctx)
	for model := iter.Next(); model != nil; model = iter.Next() {
		view, err := self.GetModelIteratorView(model, ctx)
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

	return children.Render(ctx)
}
