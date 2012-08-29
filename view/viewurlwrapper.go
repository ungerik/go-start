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

func (self *ViewURLWrapper) ID() string {
	return self.View.ID()
}

func (self *ViewURLWrapper) IterateChildren(callback IterateChildrenCallback) {
	self.View.IterateChildren(callback)
}

func (self *ViewURLWrapper) Render(ctx *Context) (err error) {
	return self.View.Render(ctx)
}

func (self *ViewURLWrapper) SetPath(path string) {
	self.path = path
}

func (self *ViewURLWrapper) URL(ctx *Context) string {
	return StringURL(self.path).URL(ctx)
}
