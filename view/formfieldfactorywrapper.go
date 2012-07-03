package view

import (
	"github.com/ungerik/go-start/model"
)

type FormFieldFactoryWrapper struct {
	Wrapped FormFieldFactory
}

func (self *FormFieldFactoryWrapper) CanCreateInput(metaData *model.MetaData, form *Form) bool {
	return self.Wrapped.CanCreateInput(metaData, form)
}

func (self *FormFieldFactoryWrapper) NewInput(withLabel bool, metaData *model.MetaData, form *Form) View {
	return self.Wrapped.NewInput(withLabel, metaData, form)
}

func (self *FormFieldFactoryWrapper) NewHiddenInput(metaData *model.MetaData, form *Form) View {
	return self.Wrapped.NewHiddenInput(metaData, form)
}

func (self *FormFieldFactoryWrapper) NewTableHeader(metaData *model.MetaData, form *Form) View {
	return self.Wrapped.NewTableHeader(metaData, form)
}

func (self *FormFieldFactoryWrapper) NewFieldDescrtiption(description string, form *Form) View {
	return self.Wrapped.NewFieldDescrtiption(description, form)
}

func (self *FormFieldFactoryWrapper) NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View {
	return self.Wrapped.NewFieldErrorMessage(message, metaData, form)
}

func (self *FormFieldFactoryWrapper) NewGeneralErrorMessage(message string, form *Form) View {
	return self.Wrapped.NewGeneralErrorMessage(message, form)
}

func (self *FormFieldFactoryWrapper) NewSuccessMessage(message string, form *Form) View {
	return self.Wrapped.NewSuccessMessage(message, form)
}

func (self *FormFieldFactoryWrapper) NewSubmitButton(text, confirmationMessage string, form *Form) View {
	return self.Wrapped.NewSubmitButton(text, confirmationMessage, form)
}

func (self *FormFieldFactoryWrapper) NewAddButton(onclick string, form *Form) View {
	return self.Wrapped.NewAddButton(onclick, form)
}

func (self *FormFieldFactoryWrapper) NewRemoveButton(onclick string, form *Form) View {
	return self.Wrapped.NewRemoveButton(onclick, form)
}

func (self *FormFieldFactoryWrapper) NewUpButton(disabled bool, onclick string, form *Form) View {
	return self.Wrapped.NewUpButton(disabled, onclick, form)
}

func (self *FormFieldFactoryWrapper) NewDownButton(disabled bool, onclick string, form *Form) View {
	return self.Wrapped.NewDownButton(disabled, onclick, form)
}
