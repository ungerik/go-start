package view

import "github.com/ungerik/go-start/utils"

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
	Text      string
	Name      string
	Size      int
	MaxLength int
	Type      TextFieldType
	Readonly  bool
	Disabled  bool
	TabIndex  int
	Class     string
	Placeholder string
}

func (self *TextField) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	writer.Attrib("name", self.Name)
	writer.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.Readonly {
		writer.Attrib("readonly", "readonly")
	}
	if self.Disabled {
		writer.Attrib("disabled", "disabled")
	}

	switch self.Type {
	case PasswordTextField:
		writer.Attrib("type", "password")
	case EmailTextField:
		writer.Attrib("type", "email")
	default:
		writer.Attrib("type", "text")
	}

	writer.AttribIfNotDefault("size", self.Size)
	writer.AttribIfNotDefault("maxlength", self.MaxLength)
	writer.Attrib("value", self.Text)
	if self.Placeholder != "" {
		writer.Attrib("placeholder", self.Placeholder)	
	}

	writer.CloseTag()
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
