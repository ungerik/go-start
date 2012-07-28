package media

import (
	"errors"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type ImageRef string

func (self *ImageRef) String() string {
	return string(*self)
}

func (self *ImageRef) SetString(str string) error {
	*self = ImageRef(str)
	return nil
}

func (self *ImageRef) IsEmpty() bool {
	return *self == ""
}

func (self *ImageRef) Required(metaData *model.MetaData) bool {
	return metaData.BoolAttrib(view.StructTagKey, "required")
}

func (self *ImageRef) Validate(metaData *model.MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return model.NewRequiredError(metaData)
	}
	return nil
}

func (self *ImageRef) Image() (*Image, error) {
	if self.IsEmpty() {
		return nil, errors.New("Can't get Image from empty ImageRef")
	}
	return Config.Backend.LoadImage(self.String())
}
