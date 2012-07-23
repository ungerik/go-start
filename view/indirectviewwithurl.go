package view

import "github.com/ungerik/go-start/utils"

func IndirectViewWithURL(viewWithURL *ViewWithURL) ViewWithURL {
	return &indirectViewWithURL{viewWithURL}
}

type indirectViewWithURL struct {
	viewWithURL *ViewWithURL
}

func (self *indirectViewWithURL) Init(thisView View) {
	(*self.viewWithURL).Init(thisView)
}

func (self *indirectViewWithURL) ID() string {
	return (*self.viewWithURL).ID()
}

func (self *indirectViewWithURL) IterateChildren(callback IterateChildrenCallback) {
	(*self.viewWithURL).IterateChildren(callback)
}

func (self *indirectViewWithURL) Render(context *Context, writer *utils.XMLWriter) (err error) {
	return (*self.viewWithURL).Render(context, writer)
}

func (self *indirectViewWithURL) URL(args ...string) string {
	return (*self.viewWithURL).URL(args...)
}

func (self *indirectViewWithURL) SetPath(path string) {
	(*self.viewWithURL).SetPath(path)
}
