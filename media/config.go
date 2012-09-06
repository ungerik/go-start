package media

import (
	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
)

/*
Check out:

https://github.com/valums/file-uploader
http://deepliquid.com/content/Jcrop.html
*/

var Config = Configuration{
	DummyImageColor: "#a8a8a8",
	ImagesAdmin: ImagesAdminConfiguration{
		Class:          "media-images-admin",
		ThumbnailSize:  100,
		ThumbnailClass: "media-images-admin-thumbnail",
	},
}

type ImagesAdminConfiguration struct {
	Class          string
	ThumbnailSize  int
	ThumbnailClass string
}

type Configuration struct {
	Backend                 Backend
	NoDynamicStyleAndScript bool
	DummyImageColor         string // web color
	dummyImageURL           string
	ImagesAdmin             ImagesAdminConfiguration
}

func (self *Configuration) Name() string {
	return "media"
}

func (self *Configuration) Init() error {
	c := model.NewColor(self.DummyImageColor)
	self.dummyImageURL = ColoredImageDataURL(c.RGBA())
	return nil
}

func (self *Configuration) Close() error {
	return nil
}
