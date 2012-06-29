package media

import (
	"image"
	"github.com/ungerik/go-start/model"
)

type ImageVersion struct {
	ID          model.String `bson:",omitempty"`
	Filename    model.String
	ContentType model.String
	Width       model.Int
	Height      model.Int
	Grayscale   model.Bool
}

func (self *ImageVersion) URL() string {
	return View.URL(self.ID.Get(), self.Filename.Get())
}

func (self *ImageVersion) SaveData(data []byte) error {
	return nil
}

func (self *ImageVersion) LoadImage() (image.Image, error) {
	return nil, nil
}
