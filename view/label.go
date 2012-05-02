package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Label

type Label struct {
	ViewBaseWithId
	Class   string
	For     View
	Content View
}

func (self *Label) Render(request *Request, session *Session, response *Response) (err error) {
	var forID string
	if self.For != nil {
		forID = self.For.ID()
	}
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("label").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.AttribIfNotDefault("for", forID)
	if self.Content != nil {
		err = self.Content.Render(request, session, response)
	}
	writer.CloseTag()
	return err
}
