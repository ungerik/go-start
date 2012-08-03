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

func (self *TextField) Render(response *Response) (err error) {
	response.XML.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	response.XML.Attrib("name", self.Name)
	response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.Readonly {
		response.XML.Attrib("readonly", "readonly")
	}
	if self.Disabled {
		response.XML.Attrib("disabled", "disabled")
	}

	switch self.Type {
	case PasswordTextField:
		response.XML.Attrib("type", "password")
	case EmailTextField:
		response.XML.Attrib("type", "email")
	default:
		response.XML.Attrib("type", "text")
	}

	response.XML.AttribIfNotDefault("size", self.Size)
	response.XML.AttribIfNotDefault("maxlength", self.MaxLength)
	response.XML.Attrib("value", self.Text)
	if self.Placeholder != "" {
		response.XML.Attrib("placeholder", self.Placeholder)
	}

	response.XML.CloseTag()
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
