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
	Onclick  string
	Content  View // Only used when Submit is false
}

func (self *Button) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Button) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.Submit {
		writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		writer.Attrib("type", "submit")
		writer.Attrib("name", self.Name)
		writer.Attrib("value", self.Value)
		if self.Disabled {
			writer.Attrib("disabled", "disabled")
		}
		writer.AttribIfNotDefault("tabindex", self.TabIndex)
		writer.AttribIfNotDefault("onclick", self.Onclick)
		writer.CloseTag()
	} else {
		writer.OpenTag("button").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		writer.Attrib("type", "button")
		writer.Attrib("name", self.Name)
		writer.Attrib("value", self.Value)
		if self.Disabled {
			writer.Attrib("disabled", "disabled")
		}
		writer.AttribIfNotDefault("tabindex", self.TabIndex)
		writer.AttribIfNotDefault("onclick", self.Onclick)
		if self.Content != nil {
			err = self.Content.Render(context, writer)
		}
		writer.ForceCloseTag()
	}
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
