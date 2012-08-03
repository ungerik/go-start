package view

///////////////////////////////////////////////////////////////////////////////
// ShortTag

// ShortTag represents an arbitrary HTML element. It has a smaller footprint than Tag.
type ShortTag struct {
	ViewBase
	Tag     string
	Class   string
	Attribs map[string]string
	Content View
}

func (self *ShortTag) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *ShortTag) Render(response *Response) (err error) {
	response.XML.OpenTag(self.Tag).AttribIfNotDefault("class", self.Class)
	for key, value := range self.Attribs {
		response.XML.Attrib(key, value)
	}
	if self.Content != nil {
		err = self.Content.Render(response)
	}
	response.XML.ForceCloseTag()
	return err
}
