package view

///////////////////////////////////////////////////////////////////////////////
// Paragraph

type Paragraph struct {
	ViewBaseWithId
	Class   string
	Content View
}

func (self *Paragraph) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Paragraph) Render(response *Response) (err error) {
	response.XML.OpenTag("p")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	if self.Content != nil {
		err = self.Content.Render(response)
	}
	response.XML.ForceCloseTag()
	return err
}
