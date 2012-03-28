package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// HTML

type HTML string

func (self HTML) Init(thisView View) {
}

func (self HTML) OnRemove() {
}

func (self HTML) ID() string {
	return ""
}

func (self HTML) IterateChildren(callback IterateChildrenCallback) {
}

func (self HTML) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.Content(string(self))
	return nil
}
