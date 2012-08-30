package view

///////////////////////////////////////////////////////////////////////////////
// Div

// Div represents a HTML div element.
type Div struct {
	ViewBaseWithId
	Class   string
	Style   string
	Content View
}

func (self *Div) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Div) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("div")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.AttribIfNotDefault("style", self.Style)
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.ForceCloseTag()
	return err
}
