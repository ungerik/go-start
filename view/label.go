package view

///////////////////////////////////////////////////////////////////////////////
// Label

type Label struct {
	ViewBaseWithId
	Class   string
	For     View
	Content View
}

func (self *Label) Render(ctx *Context) (err error) {
	var forID string
	if self.For != nil {
		forID = self.For.ID()
	}
	ctx.Response.XML.OpenTag("label")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.AttribIfNotDefault("for", forID)
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.CloseTag()
	return err
}
