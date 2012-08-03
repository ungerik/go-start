package view

import (
	"fmt"
	"net/url"
)

///////////////////////////////////////////////////////////////////////////////
// DummyImage

func NewDummyImage(width, height int) *DummyImage {
	return &DummyImage{Width: width, Height: height}
}

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

func (self *DummyImage) Render(response *Response) (err error) {
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

	response.XML.OpenTag("img").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("src", src)
	response.XML.AttribIfNotDefault("width", self.Width)
	response.XML.AttribIfNotDefault("height", self.Height)
	response.XML.AttribIfNotDefault("alt", self.Text)
	response.XML.CloseTag()
	return nil
}

//func (self *DummyImage) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
