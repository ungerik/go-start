package view

type TextFieldType int

const (
	NormalTextField TextFieldType = iota
	PasswordTextField
	EmailTextField
)

///////////////////////////////////////////////////////////////////////////////
// TextField

type TextField struct {
	ViewBaseWithId
	Text        string
	Name        string
	Size        int
	MaxLength   int
	Type        TextFieldType
	Readonly    bool
	Disabled    bool
	TabIndex    int
	Class       string
	Placeholder string
}

func (self *TextField) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("input")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)

	ctx.Response.XML.Attrib("name", self.Name)
	ctx.Response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.Readonly {
		ctx.Response.XML.Attrib("readonly", "readonly")
	}
	if self.Disabled {
		ctx.Response.XML.Attrib("disabled", "disabled")
	}

	switch self.Type {
	case PasswordTextField:
		ctx.Response.XML.Attrib("type", "password")
	case EmailTextField:
		ctx.Response.XML.Attrib("type", "email")
	default:
		ctx.Response.XML.Attrib("type", "text")
	}

	ctx.Response.XML.AttribIfNotDefault("size", self.Size)
	ctx.Response.XML.AttribIfNotDefault("maxlength", self.MaxLength)
	ctx.Response.XML.AttribIfNotDefault("placeholder", self.Placeholder)

	ctx.Response.XML.Attrib("value", self.Text)

	ctx.Response.XML.CloseTag()
	return nil
}

//func (self *TextField) SetText(text string) {
//	self.Text = text
//	ViewChanged(self)
//}
//
//func (self *TextField) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *TextField) SetReadonly(readonly bool) {
//	self.Readonly = readonly
//	ViewChanged(self)
//}
//
//func (self *TextField) SetMaxLength(maxlength int) {
//	self.MaxLength = maxlength
//	ViewChanged(self)
//}
//
//func (self *TextField) SetTabIndex(tabindex int) {
//	self.TabIndex = tabindex
//	ViewChanged(self)
//}
//
//func (self *TextField) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
