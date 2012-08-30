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

func (self *indirectViewWithURL) Render(ctx *Context) (err error) {
	return (*self.viewWithURL).Render(ctx)
}

func (self *indirectViewWithURL) URL(ctx *Context) string {
	return (*self.viewWithURL).URL(ctx)
}

func (self *indirectViewWithURL) SetPath(path string) {
	(*self.viewWithURL).SetPath(path)
}
