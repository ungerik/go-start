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

func (self *Checkbox) Render(response *Response) (err error) {
	response.XML.OpenTag("input")
	response.XML.Attrib(self.ID(), self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("type", "checkbox")
	response.XML.Attrib("name", self.Name)
	response.XML.Attrib("value", "true")
	if self.Disabled {
		response.XML.Attrib("disabled", "disabled")
	}
	if self.Checked {
		response.XML.Attrib("checked", "checked")
	}
	response.XML.CloseTag()

	if self.Label != "" {
		response.XML.OpenTag("label").Attrib("for", self.ID()).Content(self.Label).CloseTag()
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
