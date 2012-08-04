package view

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

func (self *indirectViewWithURL) Render(response *Response) (err error) {
	return (*self.viewWithURL).Render(response)
}

func (self *indirectViewWithURL) URL(response *Response) string {
	return (*self.viewWithURL).URL(response)
}

func (self *indirectViewWithURL) SetPath(path string) {
	(*self.viewWithURL).SetPath(path)
}
