package media

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

/*
Check out:

https://github.com/valums/file-uploader
http://deepliquid.com/content/Jcrop.html
*/

var Config = Configuration{
	DummyImageColor: "#a8a8a8",
	Admin: AdminConfiguration{
		ImageEditorClass:    "media-image-editor",
		ThumbnailSize:       200,
		ThumbnailFrameClass: "thumbnail-frame",
		ActionsClass:        "actions",
	},
	ImageRefController: ImageRefControllerConfiguration{
		Class:               "media-imageref-editor",
		ThumbnailSize:       100,
		ThumbnailFrameClass: "thumbnail-frame",
		ActionsClass:        "actions",
	},
	BlobRefController: BlobRefControllerConfiguration{
		Class:        "media-blobref-editor",
		ActionsClass: "actions",
	},
}

type AdminConfiguration struct {
	ImageEditorClass    string
	ThumbnailSize       int
	ThumbnailFrameClass string
	ActionsClass        string
	ButtonClass         string
}

type ImageRefControllerConfiguration struct {
	Class               string
	ThumbnailFrameClass string
	ThumbnailSize       int
	ActionsClass        string
}

type BlobRefControllerConfiguration struct {
	Class        string
	ActionsClass string
}

type Configuration struct {
	Backend                 Backend
	NoDynamicStyleAndScript bool
	DummyImageColor         string // web color
	dummyImageURL           string
	Admin                   AdminConfiguration
	ImageRefController      ImageRefControllerConfiguration
	BlobRefController       BlobRefControllerConfiguration
}

func (self *Configuration) Name() string {
	return "media"
}

func (self *Configuration) Init() error {
	c := model.NewColor(self.DummyImageColor)
	self.dummyImageURL = ColoredImageDataURL(c.RGBA())
	view.Config.Form.DefaultFieldControllers = append(view.Config.Form.DefaultFieldControllers,
		ImageRefController{},
		BlobRefController{},
	)
	return nil
}

func (self *Configuration) Close() error {
	return nil
}
