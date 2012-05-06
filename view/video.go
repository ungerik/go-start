package view

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/utils"
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
	case utils.StringStartsWith(self.URL, "http://youtu.be/"):
		i := len("http://youtu.be/")
		youtubeId = self.URL[i : i+11]

	case utils.StringStartsWith(self.URL, "http://www.youtube.com/watch?v="):
		i := len("http://www.youtube.com/watch?v=")
		youtubeId = self.URL[i : i+11]
	}

	if youtubeId != "" {
		writer := utils.NewXMLWriter(response)
		writer.OpenTag("iframe").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		width := self.Width
		if width == 0 {
			width = 640
		}
		height := self.Height
		if height == 0 {
			height = 390
		}
		writer.Attrib("src", "http://www.youtube.com/embed/", youtubeId)
		writer.Attrib("width", width)
		writer.Attrib("height", height)
		writer.Attrib("frameborder", "0")
		writer.Attrib("allowfullscreen", "allowfullscreen")
		writer.CloseTag()
		return nil
	}

	return errs.Format("Unsupported video URL: %s", self.URL)
}
