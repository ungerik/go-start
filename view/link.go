package view

import ()

///////////////////////////////////////////////////////////////////////////////
// Link

type Link struct {
	ViewBaseWithId
	Class     string
	Model     LinkModel
	NewWindow bool
}

func (self *Link) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("a")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	if self.NewWindow {
		ctx.Response.XML.Attrib("target", "_blank")
	}
	if self.Model != nil {
		ctx.Response.XML.Attrib("href", self.Model.URL(ctx))
		ctx.Response.XML.AttribIfNotDefault("title", self.Model.LinkTitle(ctx))
		ctx.Response.XML.AttribIfNotDefault("rel", self.Model.LinkRel(ctx))
		content := self.Model.LinkContent(ctx)
		if content != nil {
			err = content.Render(ctx)
		}
	}
	ctx.Response.XML.ForceCloseTag() // a
	return err
}
