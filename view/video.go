package view

import (
	"github.com/ungerik/go-start/errs"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// Video

// Video shows a Youtube Video, other formats to come.
type Video struct {
	ViewBaseWithId
	Class  string
	URL    string
	Width  int
	Height int
	//Description string
}

func (self *Video) Render(response *Response) (err error) {
	youtubeId := ""

	switch {
	case strings.HasPrefix(self.URL, "http://youtu.be/"):
		i := len("http://youtu.be/")
		youtubeId = self.URL[i : i+11]

	case strings.HasPrefix(self.URL, "http://www.youtube.com/watch?v="):
		i := len("http://www.youtube.com/watch?v=")
		youtubeId = self.URL[i : i+11]
	}

	if youtubeId != "" {
		response.XML.OpenTag("iframe").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		width := self.Width
		if width == 0 {
			width = 640
		}
		height := self.Height
		if height == 0 {
			height = 390
		}
		response.XML.Attrib("src", "http://www.youtube.com/embed/", youtubeId)
		response.XML.Attrib("width", width)
		response.XML.Attrib("height", height)
		response.XML.Attrib("frameborder", "0")
		response.XML.Attrib("allowfullscreen", "allowfullscreen")
		response.XML.CloseTag()
		return nil
	}

	return errs.Format("Unsupported video URL: %s", self.URL)
}
