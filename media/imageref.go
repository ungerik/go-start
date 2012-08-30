package media

import (
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

// SetImage sets the ID of image, or an empty reference if image is nil.
func (self *ImageRef) SetImage(image *Image) {
	if image != nil {
		self.SetString(image.ID.String())
	} else {
		*self = ""
	}
}

// GetImage loads the referenced image, or returns nil if the reference is empty.
func (self *ImageRef) GetImage() (*Image, error) {
	if self.IsEmpty() {
		return nil, nil
	}
	return Config.Backend.LoadImage(self.String())
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
