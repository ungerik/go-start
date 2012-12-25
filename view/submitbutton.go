package view

type SubmitButton struct {
	ViewBaseWithId
	Name           string
	Value          interface{}
	Class          string
	Disabled       bool
	TabIndex       int
	OnClick        string
	OnClickConfirm string // Will add a confirmation dialog for onclick
}

func (self *SubmitButton) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("input")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("type", "submit")
	ctx.Response.XML.AttribIfNotDefault("name", self.Name)
	ctx.Response.XML.AttribIfNotDefault("value", self.Value)
	ctx.Response.XML.AttribFlag("disabled", self.Disabled)
	ctx.Response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.OnClickConfirm != "" {
		ctx.Response.XML.Attrib("onclick", self.OnClick, "; return confirm('", self.OnClickConfirm, "');")
	} else {
		ctx.Response.XML.AttribIfNotDefault("onclick", self.OnClick)
	}
	ctx.Response.XML.CloseTag()
	return nil
}
