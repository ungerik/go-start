package view

///////////////////////////////////////////////////////////////////////////////
// Label

type Label struct {
	ViewBaseWithId
	Class   string
	For     View
	Content View
}

func (self *Label) Render(response *Response) (err error) {
	var forID string
	if self.For != nil {
		forID = self.For.ID()
	}
	response.XML.OpenTag("label")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	response.XML.AttribIfNotDefault("for", forID)
	if self.Content != nil {
		err = self.Content.Render(response)
	}
	response.XML.CloseTag()
	return err
}
