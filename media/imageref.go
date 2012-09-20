package media

import (
	"image"
	"image/color"

	"github.com/ungerik/go-start/model"
)

type ImageRef string

func (self *ImageRef) String() string {
	return string(*self)
}

func (self *ImageRef) SetString(str string) error {
	*self = ImageRef(str)
	return nil
}

// Image loads the referenced image, or returns nil if the reference is empty.
func (self *ImageRef) Get() (*Image, error) {
	if self.IsEmpty() {
		return nil, nil
	}
	return Config.Backend.LoadImage(self.String())
}

// SetImage sets the ID of image, or an empty reference if image is nil.
func (self *ImageRef) Set(image *Image) {
	if image != nil {
		self.SetString(image.ID.String())
	} else {
		*self = ""
	}
}

func (self *ImageRef) IsEmpty() bool {
	return *self == ""
}

func (self *ImageRef) Required(metaData *model.MetaData) bool {
	return metaData.BoolAttrib(model.StructTagKey, "required")
}

func (self *ImageRef) Validate(metaData *model.MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return model.NewRequiredError(metaData)
	}
	return nil
}

func (self *ImageRef) OriginalVersion() (*ImageVersion, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.OriginalVersion(), nil
}

func (self *ImageRef) VersionSourceRect(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionSourceRect(sourceRect, width, height, grayscale, outsideColor)
}

func (self *ImageRef) Version(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.Version(width, height, horAlign, verAlign, grayscale)
}

func (self *ImageRef) VersionCentered(width, height int, grayscale bool) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionCentered(width, height, grayscale)
}

func (self *ImageRef) VersionTouchOrigFromOutside(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionTouchOrigFromOutside(width, height, horAlign, verAlign, grayscale, outsideColor)
}

func (self *ImageRef) VersionTouchOrigFromOutsideCentered(width, height int, grayscale bool, outsideColor color.Color) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionTouchOrigFromOutsideCentered(width, height, grayscale, outsideColor)
}

func (self *ImageRef) Thumbnail(size int) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.Thumbnail(size)
}
