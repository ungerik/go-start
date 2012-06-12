package view

import (
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"strconv"
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

func (self *StandardFormLayout) structFieldInArrayOrSlice(arrayOrSlice, field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	fieldFactory := form.GetFieldFactory()
	// We expect a Table as last form field.
	// If it doesn't exist yet because this is the first visible
	// struct field in the first array field, then create it
	var table *Table
	if len(formFields) > 0 {
		table, _ = formFields[len(formFields)-1].(*Table)
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
		// Also a hidden input to get back the changed length of the table
		// if the user changes it via javascript in the client
		lengthInput := &HiddenInput{
			Name:  arrayOrSlice.Selector() + ".length",
			Value: strconv.Itoa(arrayOrSlice.Value.Len()),
		}
		formFields = append(formFields, lengthInput, table)
		// Add script for manipulating table rows
		context.AddScript(EditFormSliceTableScript, 0)
	}
	tableModel := table.Model.(ViewsTableModel)
	if field.Parent.Index == 0 {
		// If first array field, add label to table header
		tableModel[0] = append(tableModel[0], fieldFactory.NewTableHeader(field, form))
	}
	if tableModel.Rows()-1 == field.Parent.Index {
		// Create row in table model for this array field
		tableModel = append(tableModel, Views{})
		table.Model = tableModel
	}
	// Append form field in last row for this struct field
	row := &tableModel[tableModel.Rows()-1]

	var td View = fieldFactory.NewInput(false, field, form)
	if validationErr != nil {
		td = Views{
			td,
			fieldFactory.NewFieldErrorMessage(validationErr.Error(), field, form),
		}
	}
	*row = append(*row, td)

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
		return self.structFieldInArrayOrSlice(grandParent, field, validationErr, form, context, formFields)
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
	// todo replace. below handels non struct array fields:
	return self.StructField(field, validationErr, form, context, formFields)
	return formFields
}

func (self *StandardFormLayout) EndArray(array *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	return self.endArrayOrSlice(array, formFields)
}

func (self *StandardFormLayout) BeginSlice(slice *model.MetaData, form *Form, context *Context, formFields Views) Views {
	return formFields
}

func (self *StandardFormLayout) SliceField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	// todo replace. below handels non struct array fields:
	return self.StructField(field, validationErr, form, context, formFields)
	return formFields
}

func (self *StandardFormLayout) EndSlice(slice *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views {
	return self.endArrayOrSlice(slice, formFields)
}

func (self *StandardFormLayout) endArrayOrSlice(arrayOrSlice *model.MetaData, formFields Views) Views {
	if len(formFields) > 0 {
		if table, ok := formFields[len(formFields)-1].(*Table); ok {
			tableModel := table.Model.(ViewsTableModel)
			tableModel[0] = append(tableModel[0], HTML("Actions"))
			rows := tableModel.Rows()
			for i := 1; i < rows; i++ {
				firstRow := (i == 1)
				lastRow := (i == rows-1)
				buttons := Views{
					HTML("&nbsp;&nbsp;"),
					&Button{
						Content:  HTML("&uarr;"),
						Disabled: firstRow,
						OnClick:  "gostart_form.moveRowUp(this);",
					},
					&Button{
						Content:  HTML("&darr;"),
						Disabled: lastRow,
						OnClick:  "gostart_form.moveRowDown(this);",
					},
				}
				if arrayOrSlice.Kind == model.SliceKind {
					if lastRow {
						buttons = append(buttons,
							&Button{
								Content: HTML("+"),
								OnClick: "gostart_form.addRow(this);",
							},
						)
					} else {
						buttons = append(buttons,
							&Button{
								Content: HTML("X"),
								OnClick: "gostart_form.removeRow(this)",
							},
						)
					}
				}
				tableModel[i] = append(tableModel[i], buttons)
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

const EditFormSliceTableScript = `
var gostart_form = {

	// function swapAttribs(a, b, attr) {
	// 	var x = a.attr(attr);
	// 	var y = b.attr(attr);
	// 	a.attr(attr, y);
	// 	b.attr(attr, x);
	// }
	// function swapChildren(a, b) {
	// 	var x = a.children().detach();
	// 	var y = b.children().detach();
	// 	x.appendTo(b);
	// 	y.appendTo(a);
	// }

	swapValues: function(a, b) {
		var x = a.val();
		var y = b.val();
		a.val(y);
		b.val(x);
	},

	swapChecked: function(a, b) {
		var x = a.prop("checked");
		var y = b.prop("checked");
		a.prop("checked", y);
		b.prop("checked", x);
	},

	swapRowValues: function(tr0, tr1) {
		var inputs0 = tr0.find("td > :input").not(":button");
		var inputs1 = tr1.find("td > :input").not(":button");
		for (i=0; i < inputs0.length; i++) {
			gostart_form.swapValues(inputs0.eq(i), inputs1.eq(i));
		}	
		inputs0 = tr0.find("td > :checkbox");
		inputs1 = tr1.find("td > :checkbox");
		for (i=0; i < inputs0.length; i++) {
			gostart_form.swapChecked(inputs0.eq(i), inputs1.eq(i));
		}	
	},

	onLengthChanged: function(table) {
		var rows = table.find("tr");
		table.prev("input[type=hidden]").val(rows.length-1);
		rows.each(function(row) {
			var firstRow = (row == 1); // ignore header row
			var lastRow = (row == rows.length-1);
			var buttons = jQuery(this).find("td:last > :button");
			buttons.eq(0).prop("disabled", firstRow);
			buttons.eq(1).prop("disabled", lastRow);
			if (lastRow) {
				buttons.eq(2).attr("onclick", "gostart_form.addRow(this);").text("+");
			} else {
				buttons.eq(2).attr("onclick", "gostart_form.removeRow(this);").text("X");
			}
		});
	},

	removeRow: function(button) {
		if (confirm("Are you sure you want to delete this row?")) {
			var tr = jQuery(button).parents("tr");
			var table = tr.parents("table");

			// Swap all values with following rows to move the values of the
			// row to be deleted to the last row and everything else one row up
			var rows = tr.add(tr.nextAll());
			rows.each(function(i) {
				if (i == 0) return;
				gostart_form.swapRowValues(rows.eq(i-1), rows.eq(i));
			});

			rows.last().remove();

			gostart_form.onLengthChanged(table);
		}
	},

	addRow: function(button) {
		var tr0 = jQuery(button).parents("tr");
		var tr1 = tr0.clone();
		var table = tr0.parents("table");

		// Set correct class for new row
		var numRows = tr0.prevAll().length + 1;
		var evenOdd = (numRows % 2 == 0) ? " even" : " odd";
		tr1.attr("class", "row"+numRows+evenOdd);

		// Correct name attributes of the new row's input elements
		var oldIndex = "."+(numRows-2)+".";
		var newIndex = "."+(numRows-1)+".";
		tr1.find("td > :input").not(":button").each(function(index) {
			var i = this.name.lastIndexOf(oldIndex);
			this.name = this.name.slice(0,i)+newIndex+this.name.slice(i+oldIndex.length);
		});

		tr1.insertAfter(tr0);

		gostart_form.onLengthChanged(table);
	},

	moveRowUp: function(button) {
		var tr1 = jQuery(button).parents("tr");
		var tr0 = tr1.prev();
		gostart_form.swapRowValues(tr0, tr1);
	},

	moveRowDown: function(button) {
		var tr0 = jQuery(button).parents("tr");
		var tr1 = tr0.next();
		gostart_form.swapRowValues(tr0, tr1);
	}
};
`
