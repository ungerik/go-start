package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// NotFoundView

type NotFoundView struct {
	ViewBase
	Message string
}

func (self *NotFoundView) Render(context *Context, writer *utils.XMLWriter) (err error) {
	return NotFound(self.Message)
}
