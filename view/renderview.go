package view

/*
RenderView implements all View methods for a View.Render compatible function.

Example:

	renderView := RenderView(
		func(response *Response) error {
			writer.Write([]byte("<html><body>Any Content</body></html>"))
			return nil
		},
	)
*/
type RenderView func(response *Response) error

func (self RenderView) Init(thisView View) {
}

func (self RenderView) ID() string {
	return ""
}

func (self RenderView) IterateChildren(callback IterateChildrenCallback) {
}

func (self RenderView) Render(response *Response) error {
	return self(response)
}
