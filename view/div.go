package view

///////////////////////////////////////////////////////////////////////////////
// Div

// Div represents a HTML div element.
type Div struct {
	ViewBaseWithId
	Class   string
	Style   string
	Content View
}

func (self *Div) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Div) Render(response *Response) (err error) {
	response.XML.OpenTag("div").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	response.XML.AttribIfNotDefault("style", self.Style)
	if self.Content != nil {
		err = self.Content.Render(response)
	}
	response.XML.ForceCloseTag()
	return err
}
