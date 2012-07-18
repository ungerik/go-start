package view

import (
	// "github.com/ungerik/go-start/debug"
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

func (self *StandardFormLayout) BeginFormContent(form *Form, context *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) SubmitSuccess(message string, form *Form, context *Context, formContent *Views) error {
	*formContent = append(*formContent, form.GetFieldFactory().NewSuccessMessage(message, form))
	return nil
}

func (self *StandardFormLayout) SubmitError(message string, form *Form, context *Context, formContent *Views) error {
	*formContent = append(*formContent, form.GetFieldFactory().NewGeneralErrorMessage(message, form))
	return nil
}

func (self *StandardFormLayout) EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, context *Context, formContent *Views) error {
	fieldFactory := form.GetFieldFactory()
	for _, err := range generalValidationErrs {
		*formContent = append(*formContent, fieldFactory.NewGeneralErrorMessage(err.Error(), form))
		*formContent = append(Views{fieldFactory.NewGeneralErrorMessage(err.Error(), form)}, *formContent...)
	}
	submitButton := fieldFactory.NewSubmitButton(form.GetSubmitButtonText(), form.SubmitButtonConfirm, form)
	*formContent = append(*formContent, submitButton)
	return nil
}

func (self *StandardFormLayout) BeginStruct(strct *model.MetaData, form *Form, context *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) StructField(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	fieldFactory := form.GetFieldFactory()
	if !fieldFactory.CanCreateInput(field, form) || form.IsFieldExcluded(field) {
		return nil
	}

	grandParent := field.Parent.Parent
	if grandParent != nil && (grandParent.Kind == model.ArrayKind || grandParent.Kind == model.SliceKind) {
		if form.IsFieldExcluded(grandParent) {
			return nil
		}
		return self.structFieldInArrayOrSlice(grandParent, field, validationErr, form, context, formContent)
	}

	if form.IsFieldHidden(field) {
		input, err := fieldFactory.NewHiddenInput(field, form)
		if err != nil {
			return err
		}
		*formContent = append(*formContent, input)
		return nil
	}

	formField, err := fieldFactory.NewInput(true, field, form)
	if err != nil {
		return err
	}
	if validationErr != nil {
		formField = Views{formField, fieldFactory.NewFieldErrorMessage(validationErr.Error(), field, form)}
	}
	*formContent = append(*formContent, DIV(Config.Form.StandardFormLayoutDivClass, formField))
	return nil
}

func (self *StandardFormLayout) EndStruct(strct *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) BeginArray(array *model.MetaData, form *Form, context *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) ArrayField(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	if field.Kind == model.ValueKind && !form.IsFieldExcluded(field.Parent) && form.GetFieldFactory().CanCreateInput(field, form) {
		return self.arrayOrSliceFieldValue(field, validationErr, form, context, formContent)
	}
	return nil
}

func (self *StandardFormLayout) EndArray(array *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	return self.endArrayOrSlice(array, form, formContent)
}

func (self *StandardFormLayout) BeginSlice(slice *model.MetaData, form *Form, context *Context, formContent *Views) error {
	return nil
}

func (self *StandardFormLayout) SliceField(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	if field.Kind == model.ValueKind && !form.IsFieldExcluded(field.Parent) && form.GetFieldFactory().CanCreateInput(field, form) {
		return self.arrayOrSliceFieldValue(field, validationErr, form, context, formContent)
	}
	return nil
}

func (self *StandardFormLayout) EndSlice(slice *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	return self.endArrayOrSlice(slice, form, formContent)
}

func (self *StandardFormLayout) structFieldInArrayOrSlice(arrayOrSlice, field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	fieldFactory := form.GetFieldFactory()
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
		context.AddScript(EditFormSliceTableScript, 0)
	}
	tableModel := table.Model.(ViewsTableModel)
	if field.Parent.Index == 0 {
		// If first array field, add label to table header
		header, err := fieldFactory.NewTableHeader(field, form)
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

	td, err := fieldFactory.NewInput(false, field, form)
	if err != nil {
		return err
	}
	if validationErr != nil {
		td = Views{
			td,
			fieldFactory.NewFieldErrorMessage(validationErr.Error(), field, form),
		}
	}
	*row = append(*row, td)
	return nil
}

func (self *StandardFormLayout) arrayOrSliceFieldValue(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error {
	arrayOrSlice := field.Parent
	fieldFactory := form.GetFieldFactory()
	// We expect a Table as last form content field.
	// If it doesn't exist yet because this is the first visible
	// struct field in the first array field, then create it
	var table *Table
	if len(*formContent) > 0 {
		table, _ = (*formContent)[len(*formContent)-1].(*Table)
	}
	if table == nil {
		// First array/slice field, create table and table model.
		header, err := fieldFactory.NewTableHeader(arrayOrSlice, form)
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
		context.AddScript(EditFormSliceTableScript, 0)
	}
	td, err := fieldFactory.NewInput(false, field, form)
	if err != nil {
		return err
	}
	if validationErr != nil {
		td = Views{
			td,
			fieldFactory.NewFieldErrorMessage(validationErr.Error(), field, form),
		}
	}
	table.Model = append(table.Model.(ViewsTableModel), Views{td})
	return nil
}

func (self *StandardFormLayout) endArrayOrSlice(arrayOrSlice *model.MetaData, form *Form, formContent *Views) error {
	if len(*formContent) > 0 {
		// Add "Actions" column with buttons to table with slice or array values
		if table, ok := (*formContent)[len(*formContent)-1].(*Table); ok {
			fieldFactory := form.GetFieldFactory()
			tableModel := table.Model.(ViewsTableModel)
			tableModel[0] = append(tableModel[0], HTML("Actions"))
			rows := tableModel.Rows()
			for i := 1; i < rows; i++ {
				// todo: script depends on buttons being HTML buttons,
				// but buttons are created by form field factory that doesn't
				// guarantee that. needs decoupling via css classes.
				firstRow := (i == 1)
				lastRow := (i == rows-1)
				buttons := Views{
					fieldFactory.NewUpButton(firstRow, "gostart_form.moveRowUp(this);", form),
					fieldFactory.NewDownButton(lastRow, "gostart_form.moveRowDown(this);", form),
				}
				if arrayOrSlice.Kind == model.SliceKind {
					if lastRow {
						buttons = append(buttons, fieldFactory.NewAddButton("gostart_form.addRow(this);", form))
					} else {
						buttons = append(buttons, fieldFactory.NewRemoveButton("gostart_form.removeRow(this)", form))
					}
				}
				tableModel[i] = append(tableModel[i], buttons)
			}
			// Add a hidden input to get back the changed length of the table
			// if the user changes it via javascript in the client
			*formContent = append(*formContent,
				&HiddenInput{
					Name:  arrayOrSlice.Selector() + ".length",
					Value: strconv.Itoa(arrayOrSlice.Value.Len()),
				},
			)
		}
	}
	return nil
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
		var oldIndex = "."+(numRows-2);
		var newIndex = "."+(numRows-1);
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
