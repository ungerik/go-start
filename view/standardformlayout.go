package view

import (
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
)

const EditFormSliceTableScript = `
function swapAttribs(a, b, attr) {
	var x = a.attr(attr);
	var y = b.attr(attr);
	a.attr(attr, y);
	b.attr(attr, x);
}

function swapNames(tr0, tr1) {
	// Swap tr classes
	swapAttribs(tr0, tr1, "class");

	// http://stackoverflow.com/questions/5897084/jquery-add-form-input-and-change-input-name
	var inputs0 = tr0.find("td > *");
	var inputs1 = tr1.find("td > *");
	// for (i=0; i < inputs0.length-3; i++) {
	// 	swapAttribs(inputs0[i], inputs1[i], "name");
	// }
}

function removeRow(button) {
	if (confirm("Are you sure you want to delete this row?")) {
		var tr = jQuery(button).parents("tr");
		var table = tr.parents("table");
		tr.remove();
	}
}

function addRow(button) {
}

function moveRowUp(button) {
	var tr1 = jQuery(button).parents("tr");
	var tr0 = tr1.prev();
	swapNames(tr0, tr1);
	tr1.after(tr0);
}

function moveRowDown(button) {
	var tr0 = jQuery(button).parents("tr");
	var tr1 = tr0.next();
	swapNames(tr0, tr1);
	tr1.after(tr0);
}
`

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

func (self *StandardFormLayout) BeginFormContent(form *Form, context *Context, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) SubmitSuccess(message string, form *Form, context *Context, formFields Views) Views {
	return append(formFields, form.GetFieldFactory().NewSuccessMessage(message, form))
}

func (self *StandardFormLayout) SubmitError(message string, form *Form, context *Context, formFields Views) Views {
	return append(formFields, form.GetFieldFactory().NewGeneralErrorMessage(message, form))
}

func (self *StandardFormLayout) EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, context *Context, formFields Views) Views {
	fieldFactory := form.GetFieldFactory()
	for _, err := range generalValidationErrs {
		formFields = append(formFields, fieldFactory.NewGeneralErrorMessage(err.Error(), form))
		formFields = append(Views{fieldFactory.NewGeneralErrorMessage(err.Error(), form)}, formFields...)
	}
	formId := &HiddenInput{Name: FormIDName, Value: form.FormID}
	submitButton := fieldFactory.NewSubmitButton(form.GetSubmitButtonText(), form.SubmitButtonConfirm, form)
	return append(formFields, formId, submitButton)
}

func (self *StandardFormLayout) BeginStruct(strct *model.MetaData, form *Form, context *Context, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) StructField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	fieldFactory := form.GetFieldFactory()
	if !fieldFactory.CanCreateInput(field, form) || form.IsFieldExcluded(field) {
		return formFields
	}

	grandParent := field.Parent.Parent
	if grandParent != nil && (grandParent.Kind == model.ArrayKind || grandParent.Kind == model.SliceKind) {
		if form.IsFieldExcluded(grandParent) {
			return formFields
		}
		// We expect a Table as last form field.
		// If it doesn't exist yet because this is the first visible
		// struct field in the first array field, then create it
		var table *Table
		if len(formFields) > 0 {
			table, _ = formFields[len(formFields)-1].(*Table)
		}
		if table == nil {
			// First struct field of first array field, create table and table model
			table = &Table{
				Caption:   form.FieldLabel(grandParent),
				HeaderRow: true,
				Model:     ViewsTableModel{Views{}},
			}
			table.Init(table) // get an ID now
			formFields = append(formFields, table)
			// Add script for manipulating table rows
			context.AddScript(EditFormSliceTableScript, 0)
		}
		tableModel := table.Model.(ViewsTableModel)
		if field.Parent.Index == 0 {
			// If first array field, add label to table header
			tableModel[0] = append(tableModel[0], Escape(form.DirectFieldLabel(field)))
		}
		if tableModel.Rows()-1 == field.Parent.Index {
			// Create row in table model for this array field
			tableModel = append(tableModel, Views{})
			table.Model = tableModel
		}
		// Append form field in last row for this struct field
		row := &tableModel[tableModel.Rows()-1]
		*row = append(*row, fieldFactory.NewInput(false, field, form))

		return formFields
	}

	if form.IsFieldHidden(field) {
		return append(formFields, fieldFactory.NewHiddenInput(field, form))
	}
	var formField View = fieldFactory.NewInput(true, field, form)
	if validationErr != nil {
		formField = Views{formField, fieldFactory.NewFieldErrorMessage(validationErr.Error(), field, form)}
	}
	return append(formFields, DIV(Config.Form.StandardFormLayoutDivClass, formField))
}

func (self *StandardFormLayout) EndStruct(strct *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) BeginArray(array *model.MetaData, form *Form, context *Context, formFields Views) Views {
	debug.Nop()
	return formFields
}

func (self *StandardFormLayout) ArrayField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	return self.StructField(field, validationErr, form, context, formFields) // todo replace
	return formFields
}

func (self *StandardFormLayout) EndArray(array *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) BeginSlice(slice *model.MetaData, form *Form, context *Context, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) SliceField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	return self.StructField(field, validationErr, form, context, formFields) // todo replace
	return formFields
}

func (self *StandardFormLayout) EndSlice(slice *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	if len(formFields) > 0 {
		if table, ok := formFields[len(formFields)-1].(*Table); ok {
			tableModel := table.Model.(ViewsTableModel)
			tableModel[0] = append(tableModel[0], HTML("Actions"))
			rows := tableModel.Rows()
			for i := 1; i < rows; i++ {
				//firstRow := i == 1
				lastRow := false //i == rows-1
				tableModel[i] = append(
					tableModel[i],
					Views{
						HTML("&nbsp;&nbsp;"),
						&Button{
							Content: HTML("&uarr;"),
							// Disabled: firstRow,
							OnClick: "moveRowUp(this);",
						},
						&Button{
							Content: HTML("&darr;"),
							// Disabled: lastRow,
							OnClick: "moveRowDown(this);",
						},
						&If{
							Condition: lastRow,
							Content: &Button{
								Content: HTML("+"),
								OnClick: "addRow(this);",
							},
							ElseContent: &Button{
								Content: HTML("X"),
								OnClick: "removeRow(this)",
							},
						},
					},
				)
			}
		}
	}
	return formFields
}

func (self *StandardFormLayout) fieldNeedsLabel(field *model.MetaData) bool {
	switch field.Value.Addr().Interface().(type) {
	case *model.Bool:
		return false
	}
	return true
}
