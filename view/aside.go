package view

// Aside represents a HTML aside element.
type Aside struct {
	ViewBaseWithId
	Class   string
	Style   string
	Content View
}

func (self *Aside) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Aside) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("aside")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.AttribIfNotDefault("style", self.Style)
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.CloseTagAlways()
	return err
}
