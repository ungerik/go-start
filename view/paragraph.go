package view

///////////////////////////////////////////////////////////////////////////////
// Paragraph

type Paragraph struct {
	ViewBaseWithId
	Class   string
	Content View
}

func (self *Paragraph) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Paragraph) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("p")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.CloseTagAlways()
	return err
}
