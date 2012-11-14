package view

type ViewWithURLBase struct {
	ViewBase
	path string
}

func (self *ViewWithURLBase) SetPath(path string) {
	self.path = path
}

func (self *ViewWithURLBase) URL(ctx *Context) string {
	return StringURL(self.path).URL(ctx)
}
