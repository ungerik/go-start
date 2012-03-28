package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Renderer

// not used atm.
type Renderer func(context *Context, writer *utils.XMLWriter) (err error)

func (self Renderer) Init(thisView View) {
}

func (self Renderer) OnRemove() {
}

func (self Renderer) ID() string {
	return ""
}

func (self Renderer) IterateChildren(callback IterateChildrenCallback) {
}

func (self Renderer) Render(context *Context, writer *utils.XMLWriter) (err error) {
	return self(context, writer)
}
