package media

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"

	_ "code.google.com/p/go.image/bmp"
	_ "code.google.com/p/go.image/tiff"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type HorAlignment int

const (
	HorCenter HorAlignment = iota
	Left
	Right
)

type VerAlignment int

const (
	VerCenter VerAlignment = iota
	Top
	Bottom
)

// LoadImage loads an existing image from Config.Backend.
func LoadImage(id string) (*Image, error) {
	return Config.Backend.LoadImage(id)
}

// NewImageFromURL creates a new Image by downloading it from url.
// GIF, TIFF, BMP images will be read, but saved as PNG.
func NewImageFromURL(url string) (*Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	filename := ""
	defer response.Body.Close()
	return NewImageFromReader(filename, response.Body)
}

// NewImageFromReader creates a new Image and saves the original version to Config.Backend.
// GIF, TIFF, BMP images will be read, but saved as PNG.
func NewImageFromReader(filename string, reader io.Reader) (*Image, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return NewImage(filename, data)
}

// NewImage creates a new Image and saves the original version to Config.Backend.
// GIF, TIFF, BMP images will be read, but saved as PNG.
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

	image := new(Image)

	version := image.addVersion(
		MakePrettyUrlFilename(filename),
		"image/"+t,
		i.Bounds(),
		i.Bounds().Dx(),
		i.Bounds().Dy(),
		i.ColorModel() == color.GrayModel || i.ColorModel() == color.Gray16Model,
	)
	err = version.SaveImageData(data)
	if err != nil {
		return nil, err
	}

	return image, nil
}

type Image struct {
	ID       model.String `bson:",omitempty"`
	Title    model.String
	Link     model.Url
	Versions []ImageVersion
}

func (self *Image) Init() {
	for i := range self.Versions {
		self.Versions[i].image = self
	}
}

func (self *Image) Save() error {
	return Config.Backend.SaveImage(self)
}

func (self *Image) Delete() error {
	return Config.Backend.DeleteImage(self)
}

func (self *Image) CountRefs() (int, error) {
	return Config.Backend.CountImageRefs(self.ID.Get())
}

func (self *Image) RemoveAllRefs() error {
	return Config.Backend.RemoveAllImageRefs(self.ID.Get())
}

func (self *Image) addVersion(filename, contentType string, sourceRect image.Rectangle, width, height int, grayscale bool) *ImageVersion {
	version := ImageVersion{
		image:       self,
		Filename:    model.String(filename),
		ContentType: model.String(contentType),
		Width:       model.Int(width),
		Height:      model.Int(height),
		Grayscale:   model.Bool(grayscale),
	}
	version.SourceRect.SetRectangle(sourceRect)
	self.Versions = append(self.Versions, version)
	return &self.Versions[len(self.Versions)-1]
}

func (self *Image) DeleteVersion(index int) error {
	err := Config.Backend.DeleteImageVersion(self.Versions[index].ID.Get())
	if err != nil {
		return err
	}
	self.Versions = append(self.Versions[:index], self.Versions[index+1:]...)
	return nil
}

func (self *Image) Filename() string {
	return self.Versions[0].Filename.Get()
}

