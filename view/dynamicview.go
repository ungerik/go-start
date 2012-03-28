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
	GetView GetViewFunc // nil Views will be ignored
	child   View
}

func (self *DynamicView) IterateChildren(callback IterateChildrenCallback) {
	if self.child != nil {
		callback(self, self.child)
	}
}

func (self *DynamicView) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.child != nil {
		self.child.OnRemove()
		self.child = nil
	}

	self.child, err = self.GetView(context)
	if err != nil || self.child == nil {
		return err
	}
	self.child.Init(self.child)

	return self.child.Render(context, writer)
}
