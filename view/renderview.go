package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// RenderView

/*
RenderView implements all View methods for a View.Render compatible function.

Example:

	renderView := RenderView(
		func(context *Context, writer *utils.XMLWriter) error {
			writer.Write([]byte("<html><body>Any Content</body></html>"))
			return nil
		},
	)
*/
type RenderView func(context *Context, writer *utils.XMLWriter) error

func (self RenderView) Init(thisView View) {
}

func (self RenderView) OnRemove() {
}

func (self RenderView) ID() string {
	return ""
}

func (self RenderView) IterateChildren(callback IterateChildrenCallback) {
}

func (self RenderView) Render(context *Context, writer *utils.XMLWriter) error {
	return self(context, writer)
}
