package media

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	_ "code.google.com/p/go.image/bmp"
	_ "code.google.com/p/go.image/tiff"
	"code.google.com/p/graphics-go/graphics"

	"github.com/ungerik/go-start/debug"
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

func ImageIterator() model.Iterator {
	return Config.Backend.ImageIterator()
}

// NewImage creates and saves a new Image.
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
	image.Init()
	err = version.SaveImageData(data)
	if err != nil {
		return nil, err
	}
	image.Size.SetInt(len(data))
	err = image.Save()
	if err != nil {
		return nil, err
	}
	return image, nil
}

// NewImageFromURL creates and saves a new Image by downloading it from url.
// GIF, TIFF, BMP images will be read, but saved as PNG.
func NewImageFromURL(url string) (*Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return NewImageFromReader(path.Base(url), response.Body)
}

// NewImageFromReader creates and saves a new Image from reader.
// GIF, TIFF, BMP images will be read, but saved as PNG.
func NewImageFromReader(filename string, reader io.Reader) (*Image, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return NewImage(filename, data)
}

// LoadImage loads an existing image from Config.Backend.
func LoadImage(id string) (*Image, error) {
	return Config.Backend.LoadImage(id)
}

type Image struct {
	ID       model.String `bson:",omitempty"`
	Title    model.String
	Link     model.Url
	Size     model.Int
	Versions []ImageVersion
}

func (self *Image) Init() *Image {
	if self.Title == "" {
		self.Title.Set(self.Filename())
	}
	for i := range self.Versions {
		self.Versions[i].image = self
	}
	return self
}

func (self *Image) Save() error {
	if len(self.Versions) == 0 {
		return fmt.Errorf("Trying to save an Image without an original version, soemthing went wrong\n%v", self)
	}
	return Config.Backend.SaveImage(self)
}

func (self *Image) Delete() error {
	return Config.Backend.DeleteImage(self)
}

func (self *Image) CountRefs() (int, error) {
	return Config.Backend.CountImageRefs(self.ID.Get())
}

func (self *Image) RemoveAllRefs() (count int, err error) {
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
	if index <= 0 || index >= len(self.Versions) {
		return fmt.Errorf("Invalid index %d", index)
	}
	err := Config.Backend.DeleteFile(self.Versions[index].ID.Get())
	if err != nil {
		return err
	}
	self.Versions = append(self.Versions[:index], self.Versions[index+1:]...)
	return nil
}

