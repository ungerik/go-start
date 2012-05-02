package view

/*
RenderView implements all View methods for a View.Render compatible function.

Example:

	renderView := RenderView(
		func(request *Request, session *Session, response *Response, writer *utils.XMLWriter) error {
			writer.Write([]byte("<html><body>Any Content</body></html>"))
			return nil
		},
	)
*/
type RenderView func(request *Request, session *Session, response *Response) error

func (self RenderView) Init(thisView View) {
}

func (self RenderView) OnRemove() {
}

func (self RenderView) ID() string {
	return ""
}

func (self RenderView) IterateChildren(callback IterateChildrenCallback) {
}

func (self RenderView) Render(request *Request, session *Session, response *Response) error {
	return self(request, session, response)
}
