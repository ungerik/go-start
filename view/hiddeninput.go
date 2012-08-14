package view

///////////////////////////////////////////////////////////////////////////////
// HiddenInput

type HiddenInput struct {
	ViewBaseWithId
	Name  string
	Value string
}

func (self *HiddenInput) Render(response *Response) (err error) {
	response.XML.OpenTag("input")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.Attrib("type", "hidden")
	response.XML.Attrib("name", self.Name)
	response.XML.Attrib("value", self.Value)
	response.XML.CloseTag()
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
