package media

import (
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
