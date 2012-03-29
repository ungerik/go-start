package view

import "github.com/ungerik/go-start/utils"

type GetViewFunc func(context *Context) (view View, err error)

func NewDynamicView(getView GetViewFunc) *DynamicView {
	return &DynamicView{GetView: getView}
}

///////////////////////////////////////////////////////////////////////////////
// DynamicView

type DynamicView struct {
	ViewBase
	// nil Views will be ignored
	GetView GetViewFunc
}

func (self *DynamicView) Render(context *Context, writer *utils.XMLWriter) error {
	child, err := self.GetView(context)
	if err != nil || child == nil {
		return err
	}
	child.Init(child)
	return child.Render(context, writer)
}
