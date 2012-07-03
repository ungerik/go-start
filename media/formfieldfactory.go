package media

import (
	"github.com/ungerik/go-start/view"
	"github.com/ungerik/go-start/model"
)

func NewFormFieldFactory(wrapped view.FormFieldFactory) *FormFieldFactory {
	return &FormFieldFactory{view.FormFieldFactoryWrapper{wrapped}}
}

type FormFieldFactory struct {
	view.FormFieldFactoryWrapper
}

func (self *FormFieldFactory) CanCreateInput(metaData *model.MetaData, form *view.Form) bool {
	if _, ok := metaData.Value.Interface().(*ImageRef); ok {
		return true
	}
	return self.FormFieldFactoryWrapper.CanCreateInput(metaData, form)
}

func (self *FormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) view.View {
	if imageRef, ok := metaData.Value.Interface().(*ImageRef); ok {
		imageRef.Image()
		return nil
	}
	return self.FormFieldFactoryWrapper.NewInput(withLabel, metaData, form)
}
