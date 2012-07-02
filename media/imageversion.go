package media

import (
	"image"
	"image/png"
	"image/jpeg"
	"errors"
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
