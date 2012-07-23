package view

import (
	"github.com/ungerik/go-start/utils"
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
	writer := utils.NewXMLWriter(response)
	if len(self.PlainText) < self.MaxLength {
		writer.Content(self.PlainText)
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

		writer.Content(self.PlainText[:shortLength])
		writer.Content("... ")
		if self.MoreLink != nil {
			writer.OpenTag("a")
			writer.Attrib("href", self.MoreLink.URL(response))
			writer.AttribIfNotDefault("title", self.MoreLink.LinkTitle(response))
			content := self.MoreLink.LinkContent(response)
			if content != nil {
				err = content.Render(response)
			}
			writer.ForceCloseTag() // a
		}
	}
	return err
}
