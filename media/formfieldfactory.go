package media

import (
	"fmt"

	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func NewFormFieldFactory(wrapped view.FormFieldFactory, imageWidgetClass string, thumbnailsize int) *FormFieldFactory {
	return &FormFieldFactory{
		FormFieldFactoryWrapper: view.FormFieldFactoryWrapper{wrapped},
		ImageWidgetClass:        imageWidgetClass,
		ThumbnailSize:           thumbnailsize,
	}
}

type FormFieldFactory struct {
	view.FormFieldFactoryWrapper
	ImageWidgetClass string
	ThumbnailSize    int
}

func (self *FormFieldFactory) CanCreateInput(metaData *model.MetaData, form *view.Form) bool {
	debug.Nop()
	if metaData.Value.CanAddr() {
		if _, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
			return true
		}
	}
	return self.FormFieldFactoryWrapper.Wrapped.CanCreateInput(metaData, form)
}

const registerUploaderScript = `jQuery(function() {
	var uploader = new qq.FileUploader({
	    element: document.getElementById("%s"),
	    action: "%s",
	    allowedExtensions: ["png", "jpg", "jpeg", "gif", "bmp", "tif", "tiff"],
	    acceptFiles: ["image/png", "image/jpeg", "image/gif", "image/bmp", "image/tiff"],
	    multiple: false
	});
})`

func (self *FormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (view.View, error) {
	if imageRef, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
		var img view.View
		var removeButton = &view.Button{Content: view.HTML("Remove"), OnClick: ""}
		var requires view.View

		if !Config.NoDynamicStyleAndScript {
			// Uses https://github.com/valums/file-uploader
			requires = view.Views{
				view.RequireStyleURL("/media/fileuploader.css", 0),
				view.RequireScriptURL("/media/fileuploader.js", 0),
				view.RequireScriptURL("/media/media.js", 10),
			}
		}

		if imageRef.IsEmpty() {
			removeButton.Disabled = true
			img = view.IMG(Config.dummyImageURL, self.ThumbnailSize, self.ThumbnailSize)
		} else {
			image, err := imageRef.GetImage()
			if err != nil {
				return nil, err
			}
			version, err := image.Thumbnail(self.ThumbnailSize)
			if err != nil {
				return nil, err
			}
			img = version.ViewImage("")
		}

		thumbnailFrame := &view.Div{
			Class:   "media-thumbnail-frame",
			Style:   fmt.Sprintf("width:%dpx;height:%dpx;", self.ThumbnailSize, self.ThumbnailSize),
			Content: img,
		}

		uploader := view.DIV("uploader")
		uploaderID := uploader.ID()
		uploaderScript := view.RenderView(
			func(ctx *view.Context) error {
				uploadURL := UploadImage.URL(ctx.ForURLArgsConvert(self.ThumbnailSize))
				script := fmt.Sprintf(registerUploaderScript, uploaderID, uploadURL)
				ctx.Response.RequireScript(script, 20)
				return nil
			},
		)

		editor := view.DIV(self.ImageWidgetClass,
			requires,
			&view.HiddenInput{Name: metaData.Selector(), Value: imageRef.String()},
			thumbnailFrame,
			view.DIV("media-actions-frame",
				view.HTML("&larr; drag &amp; drop files here"),
				view.BR(),
				removeButton,
				view.BR(),
				uploader,
				uploaderScript,
			),
			view.DivClearBoth(),
		)

		if withLabel {
			return view.AddStandardLabel(form, editor, metaData), nil
		}
		return editor, nil
	}

	return self.FormFieldFactoryWrapper.Wrapped.NewInput(withLabel, metaData, form)
}
