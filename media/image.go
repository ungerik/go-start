package media

import (
	"bytes"
	"image"
	"image/png"
	_ "image/gif"
	_ "code.google.com/p/go.image/tiff"
	_ "code.google.com/p/go.image/bmp"
	"image/color"
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/view"
)

// NewImage creates a new Image and saves the original version to Config.Backend.
// GIF, TIFF, BMP images will be read, but written as PNG.
func NewImage(filename string, data []byte) (*Image, error) {
	i, t, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if t == "gif" || t == "bmp" || t == "tiff" {
		var buf bytes.Buffer
		err = png.Encode(&buf, i)
		if err != nil {
			return nil, err
		}
		data = buf.Bytes()
		filename += ".png"
		i, t, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	}
	result := &Image{
		Versions: []ImageVersion{{
			Filename:    model.String(ValidUrlFilename(filename)),
			ContentType: model.String("image/" + t),
			Width:       model.Int(i.Bounds().Dx()),
			Height:      model.Int(i.Bounds().Dy()),
			Grayscale:   model.Bool(i.ColorModel() == color.GrayModel || i.ColorModel() == color.Gray16Model),
		}},
	}
	err = result.Versions[0].SaveImageData(data)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func (self *Image) Width() int {
	return self.Versions[0].Width.GetInt()
}

func (self *Image) Height() int {
	return self.Versions[0].Height.GetInt()
}

func (self *Image) Grayscale() bool {
	return self.Versions[0].Grayscale.Get()
}

// AspectRatio returns Width / Height
func (self *Image) AspectRatio() float64 {
	return self.Versions[0].AspectRatio()
}

func (self *Image) touchFromOutsideWithOriginalAspectRatio(width, height int) (int, int) {
	aspectRatio := float64(width) / float64(height)
	originalAspectRatio := self.AspectRatio()
	if aspectRatio > originalAspectRatio {
		// Wider than original
		return width, int(float64(width) / originalAspectRatio)
	}
	// Heigher than original
	return int(float64(height) * originalAspectRatio), height
}

// func (self *Image) newVersion(width, height int, grayscale bool) (*ImageVersion, error) {
// 	image, err := self.Versions[0].LoadImage()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if image == nil {
// 	}
// 	version := &ImageVersion{
// 		Filename:    self.Versions[0].Filename,
// 		ContentType: self.Versions[0].ContentType,
// 		Width:       model.Int(width),
// 		Height:      model.Int(height),
// 		Grayscale:   model.Bool(grayscale),
// 	}
// 	return version, nil
// }

func (self *Image) Version(width, height int, grayscale bool) (im *ImageVersion, err error) {
	if self.Grayscale() {
		// Ignore color requests when original image is grayscale
		grayscale = true
	}

	// aspectRatio := float64(width) / float64(height)

	// If requested image is larger than original size, return original
	if width > self.Width() || height > self.Height() {
		// todo
	}

	// Search for exact match
	for i := range self.Versions {
		version := &self.Versions[i]
		if width == version.Width.GetInt() && height == version.Height.GetInt() && version.Grayscale.Get() == grayscale {
			return version, nil
		}
	}
	// 

	// outerWidth, outerHeight := self.touchFromOutsideWithOriginalAspectRatio(width, height)
	// orig, err := self.Versions[0].LoadImage()
	if err != nil {
		return nil, err
	}
	// var r image.Rectangle
	// scaled := ResizeImage(orig, r, width, height)

	version := &ImageVersion{
		Filename:    self.Versions[0].Filename,
		ContentType: self.Versions[0].ContentType,
		Width:       model.Int(width),
		Height:      model.Int(height),
		Grayscale:   model.Bool(grayscale),
	}
	self.Versions = append(self.Versions, *version)

	return nil, nil
}
