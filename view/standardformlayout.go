package view

import (
	"fmt"
	"strconv"

	// "github.com/ungerik/go-start/debug"
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
	DefaultInputSize      int
	DefaultTableInputSize int
}

func (self *StandardFormLayout) GetDefaultInputSize(metaData *model.MetaData) int {
	grandParent := metaData.Parent.Parent
	if grandParent != nil && (grandParent.Kind == model.ArrayKind || grandParent.Kind == model.SliceKind) {
		return self.DefaultTableInputSize
	}
	return self.DefaultInputSize
}

func (self *StandardFormLayout) BeginFormContent(form *Form, ctx *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) SubmitSuccess(message string, form *Form, ctx *Context, formContent *Views) error {
	*formContent = append(*formContent, self.NewSuccessMessage(message, form))
	return nil
}

func (self *StandardFormLayout) SubmitError(message string, form *Form, ctx *Context, formContent *Views) error {
	*formContent = append(*formContent, self.NewGeneralErrorMessage(message, form))
	return nil
}

func (self *StandardFormLayout) EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, ctx *Context, formContent *Views) error {
	for _, err := range generalValidationErrs {
		*formContent = append(*formContent, self.NewGeneralErrorMessage(err.Error(), form))
		*formContent = append(Views{self.NewGeneralErrorMessage(err.Error(), form)}, *formContent...)
	}
	if form.GeneralErrorOnFieldError && len(fieldValidationErrs) > 0 {
		e := Config.Form.GeneralErrorMessageOnFieldError
		*formContent = append(*formContent, self.NewGeneralErrorMessage(e, form))
		*formContent = append(Views{self.NewGeneralErrorMessage(e, form)}, *formContent...)
	}
	submitButton := self.NewSubmitButton(form.GetSubmitButtonText(), form.SubmitButtonConfirm, form)
	*formContent = append(*formContent, submitButton)
	return nil
}

