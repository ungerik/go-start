package view

///////////////////////////////////////////////////////////////////////////////
// HiddenInput

type HiddenInput struct {
	ViewBaseWithId
	Name  string
	Value string
}

func (self *HiddenInput) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("input")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.Attrib("type", "hidden")
	ctx.Response.XML.Attrib("name", self.Name)
	ctx.Response.XML.Attrib("value", self.Value)
	ctx.Response.XML.CloseTag()
	return nil
}

//func (self *HiddenInput) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *HiddenInput) SetValue(value string) {
//	self.Value = value
//	ViewChanged(self)
//}
