package media

import (
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

func (self *FormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (view.View, error) {
	if imageRef, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
		editor, err := ImageRefEditor(imageRef, metaData.Selector(), self.ImageWidgetClass, self.ThumbnailSize)
		if err != nil {
			return nil, err
		}
		if withLabel {
			return view.AddStandardLabel(form, editor, metaData), nil
		}
		return editor, nil
	}
	return self.FormFieldFactoryWrapper.Wrapped.NewInput(withLabel, metaData, form)
}
