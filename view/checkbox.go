package view

///////////////////////////////////////////////////////////////////////////////
// Checkbox

// Checkbox represents a HTML input element of type checkbox.
type Checkbox struct {
	ViewBaseWithId
	Name     string
	Label    string
	Checked  bool
	Disabled bool
	Class    string
}

func (self *Checkbox) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("input")
	ctx.Response.XML.Attrib(self.ID(), self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("type", "checkbox")
	ctx.Response.XML.Attrib("name", self.Name)
	ctx.Response.XML.Attrib("value", "true")
	if self.Disabled {
		ctx.Response.XML.Attrib("disabled", "disabled")
	}
	if self.Checked {
		ctx.Response.XML.Attrib("checked", "checked")
	}
	ctx.Response.XML.CloseTag()

	if self.Label != "" {
		ctx.Response.XML.OpenTag("label").Attrib("for", self.ID()).Content(self.Label).CloseTag()
	}
	return nil
}

//func (self *Checkbox) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *Checkbox) SetLabel(label string) {
//	self.Label = label
//	ViewChanged(self)
//}
