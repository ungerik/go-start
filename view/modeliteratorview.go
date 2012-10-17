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
	// Model will be used by Next() from the model iterator
	Model                interface{}
	GetModelIteratorView GetModelIteratorViewFunc // nil Views will be ignored
}

func (self *ModelIteratorView) Render(ctx *Context) (err error) {
	if self.Model == nil {
		panic("view.ModelIteratorTableView.RowModel is nil")
	}

	iter := self.GetModelIterator(ctx)
	for iter.Next(self.Model) {
		view, err := self.GetModelIteratorView(ctx, self.Model)
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
	}
	return iter.Err()
}
