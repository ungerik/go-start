package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// DynamicView

type DynamicView func(context *Context) (view View, err error)

func (self DynamicView) Init(thisView View) {
}

func (self DynamicView) OnRemove() {
}

func (self DynamicView) ID() string {
	return ""
}

func (self DynamicView) IterateChildren(callback IterateChildrenCallback) {
}

func (self DynamicView) Render(context *Context, writer *utils.XMLWriter) error {
	child, err := self(context)
	if err != nil || child == nil {
		return err
	}
	child.Init(child)
	return child.Render(context, writer)
}
