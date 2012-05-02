package view

func NewViewURLWrapper(view View) *ViewURLWrapper {
	return &ViewURLWrapper{View: view}
}

type ViewURLWrapper struct {
	View View
	path string
}

func (self *ViewURLWrapper) Init(thisView View) {
	self.View.Init(self.View)
}

func (self *ViewURLWrapper) OnRemove() {
	self.View.OnRemove()
}

func (self *ViewURLWrapper) ID() string {
	return self.View.ID()
}

func (self *ViewURLWrapper) IterateChildren(callback IterateChildrenCallback) {
	self.View.IterateChildren(callback)
}

func (self *ViewURLWrapper) Render(request *Request, session *Session, response *Response) (err error) {
	return self.View.Render(request, session, response)
}

func (self *ViewURLWrapper) SetPath(path string) {
	self.path = path
}

func (self *ViewURLWrapper) URL(request *Request, session *Session, response *Response, args ...string) string {
	path := StringURL(self.path).URL(request, session, response, args...)
	return "http://" + request.Host + path
}