func (self *StandardFormLayout) BeginNamedFields(namedFields *model.MetaData, form *Form, ctx *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) NamedField(field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error {
	controllers := form.GetFieldControllers()
	if !controllers.Supports(field, form) {
		return nil
	}

	grandParent := field.Parent.Parent
	if grandParent != nil && grandParent.Kind.HasIndexedFields() {
		return self.structFieldInArrayOrSlice(grandParent, field, validationErr, form, ctx, formContent)
	}

	if form.IsFieldHidden(field) {
		input, err := self.NewHiddenInput(field, form)
		if err != nil {
			return err
		}
		*formContent = append(*formContent, input)
		return nil
	}

	formField, err := controllers.NewInput(true, field, form)
	if err != nil {
		return err
	}
	if validationErr != nil {
		formField = Views{formField, self.NewFieldErrorMessage(validationErr.Error(), field, form)}
	}
	*formContent = append(*formContent, DIV(Config.Form.StandardFormLayoutDivClass, formField))
	return nil
}

func (self *StandardFormLayout) EndNamedFields(namedFields *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) BeginIndexedFields(indexedFields *model.MetaData, form *Form, ctx *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) IndexedField(field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error {
	if field.Kind != model.ValueKind {
		return nil
	}
	arrayOrSlice := field.Parent
	// We expect a Table as last form content field.
	// If it doesn't exist yet because this is the first visible
	// struct field in the first array field, then create it
	var table *Table
	if len(*formContent) > 0 {
		table, _ = (*formContent)[len(*formContent)-1].(*Table)
	}
	if table == nil {
		// First array/slice field, create table and table model.
		header, err := self.NewTableHeader(arrayOrSlice, form)
		if err != nil {
			return err
		}
		table = &Table{
			HeaderRow: true,
			Model:     ViewsTableModel{Views{header}},
		}
		table.Init(table) // get an ID now
		*formContent = append(*formContent, table)
		// Add script for manipulating table rows
		ctx.Response.RequireScriptURL("/js/form.js", 0)
	}
	td, err := form.GetFieldControllers().NewInput(false, field, form)
	if err != nil {
		return err
	}
	if validationErr != nil {
		td = Views{
			td,
			self.NewFieldErrorMessage(validationErr.Error(), field, form),
		}
	}
	table.Model = append(table.Model.(ViewsTableModel), Views{td})
	return nil
}

func (self *StandardFormLayout) EndIndexedFields(indexedFields *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error {
	if len(*formContent) > 0 {
		// Add "Actions" column with buttons to table with slice or array values
		if table, ok := (*formContent)[len(*formContent)-1].(*Table); ok {
			tableModel := table.Model.(ViewsTableModel)
			tableModel[0] = append(tableModel[0], HTML("<div class='table-actions'>Actions</div>"))
			rows := tableModel.Rows()
			for i := 1; i < rows; i++ {
				// todo: script depends on buttons being HTML buttons,
				// but buttons are created by form field factory that doesn't
				// guarantee that. needs decoupling via css classes.
				firstRow := (i == 1)
				lastRow := (i == rows-1)
				buttons := Views{
					self.NewUpButton(firstRow, "gostart_form.moveRowUp(this);", form),
					self.NewDownButton(lastRow, "gostart_form.moveRowDown(this);", form),
				}
				if indexedFields.Kind == model.SliceKind {
					if lastRow {
						buttons = append(buttons, self.NewAddButton("gostart_form.addRow(this);", form))
					} else {
						buttons = append(buttons, self.NewRemoveButton("gostart_form.removeRow(this)", form))
					}
				}
				tableModel[i] = append(tableModel[i], buttons)
			}
			// Add a hidden input to get back the changed length of the table
			// if the user changes it via javascript in the client
			*formContent = append(*formContent,
				&HiddenInput{
					Name:  indexedFields.Selector() + ".length",
					Value: strconv.Itoa(indexedFields.Value.Len()),
				},
			)
		}
	}
	return nil
}

func (self *StandardFormLayout) structFieldInArrayOrSlice(arrayOrSlice, field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error {
	// We expect a Table as last form content field.
	// If it doesn't exist yet because this is the first visible
	// struct field in the first array field, then create it
	var table *Table
	if len(*formContent) > 0 {
		table, _ = (*formContent)[len(*formContent)-1].(*Table)
	}
	if table == nil {
		// First struct field of first array field, create table
		// and table model.
		table = &Table{
			Caption:   form.FieldLabel(arrayOrSlice),
			HeaderRow: true,
			Model:     ViewsTableModel{Views{}},
		}
		table.Init(table) // get an ID now
		*formContent = append(*formContent, table)
		// Add script for manipulating table rows
		ctx.Response.RequireScriptURL("/js/form.js", 0)
	}
	tableModel := table.Model.(ViewsTableModel)
	if field.Parent.Index == 0 {
		// If first array field, add label to table header
		header, err := self.NewTableHeader(field, form)
		if err != nil {
			return err
		}
		tableModel[0] = append(tableModel[0], header)
	}
	if tableModel.Rows()-1 == field.Parent.Index {
		// Create row in table model for this array field
		tableModel = append(tableModel, Views{})
		table.Model = tableModel
	}
	// Append form field in last row for this struct field
	row := &tableModel[tableModel.Rows()-1]

	td, err := form.GetFieldControllers().NewInput(false, field, form)
	if err != nil {
		return err
	}
	if validationErr != nil {
		td = Views{
			td,
			self.NewFieldErrorMessage(validationErr.Error(), field, form),
		}
	}
	*row = append(*row, td)
	return nil
}

func (self *StandardFormLayout) fieldNeedsLabel(field *model.MetaData) bool {
	if field.Value.CanAddr() {
		switch field.Value.Addr().Interface().(type) {
		case *model.Bool:
			return false
		}
	}
	return true
}

func (self *StandardFormLayout) NewHiddenInput(metaData *model.MetaData, form *Form) (View, error) {
	return &HiddenInput{
		Name:  metaData.Selector(),
		Value: fmt.Sprint(metaData.Value.Interface()),
	}, nil
}

func (self *StandardFormLayout) NewTableHeader(metaData *model.MetaData, form *Form) (View, error) {
	switch metaData.Kind {
	case model.ValueKind:
		var label View = Escape(form.DirectFieldLabel(metaData))
		if form.IsFieldRequired(metaData) {
			label = Views{label, form.GetRequiredMarker()}
		}
		return label, nil
	case model.ArrayKind, model.SliceKind:
		return Escape(form.FieldLabel(metaData)), nil
	}
	return nil, fmt.Errorf("Unsupported metaData.Kind for NewTableHeader(): %s", metaData.Kind)
}

func (self *StandardFormLayout) NewFieldDescrtiption(description string, form *Form) View {
	return SPAN(form.GetFieldDescriptionClass(), Escape(description))
}

func (self *StandardFormLayout) NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View {
	return SPAN(form.GetErrorMessageClass(), Escape(message))
}

func (self *StandardFormLayout) NewGeneralErrorMessage(message string, form *Form) View {
	return SPAN(form.GetErrorMessageClass(), Escape(message))
}

func (self *StandardFormLayout) NewSuccessMessage(message string, form *Form) View {
	return SPAN(form.GetSuccessMessageClass(), Escape(message))
}

func (self *StandardFormLayout) NewSubmitButton(text, confirmationMessage string, form *Form) *SubmitButton {
	return &SubmitButton{
		Class:          form.GetSubmitButtonClass(),
		Value:          text,
		OnClickConfirm: confirmationMessage,
	}
}

func (self *StandardFormLayout) NewAddButton(onclick string, form *Form) View {
	return &Button{
		Content: HTML("+"),
		OnClick: onclick,
		Class:   "form-table-button",
	}
}

func (self *StandardFormLayout) NewRemoveButton(onclick string, form *Form) View {
	return &Button{
		Content: HTML("&times;"),
		OnClick: onclick,
		Class:   "form-table-button",
	}
}

func (self *StandardFormLayout) NewUpButton(disabled bool, onclick string, form *Form) View {
	return &Button{
		Content:  HTML("&uarr;"),
		Disabled: disabled,
		OnClick:  onclick,
		Class:    "form-table-button",
	}
}

func (self *StandardFormLayout) NewDownButton(disabled bool, onclick string, form *Form) View {
	return &Button{
		Content:  HTML("&darr;"),
		Disabled: disabled,
		OnClick:  onclick,
		Class:    "form-table-button",
	}
}

func AddStandardLabel(form *Form, forView View, metaData *model.MetaData) Views {
	var labelContent View = Escape(form.FieldLabel(metaData))
	if form.IsFieldRequired(metaData) {
		labelContent = Views{labelContent, form.GetRequiredMarker()}
	}
	label := &Label{For: forView, Content: labelContent}
	return Views{label, forView}
}
