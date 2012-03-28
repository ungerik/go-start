package view

import (
	"fmt"
	"github.com/ungerik/go-start/utils"
	"net/url"
)

///////////////////////////////////////////////////////////////////////////////
// DummyImage

// DummyImage represents a HTML img element with src utilizing http://dummyimage.com.
type DummyImage struct {
	ViewBaseWithId
	Class           string
	Width           int
	Height          int
	BackgroundColor string
	ForegroundColor string
	Text            string
}

func (self *DummyImage) Render(context *Context, writer *utils.XMLWriter) (err error) {
	src := fmt.Sprintf("http://dummyimage.com/%dx%d", self.Width, self.Height)

	if self.BackgroundColor != "" || self.ForegroundColor != "" {
		if self.BackgroundColor != "" {
			src += "/" + self.BackgroundColor
		} else {
			src += "/ccc"
		}

		if self.ForegroundColor != "" {
			src += "/" + self.ForegroundColor
		}
	}

	src += ".png"

	if self.Text != "" {
		src += "&text=" + url.QueryEscape(self.Text)
	}

	writer.OpenTag("img").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("src", src)
	writer.AttribIfNotDefault("width", self.Width)
	writer.AttribIfNotDefault("height", self.Height)
	writer.AttribIfNotDefault("alt", self.Text)
	writer.CloseTag()
	return nil
}

//func (self *DummyImage) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
