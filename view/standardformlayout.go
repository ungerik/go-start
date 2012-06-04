package view

import (
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
)

/*
StandardFormLayout.

CSS needed for StandardFormLayout:

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
type StandardFormLayout struct {
}

func (self *StandardFormLayout) BeginFormContent(form *Form, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) SubmitSuccess(message string, form *Form, formFields Views) Views {
	return append(formFields, form.GetFieldFactory().NewSuccessMessage(message, form))
}

func (self *StandardFormLayout) SubmitError(message string, form *Form, formFields Views) Views {
	return append(formFields, form.GetFieldFactory().NewGeneralErrorMessage(message, form))
}

func (self *StandardFormLayout) EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, formFields Views) Views {
	fieldFactory := form.GetFieldFactory()
	for _, err := range generalValidationErrs {
		formFields = append(formFields, fieldFactory.NewGeneralErrorMessage(err.Error(), form))
		formFields = append(Views{fieldFactory.NewGeneralErrorMessage(err.Error(), form)}, formFields...)
	}
	formId := &HiddenInput{Name: FormIDName, Value: form.FormID}
	submitButton := fieldFactory.NewSubmitButton(form.GetSubmitButtonText(), form)
	return append(formFields, formId, submitButton)
}

func (self *StandardFormLayout) BeginStruct(strct *model.MetaData, form *Form, formFields Views) Views {
	if strct.Parent != nil && strct.Parent.Kind == model.ArrayKind {
		formFields = append(formFields, &Table{HeaderRow: true, Model: new(ViewsTableModel)})
	}
	return formFields
}

func (self *StandardFormLayout) StructField(field *model.MetaData, validationErr error, form *Form, formFields Views) Views {
	fieldFactory := form.GetFieldFactory()
	if !fieldFactory.CanCreateInput(field, form) || form.IsFieldExcluded(field) {
		return formFields
	}
	if form.IsFieldHidden(field) {
		return append(formFields, fieldFactory.NewHiddenInput(field, form))
	}
	views := make(Views, 0, 2)
	input := fieldFactory.NewInput(field, form)
	if self.fieldNeedsLabel(field) {
		label := fieldFactory.NewLabel(input, field, form)
		views = append(views, label)
	}
	if validationErr != nil {
		views = append(views, fieldFactory.NewFieldErrorMessage(validationErr.Error(), field, form))
	}
	views = append(views, input)
	return append(formFields, DIV(Config.Form.StandardFormLayoutDivClass, views))
}

func (self *StandardFormLayout) EndStruct(strct *model.MetaData, validationErr error, form *Form, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) BeginArray(array *model.MetaData, form *Form, formFields Views) Views {
	debug.Nop()
	return formFields
}

func (self *StandardFormLayout) ArrayField(field *model.MetaData, validationErr error, form *Form, formFields Views) Views {
	return self.StructField(field, validationErr, form, formFields) // todo replace
	return formFields
}

func (self *StandardFormLayout) EndArray(array *model.MetaData, validationErr error, form *Form, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) BeginSlice(slice *model.MetaData, form *Form, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) SliceField(field *model.MetaData, validationErr error, form *Form, formFields Views) Views {

	return self.StructField(field, validationErr, form, formFields) // todo replace

	return formFields
}

func (self *StandardFormLayout) EndSlice(slice *model.MetaData, validationErr error, form *Form, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) fieldNeedsLabel(field *model.MetaData) bool {
	switch field.Value.Addr().Interface().(type) {
	case *model.Bool:
		return false
	}
	return true
}
