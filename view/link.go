package view

import

//	"bytes"
"github.com/ungerik/go-start/utils"

//func RenderLink(response *Response, link *Link) (html string, err error) {
//	var buf bytes.Buffer
//	err = link.Render(context, utils.NewXMLWriter(&buf))
//	return buf.String(), err
//}
//
//func RenderLinkModel(response *Response, model LinkModel) (html string, err error) {
//	return RenderLink(context, &Link{Model: model})
//}

///////////////////////////////////////////////////////////////////////////////
// Link

type Link struct {
	ViewBaseWithId
	Class     string
	Model     LinkModel
	NewWindow bool
}

func (self *Link) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("a").AttribIfNotDefault("id", self.id)
	writer.AttribIfNotDefault("class", self.Class)
	if self.NewWindow {
		writer.Attrib("target", "_blank")
	}
	if self.Model != nil {
		writer.Attrib("href", self.Model.URL(response))
		writer.AttribIfNotDefault("title", self.Model.LinkTitle(response))
		writer.AttribIfNotDefault("rel", self.Model.LinkRel(response))
		content := self.Model.LinkContent(response)
		if content != nil {
			err = content.Render(response)
		}
	}
	writer.ExtraCloseTag() // a
	return err
}
