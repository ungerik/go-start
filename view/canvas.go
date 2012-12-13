package view

type Canvas struct {
	ViewBaseWithId
	Class  string
	Width  int
	Height int
}

func (self *Canvas) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("label")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("width", self.Width).Attrib("height", self.Height)
	ctx.Response.XML.CloseTagAlways()
	return err
}
