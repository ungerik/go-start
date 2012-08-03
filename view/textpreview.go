package view

import (
	"unicode"
)

///////////////////////////////////////////////////////////////////////////////
// TextPreview

type TextPreview struct {
	ViewBase
	PlainText   string
	MaxLength   int
	ShortLength int // Shortened length if len(Text) > MaxLength. If zero, MaxLength will be used
	MoreLink    LinkModel
}

func (self *TextPreview) Render(response *Response) (err error) {
	if len(self.PlainText) < self.MaxLength {
		response.XML.Content(self.PlainText)
	} else {
		shortLength := self.ShortLength
		if shortLength == 0 {
			shortLength = self.MaxLength
		}

		// If in the middle of a word, go back to space before it
		for shortLength > 0 && !unicode.IsSpace(rune(self.PlainText[shortLength-1])) {
			shortLength--
		}

		// If in the middle of space, go back to word before it
		for shortLength > 0 && unicode.IsSpace(rune(self.PlainText[shortLength-1])) {
			shortLength--
		}

		response.XML.Content(self.PlainText[:shortLength])
		response.XML.Content("... ")
		if self.MoreLink != nil {
			response.XML.OpenTag("a")
			response.XML.Attrib("href", self.MoreLink.URL(response.Request.URLArgs...))
			response.XML.AttribIfNotDefault("title", self.MoreLink.LinkTitle())
			content := self.MoreLink.LinkContent()
			if content != nil {
				err = content.Render(response)
			}
			response.XML.ForceCloseTag() // a
		}
	}
	return err
}
