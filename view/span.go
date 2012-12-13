package view

///////////////////////////////////////////////////////////////////////////////
// Span

// Span represents a HTML span element.
type Span struct {
	ViewBaseWithId
	Class   string
	Content View
}

func (self *Span) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Span) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("span")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.CloseTagAlways()
	return err
}
