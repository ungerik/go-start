package view

type Button struct {
	ViewBaseWithId
	Name           string
	Class          string
	Disabled       bool
	TabIndex       int
	OnClick        string
	OnClickConfirm string // Will add a confirmation dialog for onclick
	Content        View   // Only used when Submit is false
}

func (self *Button) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Button) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("button")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("type", "button")
	ctx.Response.XML.AttribIfNotDefault("name", self.Name)
	ctx.Response.XML.AttribFlag("disabled", self.Disabled)
	ctx.Response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.OnClickConfirm != "" {
		ctx.Response.XML.Attrib("onclick", "return confirm('", self.OnClickConfirm, "');")
	} else {
		ctx.Response.XML.AttribIfNotDefault("onclick", self.OnClick)
	}
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.CloseTagAlways()
	return err
}
