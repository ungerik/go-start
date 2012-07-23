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
	if _, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
		return true
	}
	return self.FormFieldFactoryWrapper.Wrapped.CanCreateInput(metaData, form)
}

func (self *FormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (view.View, error) {
	if imageRef, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
		thumbnailFrame := &view.Div{
			Class: "thumbnail-frame",
			Style: fmt.Sprintf("width:%dpx;height:%dpx;", self.ThumbnailSize, self.ThumbnailSize),
		}
		if !imageRef.IsEmpty() {
			image, err := imageRef.Image()
			if err != nil {
				return nil, err
			}
			version, err := image.Thumbnail(self.ThumbnailSize)
			if err != nil {
				return nil, err
			}
			thumbnailFrame.Content = version.ViewImage("")
		}
		result := view.DIV(self.ImageWidgetClass,
			thumbnailFrame,
			&view.Button{Content: view.HTML("Change")},
			&view.Button{Content: view.HTML("Remove")},
			// description
			// link
			view.DivClearBoth(),
		)
		if withLabel {
			return view.AddStandardLabel(form, result, metaData), nil
		}
		return result, nil
	}
	return self.FormFieldFactoryWrapper.Wrapped.NewInput(withLabel, metaData, form)
}
