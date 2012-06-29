package media

import (
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/view"
	"bytes"
	"image"
	"image/png"
	_ "image/gif"
	_ "code.google.com/p/go.image/tiff"
	_ "code.google.com/p/go.image/bmp"
	"image/color"
)

// NewImage creates a new Image and saves the original version to Config.Backend.
// GIF, TIFF, BMP images will be read, but written as PNG.
func NewImage(file model.File) (*Image, error) {
	i, t, err := image.Decode(bytes.NewReader(file.Data))
	if err != nil {
		return nil, err
	}
	if t == "gif" || t == "bmp" || t == "tiff" {
		var buf bytes.Buffer
		err = png.Encode(&buf, i)
		if err != nil {
			return nil, err
		}
		file.Data = buf.Bytes()
		file.Name += ".png"
		i, t, err = image.Decode(bytes.NewReader(file.Data))
		if err != nil {
			return nil, err
		}
	}
	result := &Image{
		Versions: []ImageVersion{{
			Filename:    model.String(ValidUrlFilename(file.Name)),
			ContentType: model.String("image/" + t),
			Width:       model.Int(i.Bounds().Dx()),
			Height:      model.Int(i.Bounds().Dy()),
			Grayscale:   model.Bool(i.ColorModel() == color.GrayModel || i.ColorModel() == color.Gray16Model),
		}},
	}
	err = result.Versions[0].SaveData(file.Data)
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

func (self *Image) outerSizeWithOriginalAspectRatio(width, height int) (int, int) {
	originalAspectRatio := float32(self.Width()) / float32(self.Height())
	aspectRatio := float32(width) / float32(height)
	if aspectRatio > originalAspectRatio {
		// Wider than original
		return width, int(float32(width) / originalAspectRatio)
	}
	// Heigher than original
	return int(float32(height) * originalAspectRatio), height
}

func (self *Image) newVersion(width, height int, grayscale bool) (*ImageVersion, error) {
	image, err := self.Versions[0].LoadImage()
	if err != nil {
		return nil, err
	}
	if image == nil {
	}
	version := &ImageVersion{
		Filename:    self.Versions[0].Filename,
		ContentType: self.Versions[0].ContentType,
		Width:       model.Int(width),
		Height:      model.Int(height),
		Grayscale:   model.Bool(grayscale),
	}
	return version, nil
}

func (self *Image) Version(width, height int, grayscale bool) (*ImageVersion, error) {
	if self.Grayscale() {
		// Ignore color requests when original image is grayscale
		grayscale = true
	}

	// If requested image is larger than original size, return original
	if width > self.Width() || height > self.Height() {

	}

	// Search for exact match
	for i := range self.Versions {
		version := &self.Versions[i]
		if width == version.Width.GetInt() && height == version.Height.GetInt() && version.Grayscale.Get() == grayscale {
			return version, nil
		}
	}
	// 

	return nil, nil
}
