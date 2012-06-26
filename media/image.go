package media

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/view"
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
)

func NewImage(file model.File) (*Image, error) {
	i, t, err := image.Decode(bytes.NewReader(file.Data))
	if err != nil {
		return nil, err
	}
	return &Image{
		Versions: []ImageVersion{{
			Filename:    model.String(ValidUrlFilename(file.Name)),
			ContentType: model.String("image/" + t),
			Width:       model.Int(i.Bounds().Dx()),
			Height:      model.Int(i.Bounds().Dy()),
			Grayscale:   model.Bool(i.ColorModel() == color.GrayModel || i.ColorModel() == color.Gray16Model),
		}},
	}, nil
}

type Image struct {
	ID          model.String `bson:",omitempty"`
	Description model.String
	Link        model.Url
	Versions    []ImageVersion
}

func (self *Image) Filename() string {
	return self.Versions[0].Filename.Get()
}

func (self *Image) ContentType() string {
	return self.Versions[0].ContentType.Get()
}

func (self *Image) VersionURL(width, height int, grayscale bool) (string, error) {
	for i := range self.Versions {
		version := &self.Versions[i]
		if width == version.Width.GetInt() && height == version.Height.GetInt() && version.Grayscale.Get() == grayscale {
			return version.URL(), nil
		}
	}
	return "", nil
}
