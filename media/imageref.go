package media

import (
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
	return metaData.BoolAttrib("required")
}

func (self *ImageRef) Validate(metaData *model.MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return model.NewRequiredError(metaData)
	}
	return nil
}

func (self *ImageRef) getImage(width, height int, grayscale bool) (*Image, *view.Image, error) {
	image, err := Config.Backend.LoadImage(self.String())
	if err != nil {
		return nil, nil, err
	}
	url, err := image.VersionURL(width, height, grayscale)
	if err != nil {
		return nil, nil, err
	}
	viewImage := &view.Image{
		URL:         url,
		Width:       width,
		Height:      height,
		Description: image.Description.Get(),
	}
	return image, viewImage, err
}

func (self *ImageRef) Image(width, height int) (*view.Image, error) {
	_, viewImage, err := self.getImage(width, height, false)
	return viewImage, err
}

func (self *ImageRef) LinkedImage(width, height int) (*view.Link, error) {
	image, viewImage, err := self.getImage(width, height, false)
	if err != nil {
		return nil, err
	}
	return &view.Link{
		Model: &view.StringLink{
			Url:     image.Link.Get(),
			Title:   image.Description.Get(),
			Content: viewImage,
		},
	}, nil
}

func (self *ImageRef) GrayscaleImage(width, height int) (*view.Image, error) {
	_, viewImage, err := self.getImage(width, height, true)
	return viewImage, err
}

func (self *ImageRef) LinkedGrayscaleImage(width, height int) (*view.Link, error) {
	image, viewImage, err := self.getImage(width, height, true)
	if err != nil {
		return nil, err
	}
	return &view.Link{
		Model: &view.StringLink{
			Url:     image.Link.Get(),
			Title:   image.Description.Get(),
			Content: viewImage,
		},
	}, nil
}

func (self *ImageRef) ImageView(width, height int, grayscale bool, imageClass, linkClass, defaultURL string) (view.View, error) {
	if self.IsEmpty() {
		return &view.Image{
			URL:    defaultURL,
			Width:  width,
			Height: height,
		}, nil
	}
	image, viewImage, err := self.getImage(width, height, false)
	if err != nil {
		return nil, err
	}
	viewImage.Class = imageClass
	if image.Link.IsEmpty() {
		return viewImage, nil
	}
	return &view.Link{
		Class: linkClass,
		Model: &view.StringLink{
			Url:     image.Link.Get(),
			Title:   image.Description.Get(),
			Content: viewImage,
		},
	}, nil
}
