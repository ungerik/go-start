package media

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/view"
)

type ImageVersion struct {
	Width     model.Int
	Height    model.Int
	Grayscale model.Bool
	BackendID model.String
}

type Image struct {
	Filename    model.String
	Description model.String
	LinkURL     model.Url
	Original    ImageVersion
	Versions    []ImageVersion
}

func (self *Image) VersionURL(width, height int, grayscale bool) (string, error) {
	if width == self.Original.Width.GetInt() && height == self.Original.Height.GetInt() && self.Original.Grayscale.Get() == grayscale {
		return "", nil
	}
	return "", nil
}
