package view

///////////////////////////////////////////////////////////////////////////////
// NotFoundView

type NotFoundView struct {
	ViewBase
	Message string
}

func (self *NotFoundView) Render(request *Request, session *Session, response *Response) (err error) {
	return NotFound(self.Message)
}
