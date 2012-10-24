package view

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// ModelIteratorView

type GetModelIteratorFunc func(ctx *Context) model.Iterator
type GetModelIteratorViewFunc func(ctx *Context, model interface{}) (view View, err error)

func ModelIterator(iter model.Iterator) GetModelIteratorFunc {
	return func(ctx *Context) model.Iterator {
		return iter
	}
}

type ModelIteratorView struct {
	ViewBase
	GetModelIterator GetModelIteratorFunc
	// GetModel returns the model used GetModelIteratorView.
	// Only return the same object if you know exactly what you are doing.
	// To be on the safe side, return a new model instance at every call.
	GetModel     func(ctx *Context) (model interface{}, err error)
	GetModelView GetModelIteratorViewFunc // nil Views will be ignored
}

func (self *ModelIteratorView) Render(ctx *Context) (err error) {
	iter := self.GetModelIterator(ctx)
	model, err := self.GetModel(ctx)
	if err != nil {
		return err
	}
	for iter.Next(model) {
		view, err := self.GetModelView(ctx, model)
		if err != nil {
			return err
		}
		if view != nil {
			view.Init(view)
			err = view.Render(ctx)
			if err != nil {
				return err
			}
		}
		model, err = self.GetModel(ctx)
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
