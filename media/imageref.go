package media

import (
	"image"
	"image/color"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func NewImageRef(image *Image) *ImageRef {
	imageRef := new(ImageRef)
	imageRef.Set(image)
	return imageRef
}

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

func (self *ImageRef) TryGet() (*Image, bool, error) {
	if self.IsEmpty() {
		return nil, false, nil
	}
	image, err := Config.Backend.LoadImage(self.String())
	if _, notFound := err.(ErrNotFound); notFound {
		return nil, false, nil
	}
	return image, err == nil, err
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

func (self *ImageRef) VersionWidth(width int, grayscale bool) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionWidth(width, grayscale)
}

func (self *ImageRef) VersionHeight(height int, grayscale bool) (im *ImageVersion, err error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionHeight(height, grayscale)
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

///

func (self *ImageRef) OriginalVersionView(class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.OriginalVersion().View(class), nil
}

func (self *ImageRef) VersionSourceRectView(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionSourceRectView(sourceRect, width, height, grayscale, outsideColor, class)
}

func (self *ImageRef) VersionView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionView(width, height, horAlign, verAlign, grayscale, class)
}

func (self *ImageRef) VersionCenteredView(width, height int, grayscale bool, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionCenteredView(width, height, grayscale, class)
}

func (self *ImageRef) VersionWidthView(width int, grayscale bool, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionWidthView(width, grayscale, class)
}

func (self *ImageRef) VersionHeightView(height int, grayscale bool, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionHeightView(height, grayscale, class)
}

func (self *ImageRef) VersionTouchOrigFromOutsideView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, outsideColor color.Color, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionTouchOrigFromOutsideView(width, height, horAlign, verAlign, grayscale, outsideColor, class)
}

func (self *ImageRef) VersionTouchOrigFromOutsideCenteredView(width, height int, grayscale bool, outsideColor color.Color, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionTouchOrigFromOutsideCenteredView(width, height, grayscale, outsideColor, class)
}

func (self *ImageRef) ThumbnailView(size int, class string) (*view.Image, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.ThumbnailView(size, class)
}

func (self *ImageRef) OriginalVersionLinkedView(imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.OriginalVersion().LinkedView(imageClass, linkClass), nil
}

func (self *ImageRef) VersionSourceRectLinkedView(sourceRect image.Rectangle, width, height int, grayscale bool, outsideColor color.Color, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionSourceRectLinkedView(sourceRect, width, height, grayscale, outsideColor, imageClass, linkClass)
}

func (self *ImageRef) VersionLinkedView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionLinkedView(width, height, horAlign, verAlign, grayscale, imageClass, linkClass)
}

func (self *ImageRef) VersionCenteredLinkedView(width, height int, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionCenteredLinkedView(width, height, grayscale, imageClass, linkClass)
}

func (self *ImageRef) VersionWidthLinkedView(width int, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionWidthLinkedView(width, grayscale, imageClass, linkClass)
}

func (self *ImageRef) VersionHeightLinkedView(height int, grayscale bool, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionHeightLinkedView(height, grayscale, imageClass, linkClass)
}

func (self *ImageRef) VersionTouchOrigFromOutsideLinkedView(width, height int, horAlign HorAlignment, verAlign VerAlignment, grayscale bool, outsideColor color.Color, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionTouchOrigFromOutsideLinkedView(width, height, horAlign, verAlign, grayscale, outsideColor, imageClass, linkClass)
}

func (self *ImageRef) VersionTouchOrigFromOutsideCenteredLinkedView(width, height int, grayscale bool, outsideColor color.Color, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.VersionTouchOrigFromOutsideCenteredLinkedView(width, height, grayscale, outsideColor, imageClass, linkClass)
}

func (self *ImageRef) ThumbnailLinkedView(size int, imageClass, linkClass string) (*view.Link, error) {
	image, err := self.Get()
	if err != nil {
		return nil, err
	}
	return image.ThumbnailLinkedView(size, imageClass, linkClass)
}
