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

func (self *TextPreview) Render(context *Context, writer *utils.XMLWriter) (err error) {
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
			writer.Attrib("href", self.MoreLink.URL(context))
			writer.AttribIfNotDefault("title", self.MoreLink.LinkTitle(context))
			content := self.MoreLink.LinkContent(context)
			if content != nil {
				err = content.Render(context, writer)
			}
			writer.ForceCloseTag() // a
		}
	}
	return err
}
