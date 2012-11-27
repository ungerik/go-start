package view

// Link represents an HTML <a> or <link> element depending on UseLinkTag.
// Content and title of the Model will only be rendered for <a>.
type Link struct {
	ViewBaseWithId
	Class      string
	Model      LinkModel
	NewWindow  bool
	UseLinkTag bool
}

func (self *Link) Render(ctx *Context) (err error) {
	if self.UseLinkTag {
		ctx.Response.XML.OpenTag("link")
	} else {
		ctx.Response.XML.OpenTag("a")
	}
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	if self.NewWindow {
		ctx.Response.XML.Attrib("target", "_blank")
	}
	if self.Model != nil {
		ctx.Response.XML.Attrib("href", self.Model.URL(ctx))
		ctx.Response.XML.AttribIfNotDefault("rel", self.Model.LinkRel(ctx))
	}
	if self.UseLinkTag {
		ctx.Response.XML.CloseTag() // link
	} else {
		ctx.Response.XML.AttribIfNotDefault("title", self.Model.LinkTitle(ctx))
		content := self.Model.LinkContent(ctx)
		if content != nil {
			err = content.Render(ctx)
		}
		ctx.Response.XML.CloseTagAlways() // a
	}
	return err
}

func (self *Link) URL(ctx *Context) string {
	return self.Model.URL(ctx)
}
