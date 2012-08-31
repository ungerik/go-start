package media

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func newImageVersion(filename, contentType string, sourceRect image.Rectangle, width, height int, grayscale bool) ImageVersion {
	version := ImageVersion{
		Filename:    model.String(filename),
		ContentType: model.String(contentType),
		Width:       model.Int(width),
		Height:      model.Int(height),
		Grayscale:   model.Bool(grayscale),
	}
	version.SourceRect.SetRectangle(sourceRect)
	return version
}

type ImageVersion struct {
	image        *Image
	ID           model.String `bson:",omitempty"`
	Filename     model.String
	ContentType  model.String
	SourceRect   ModelRect
	OutsideColor model.Color
	Width        model.Int
	Height       model.Int
	Grayscale    model.Bool
}

func (self *ImageVersion) URL() view.URL {
	return view.NewURLWithArgs(ImageView, self.ID.Get(), self.Filename.Get())
}

// AspectRatio returns Width / Height
func (self *ImageVersion) AspectRatio() float64 {
	return float64(self.Width) / float64(self.Height)
}

func (self *ImageVersion) SaveImageData(data []byte) error {
	writer, err := Config.Backend.ImageVersionWriter(self)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (self *ImageVersion) SaveImage(im image.Image) error {
	writer, err := Config.Backend.ImageVersionWriter(self)
	if err != nil {
		return err
	}
	switch self.ContentType {
	case "image/jpeg":
		err = jpeg.Encode(writer, im, nil)
	case "image/png":
		err = png.Encode(writer, im)
	default:
		return errors.New("Can't save content-type: " + self.ContentType.Get())
	}
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (self *ImageVersion) LoadImage() (image.Image, error) {
	reader, _, err := Config.Backend.ImageVersionReader(self.ID.Get())
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(reader)
	if err != nil {
		reader.Close()
		return nil, err
	}
	err = reader.Close()
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (self *ImageVersion) ViewImage(class string) *view.Image {
	return &view.Image{
		URL:         self.URL(),
		Width:       self.Width.GetInt(),
		Height:      self.Height.GetInt(),
		Description: self.image.Description.Get(),
		Class:       class,
	}
}

func (self *ImageVersion) ViewImageLink(imageClass, linkClass string) *view.Link {
	return &view.Link{
		Model: &view.StringLink{
			Url:     self.image.Link.Get(),
			Title:   self.image.Description.Get(),
			Content: self.ViewImage(imageClass),
		},
		Class: linkClass,
	}
}
