package view

import (
	"github.com/ungerik/go-start/model"
)

/*
FormLayout is responsible for creating and structuring all dynamic content
of the form including the submit button.
It uses Form.GetFieldFactory() to create the field views.
*/
type FormLayout interface {
	GetDefaultInputSize(metaData *model.MetaData) int

	BeginFormContent(form *Form, ctx *Context, formContent *Views) error
	// SubmitSuccess will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// not returned an error.
	SubmitSuccess(message string, form *Form, ctx *Context, formContent *Views) error
	// SubmitError will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// returned an error.
	SubmitError(message string, form *Form, ctx *Context, formContent *Views) error
	EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, ctx *Context, formContent *Views) error

	BeginNamedFields(namedFields *model.MetaData, form *Form, ctx *Context, formContent *Views) error
	NamedField(field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error
	EndNamedFields(namedFields *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error

	BeginIndexedFields(indexedFields *model.MetaData, form *Form, ctx *Context, formContent *Views) error
	IndexedField(field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error
	EndIndexedFields(indexedFields *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error

	NewHiddenInput(metaData *model.MetaData, form *Form) (View, error)
	NewTableHeader(metaData *model.MetaData, form *Form) (View, error)
	NewFieldDescrtiption(description string, form *Form) View
	NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View
	NewGeneralErrorMessage(message string, form *Form) View
	NewSuccessMessage(message string, form *Form) View
	NewSubmitButton(text, confirmationMessage string, form *Form) *SubmitButton
	NewAddButton(onclick string, form *Form) View
	NewRemoveButton(onclick string, form *Form) View
	NewUpButton(disabled bool, onclick string, form *Form) View
	NewDownButton(disabled bool, onclick string, form *Form) View
}
