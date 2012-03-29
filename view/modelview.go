package view

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// ModelView

type GetModelIteratorFunc func(context *Context) model.Iterator
type GetModelViewFunc func(model interface{}, context *Context) (view View, err error)

func ModelIterator(iter model.Iterator) GetModelIteratorFunc {
	return func(context *Context) model.Iterator {
		return iter
	}
}

type ModelView struct {
	ViewBase
	GetModelIterator GetModelIteratorFunc
	GetModelView     GetModelViewFunc // nil Views will be ignored
}

func (self *ModelView) Render(context *Context, writer *utils.XMLWriter) (err error) {
	var children Views

	iter := self.GetModelIterator(context)
	for model := iter.Next(); model != nil; model = iter.Next() {
		view, err := self.GetModelView(model, context)
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

	return children.Render(context, writer)
}
