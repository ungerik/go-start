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

func (self *TextPreview) Render(ctx *Context) (err error) {
	if len(self.PlainText) < self.MaxLength {
		ctx.Response.XML.Content(self.PlainText)
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

		ctx.Response.XML.Content(self.PlainText[:shortLength])
		ctx.Response.XML.Content("... ")
		if self.MoreLink != nil {
			ctx.Response.XML.OpenTag("a")
			ctx.Response.XML.Attrib("href", self.MoreLink.URL(ctx))
			ctx.Response.XML.AttribIfNotDefault("title", self.MoreLink.LinkTitle(ctx))
			content := self.MoreLink.LinkContent(ctx)
			if content != nil {
				err = content.Render(ctx)
			}
			ctx.Response.XML.CloseTagAlways() // a
		}
	}
	return err
}
