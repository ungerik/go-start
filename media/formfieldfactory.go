package media

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func NewFormFieldFactory(wrapped view.FormFieldFactory, thumbnailsize int) *FormFieldFactory {
	return &FormFieldFactory{
		FormFieldFactoryWrapper: view.FormFieldFactoryWrapper{wrapped},
		ThumbnailSize:           thumbnailsize,
	}
}

type FormFieldFactory struct {
	view.FormFieldFactoryWrapper
	ThumbnailSize int
}

func (self *FormFieldFactory) CanCreateInput(metaData *model.MetaData, form *view.Form) bool {
	if _, ok := metaData.Value.Interface().(*ImageRef); ok {
		return true
	}
	return self.FormFieldFactoryWrapper.CanCreateInput(metaData, form)
}

func (self *FormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (view.View, error) {
	if imageRef, ok := metaData.Value.Interface().(*ImageRef); ok {
		imageRef.Image()
		return nil, nil
	}
	return self.FormFieldFactoryWrapper.NewInput(withLabel, metaData, form)
}
