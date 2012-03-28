package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Button

// Button represents a HTML input element of type button or submit.
type Button struct {
	ViewBaseWithId
	Name     string
	Value    interface{}
	Submit   bool
	Class    string
	Disabled bool
	TabIndex int
}

func (self *Button) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	if self.Submit {
		writer.Attrib("type", "submit")
	} else {
		writer.Attrib("type", "button")
	}
	writer.Attrib("name", self.Name)
	writer.Attrib("value", self.Value)
	if self.Disabled {
		writer.Attrib("disabled", "disabled")
	}
	writer.AttribIfNotDefault("tabindex", self.TabIndex)
	writer.CloseTag()
	return nil
}

//func (self *Button) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *Button) SetValue(value string) {
//	self.Value = value
//	ViewChanged(self)
//}
