package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// HiddenInput

type HiddenInput struct {
	ViewBaseWithId
	Name  string
	Value string
}

func (self *HiddenInput) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("input").Attrib("id", self.id)
	writer.Attrib("type", "hidden")
	writer.Attrib("name", self.Name)
	writer.Attrib("value", self.Value)
	writer.CloseTag()
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