// TitleOrFilename returns Title if not empty,
// or else Filename().
func (self *Image) TitleOrFilename() string {
	if self.Title.IsEmpty() {
		return self.Filename()
	}
	return self.Title.Get()
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

func (self *Image) Rectangle() image.Rectangle {
	return self.Versions[0].SourceRect.Rectangle()
}

func (self *Image) Grayscale() bool {
	return self.Versions[0].Grayscale.Get()
}

func (self *Image) GetURL() view.URL {
	return self.Versions[0].GetURL()
}

// AspectRatio returns Width / Height
func (self *Image) AspectRatio() float64 {
	return self.Versions[0].AspectRatio()
}

func (self *Image) sourceRectTouchOriginalFromOutside(width, height int, horAlign HorAlignment, verAlign VerAlignment) (r image.Rectangle) {
	var offset image.Point
	aspectRatio := float64(width) / float64(height)
	if aspectRatio > self.AspectRatio() {
		// Wider than original
		// so touchOriginalFromOutside means
		// that the source rect is as high as the original
		r.Max.X = int(float64(self.Height()) * aspectRatio)
		r.Max.Y = self.Height()
		switch horAlign {
		case HorCenter:
			offset.X = (self.Width() - r.Max.X) / 2
		case Right:
			offset.X = self.Width() - r.Max.X
		}
	} else {
		// Heigher than original,
		// so touchOriginalFromOutside means
		// that the source rect is as wide as the original
		r.Max.X = self.Width()
		r.Max.Y = int(float64(self.Width()) / aspectRatio)
		switch verAlign {
		case VerCenter:
			offset.Y = (self.Height() - r.Max.Y) / 2
		case Bottom:
			offset.Y = self.Height() - r.Max.Y
		}
	}
	return r.Add(offset)
}

func (self *Image) sourceRectTouchOriginalFromInside(width, height int, horAlign HorAlignment, verAlign VerAlignment) (r image.Rectangle) {
	var offset image.Point
	aspectRatio := float64(width) / float64(height)
	if aspectRatio > self.AspectRatio() {
		// Wider than original
		// so touchOriginalFromInside means
		// that the source rect is as wide as the original
		r.Max.X = self.Width()
		r.Max.Y = int(float64(self.Width()) / aspectRatio)
		switch verAlign {
		case VerCenter:
			offset.Y = (self.Height() - r.Max.Y) / 2
		case Bottom:
			offset.Y = self.Height() - r.Max.Y
		}
	} else {
		// Heigher than original,
		// so touchOriginalFromInside means
		// that the source rect is as high as the original
		r.Max.X = int(float64(self.Height()) * aspectRatio)
		r.Max.Y = self.Height()
		switch horAlign {
		case HorCenter:
			offset.X = (self.Width() - r.Max.X) / 2
		case Right:
			offset.X = self.Width() - r.Max.X
		}
	}
	return r.Add(offset)
}

func (self *Image) OriginalVersion() *ImageVersion {
	return &self.Versions[0]
}

// SourceRectVersion searches and returns an existing matching version,
// or a new one will be created and saved.
func (self *Image) VersionSourceRect(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	if self.Grayscale() {
		grayscale = true // Ignore color requests when original image is grayscale
	}

	// Search for exact match
	for i := range self.Versions {
		v := &self.Versions[i]
		match := v.SourceRect.Rectangle() == sourceRect &&
			v.Width.GetInt() == width &&
			v.Height.GetInt() == height &&
			v.OutsideColor.EqualsColor(outsideColor) &&
			v.Grayscale.Get() == grayscale
		if match {
			return v, nil
		}
	}

	// No exact match, create version
	origImage, err := self.Versions[0].LoadImage()
	if err != nil {
		return nil, err
	}

	var versionImage image.Image
	if sourceRect.In(self.Rectangle()) {
		versionImage = ResampleImage(origImage, sourceRect, width, height)
		if grayscale && !self.Grayscale() {
			var grayVersion image.Image = image.NewGray(versionImage.Bounds())
			draw.Draw(grayVersion.(draw.Image), versionImage.Bounds(), versionImage, image.ZP, draw.Src)
			versionImage = grayVersion
		}
	} else {
		if grayscale {
			versionImage = image.NewGray(image.Rect(0, 0, width, height))
		} else {
			versionImage = image.NewRGBA(image.Rect(0, 0, width, height))
		}
		// Fill version with outsideColor
		draw.Draw(versionImage.(draw.Image), versionImage.Bounds(), image.NewUniform(outsideColor), image.ZP, draw.Src)
		// Where to draw the source image into the version image
		var destRect image.Rectangle
		if !(sourceRect.Min.X < 0 || sourceRect.Min.Y < 0) {
			panic("touching from outside means that sourceRect x or y must be negative")
		}
		sourceW := float64(sourceRect.Dx())
		sourceH := float64(sourceRect.Dy())
		destRect.Min.X = int(float64(-sourceRect.Min.X) / sourceW * float64(width))
		destRect.Min.Y = int(float64(-sourceRect.Min.Y) / sourceH * float64(height))
		destRect.Max.X = destRect.Min.X + int(float64(self.Width())/sourceW*float64(width))
		destRect.Max.Y = destRect.Min.Y + int(float64(self.Height())/sourceH*float64(height))
		destImage := ResampleImage(origImage, origImage.Bounds(), destRect.Dx(), destRect.Dy())
		draw.Draw(versionImage.(draw.Image), destRect, destImage, image.ZP, draw.Src)
	}

	// Save new image version
	version := self.addVersion(self.Filename(), self.ContentType(), sourceRect, width, height, grayscale)
	err = version.SaveImage(versionImage)
	if err != nil {
		return nil, err
	}
	err = Config.Backend.SaveImage(self)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (self *Image) Version(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool) (im *ImageVersion, err error) {
	return self.VersionSourceRect(self.sourceRectTouchOriginalFromInside(width, height, horAlign, verAlign), width, height, grayscale, color.RGBA{})
}

func (self *Image) VersionCentered(width, height int, grayscale bool) (im *ImageVersion, err error) {
	return self.Version(width, height, HorCenter, VerCenter, grayscale)
}

func (self *Image) VersionWidth(width int, grayscale bool) (im *ImageVersion, err error) {
	height := int(float64(width)/self.AspectRatio() + 0.5)
	return self.Version(width, height, HorCenter, VerCenter, grayscale)
}

func (self *Image) VersionHeight(height int, grayscale bool) (im *ImageVersion, err error) {
	width := int(float64(height)*self.AspectRatio() + 0.5)
	return self.Version(width, height, HorCenter, VerCenter, grayscale)
}

func (self *Image) VersionTouchOrigFromOutside(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	return self.VersionSourceRect(self.sourceRectTouchOriginalFromOutside(width, height, horAlign, verAlign), width, height, grayscale, outsideColor)
}

func (self *Image) VersionTouchOrigFromOutsideCentered(width, height int, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	return self.VersionTouchOrigFromOutside(width, height, HorCenter, VerCenter, grayscale, outsideColor)
}

func (self *Image) Thumbnail(size int) (im *ImageVersion, err error) {
	return self.VersionCentered(size, size, self.Grayscale())
}
