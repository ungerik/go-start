package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Button

// Button represents a HTML input element of type button or submit.
type Button struct {
	ViewBaseWithId
	Name           string
	Value          interface{}
	Submit         bool
	Class          string
	Disabled       bool
	TabIndex       int
	OnClick        string
	OnClickConfirm string // Will add a confirmation dialog for onclick
	Content        View   // Only used when Submit is false
}

func (self *Button) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Button) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	if self.Submit {
		writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		writer.Attrib("type", "submit")
		writer.AttribIfNotDefault("name", self.Name)
		writer.AttribIfNotDefault("value", self.Value)
		if self.Disabled {
			writer.Attrib("disabled", "disabled")
		}
		writer.AttribIfNotDefault("tabindex", self.TabIndex)
		if self.OnClickConfirm != "" {
			writer.Attrib("onclick", "return confirm('", self.OnClickConfirm, "');")
		} else {
			writer.AttribIfNotDefault("onclick", self.OnClick)
		}
		writer.CloseTag()
	} else {
		writer.OpenTag("button").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		writer.Attrib("type", "button")
		writer.AttribIfNotDefault("name", self.Name)
		writer.AttribIfNotDefault("value", self.Value)
		if self.Disabled {
			writer.Attrib("disabled", "disabled")
		}
		writer.AttribIfNotDefault("tabindex", self.TabIndex)
		if self.OnClickConfirm != "" {
			writer.Attrib("onclick", "return confirm('", self.OnClickConfirm, "');")
		} else {
			writer.AttribIfNotDefault("onclick", self.OnClick)
		}
		if self.Content != nil {
			err = self.Content.Render(response)
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