func (self *Image) DeleteVersions() error {

	for len(self.Versions) > 1 {
		err := self.DeleteVersion(1)
		if err != nil {
			return err
		}

	}
	return nil
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

func (self *Image) Rectangle() image.Rectangle {
	return self.Versions[0].SourceRect.Rectangle()
}

func (self *Image) Grayscale() bool {
	return self.Versions[0].Grayscale.Get()
}

func (self *Image) FileURL() view.URL {
	return self.Versions[0].FileURL()
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

// VersionSourceRect searches and returns an existing matching version,
// or a new one will be created and saved.
func (self *Image) VersionSourceRect(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	debug.Nop()
	// debug.Printf(
	// 	"VersionSourceRect: from %dx%d image take rectangle [%d,%d,%d,%d] (%dx%d) and scale it to %dx%d",
	// 	self.Width(),
	// 	self.Height(),
	// 	sourceRect.Min.X,
	// 	sourceRect.Min.Y,
	// 	sourceRect.Max.X,
	// 	sourceRect.Max.Y,
	// 	sourceRect.Dx(),
	// 	sourceRect.Dy(),
	// 	width,
	// 	height,
	// )

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
	if grayscale {
		versionImage = image.NewGray(image.Rect(0, 0, width, height))
	} else {
		versionImage = image.NewRGBA(image.Rect(0, 0, width, height))
	}
	if sourceRect.In(self.Rectangle()) {
		// debug.Print("VersionSourceRect: rectangle is within image")

		// versionImage = ResampleImage(origImage, sourceRect, width, height)
		subImage := SubImageWithoutOffset(origImage, sourceRect)
		err = graphics.Scale(versionImage.(draw.Image), subImage)
		if err != nil {
			return nil, err
		}
		if grayscale && !self.Grayscale() {
			var grayVersion image.Image = image.NewGray(versionImage.Bounds())
			draw.Draw(grayVersion.(draw.Image), versionImage.Bounds(), versionImage, image.ZP, draw.Src)
			versionImage = grayVersion
		}
	} else {
		// debug.Print("VersionSourceRect: rectangle is not completely within image, using outsideColor")

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

		// destImage := ResampleImage(origImage, origImage.Bounds(), destRect.Dx(), destRect.Dy())
		// draw.Draw(versionImage.(draw.Image), destRect, destImage, image.ZP, draw.Src)
		subImage := SubImageWithoutOffset(origImage, sourceRect)
		destImage := SubImageWithoutOffset(versionImage, destRect)
		err = graphics.Scale(destImage.(draw.Image), subImage)
		if err != nil {
			return nil, err
		}
	}

	// Save new image version
	version := self.addVersion(self.Filename(), self.ContentType(), sourceRect, width, height, grayscale)
	err = version.SaveImage(versionImage)
	if err != nil {
		return nil, err
	}
	err = self.Save()
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

func (self *Image) VersionSourceRectView(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color, class string) (*view.Image, error) {
	version, err := self.VersionSourceRect(sourceRect, width, height, grayscale, outsideColor)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, class string) (*view.Image, error) {
	version, err := self.Version(width, height, horAlign, verAlign, grayscale)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionCenteredView(width, height int, grayscale bool, class string) (*view.Image, error) {
	version, err := self.VersionCentered(width, height, grayscale)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionWidthView(width int, grayscale bool, class string) (*view.Image, error) {
	version, err := self.VersionWidth(width, grayscale)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionHeightView(height int, grayscale bool, class string) (*view.Image, error) {
	version, err := self.VersionHeight(height, grayscale)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionTouchOrigFromOutsideView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, outsideColor color.Color, class string) (*view.Image, error) {
	version, err := self.VersionTouchOrigFromOutside(width, height, horAlign, verAlign, grayscale, outsideColor)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionTouchOrigFromOutsideCenteredView(width, height int, grayscale bool, outsideColor color.Color, class string) (*view.Image, error) {
	version, err := self.VersionTouchOrigFromOutsideCentered(width, height, grayscale, outsideColor)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) ThumbnailView(size int, class string) (*view.Image, error) {
	version, err := self.Thumbnail(size)
	if err != nil {
		return nil, err
	}
	return version.View(class), nil
}

func (self *Image) VersionSourceRectLinkedView(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.VersionSourceRect(sourceRect, width, height, grayscale, outsideColor)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) VersionLinkedView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.Version(width, height, horAlign, verAlign, grayscale)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) VersionCenteredLinkedView(width, height int, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.VersionCentered(width, height, grayscale)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) VersionWidthLinkedView(width int, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.VersionWidth(width, grayscale)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) VersionHeightLinkedView(height int, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.VersionHeight(height, grayscale)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) VersionTouchOrigFromOutsideLinkedView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, outsideColor color.Color, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.VersionTouchOrigFromOutside(width, height, horAlign, verAlign, grayscale, outsideColor)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) VersionTouchOrigFromOutsideCenteredLinkedView(width, height int, grayscale bool, outsideColor color.Color, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.VersionTouchOrigFromOutsideCentered(width, height, grayscale, outsideColor)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}

func (self *Image) ThumbnailLinkedView(size int, imageClass, linkClass string) (*view.Link, error) {
	version, err := self.Thumbnail(size)
	if err != nil {
		return nil, err
	}
	return version.LinkedView(imageClass, linkClass), nil
}
