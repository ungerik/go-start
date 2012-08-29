package view

///////////////////////////////////////////////////////////////////////////////
// NotFoundView

type NotFoundView struct {
	ViewBase
	Message string
}

func (self *NotFoundView) Render(ctx *Context) (err error) {
	return NotFound(self.Message)
}
