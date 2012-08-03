package view

import ()

///////////////////////////////////////////////////////////////////////////////
// Link

type Link struct {
	ViewBaseWithId
	Class     string
	Model     LinkModel
	NewWindow bool
}

func (self *Link) Render(response *Response) (err error) {
	response.XML.OpenTag("a").AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	if self.NewWindow {
		response.XML.Attrib("target", "_blank")
	}
	if self.Model != nil {
		response.XML.Attrib("href", self.Model.URL(response.Request.URLArgs...))
		response.XML.AttribIfNotDefault("title", self.Model.LinkTitle())
		response.XML.AttribIfNotDefault("rel", self.Model.LinkRel())
		content := self.Model.LinkContent()
		if content != nil {
			err = content.Render(response)
		}
	}
	response.XML.ForceCloseTag() // a
	return err
}
