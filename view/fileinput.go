package view

type FileInput struct {
	ViewBaseWithId
	Class    string
	Name     string
	Disabled bool
}

func (self *FileInput) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("input")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("type", "file").Attrib("name", self.Name)
	ctx.Response.XML.AttribFlag("disabled", self.Disabled)
	ctx.Response.XML.CloseTag()
	return err
}
