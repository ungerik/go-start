package view

import "github.com/ungerik/go-start/utils"

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

func (self *Checkbox) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("type", "checkbox")
	writer.Attrib("name", self.Name)
	writer.Attrib("value", "true")
	if self.Disabled {
		writer.Attrib("disabled", "disabled")
	}
	if self.Checked {
		writer.Attrib("checked", "checked")
	}
	writer.CloseTag()

	writer.OpenTag("label").Attrib("for", self.id).Content(self.Label).CloseTag()
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
