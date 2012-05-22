package view

import (
	// "fmt"
	"github.com/ungerik/go-start/model"
	// "reflect"
)

/*
VerticalFormLayout.

CSS needed for VerticalFormLayout:

	form label:after {
		content: ":";
	}

	form input[type=checkbox] + label:after {
		content: "";
	}

Additional CSS for labels above input fields (except checkboxes):

	form label {
		display: block;
	}

	form input[type=checkbox] + label {
		display: inline;
	}

DIV classes for coloring:

	form .required {}
	form .error {}
	form .success {}

*/
type VerticalFormLayout struct {
	form *Form
}

func (self *VerticalFormLayout) BeginFormContent(form *Form, formFields Views) Views {
	return formFields
}

func (self *VerticalFormLayout) EndFormContent(form *Form, formFields Views) Views {
	formId := &HiddenInput{Name: "form_id", Value: form.FormID}
	submitButton := form.GetFieldFactory().NewSubmitButton(form.GetSubmitButtonText(), form)
	return append(formFields, formId, submitButton)
}

func (self *VerticalFormLayout) BeginStruct(strct *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}

func (self *VerticalFormLayout) StructField(field *model.MetaData, form *Form, formFields Views) Views {
	if field.Kind == model.ValueKind {
		views := make(Views, 0, 2)
		fieldFactory := form.GetFieldFactory()
		input := fieldFactory.NewInput(field, form)
		label := fieldFactory.NewLabel(input, field, form)
		views = append(views, label)
		// for _, error := range errors {
		// 	views = append(views)
		// }
		views = append(views, input)
		formFields = append(formFields, P(views))
	}
	return formFields
}

func (self *VerticalFormLayout) EndStruct(strct *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}

func (self *VerticalFormLayout) BeginArray(array *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}

func (self *VerticalFormLayout) ArrayField(field *model.MetaData, form *Form, formFields Views) Views {

	return self.StructField(field, form, formFields) // todo replace

	return formFields
}

func (self *VerticalFormLayout) EndArray(array *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}

func (self *VerticalFormLayout) BeginSlice(slice *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}

func (self *VerticalFormLayout) SliceField(field *model.MetaData, form *Form, formFields Views) Views {

	return self.StructField(field, form, formFields) // todo replace

	return formFields
}

func (self *VerticalFormLayout) EndSlice(slice *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}
