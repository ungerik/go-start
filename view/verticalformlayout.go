package view

import (
	"fmt"
	"github.com/ungerik/go-start/model"
	"reflect"
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
}

func (self *VerticalFormLayout) NewField_old(form *Form, modelValue model.Value, metaData *model.MetaData, errors []*model.ValidationError) View {
	switch s := modelValue.(type) {
	case *model.Bool:
		value := s.Get()
		checkbox := &Checkbox{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Label:    form.FieldLabel(metaData),
			Disabled: form.IsFieldDisabled(metaData),
			Checked:  value,
		}
		return &Paragraph{Content: checkbox}

	case *model.Choice:
		choice := s
		selectModel := &StringsSelectModel{choice.Options(metaData), choice.Get()}
		sel := &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    selectModel,
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, sel)

	case *model.DynamicChoice:
		dynamicChoice := s
		selectModel := &IndexedStringsSelectModel{dynamicChoice.Options(), dynamicChoice.Index()}
		sel := &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    selectModel,
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, sel)

	case *model.Country:
		// todo
		value := modelValue.(fmt.Stringer).String()
		if value == "" {
			value = "[empty]"
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, Escape(value))

	case *model.Date:
		date := s
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     date.Get(),
			Size:     len(model.DateFormat),
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField, HTML("(Format: "+model.DateFormat+")<br/>"))

	case *model.DateTime:
		dateTime := s
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     dateTime.Get(),
			Size:     len(model.DateTimeFormat),
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField, HTML("(Format: "+model.DateTimeFormat+")<br/>"))

	case *model.Email:
		value := s.Get()
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Type:     EmailTextField,
			Text:     value,
			Size:     40,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Float:
		str := s
		if str.Hidden(metaData) {
			return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
		}
		value := modelValue.(fmt.Stringer).String()
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Int:
		str := s
		if str.Hidden(metaData) {
			return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
		}
		value := modelValue.(fmt.Stringer).String()
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Language:
		// todo
		value := modelValue.(fmt.Stringer).String()
		if value == "" {
			value = "[empty]"
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, HTML(value))

	case *model.Password:
		value := s.Get()
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Type:     PasswordTextField,
			Text:     value,
			Size:     40,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Phone:
		value := s.Get()
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Size:     20,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case model.Reference:
		if !form.ShowRefIDs {
			return nil
		}
		value := modelValue.(fmt.Stringer).String()
		if value == "" {
			value = "[empty]"
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, HTML(value))

	case *model.String:
		str := s
		if str.Hidden(metaData) {
			return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
		}
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     str.Get(),
			Size:     80,
			Disabled: form.IsFieldDisabled(metaData),
		}
		if maxlen, ok, _ := str.Maxlen(metaData); ok {
			textField.Size = maxlen
			textField.MaxLength = maxlen
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Text:
		text := s
		cols, _, _ := text.Cols(metaData) // will be zero if not available, which is OK
		rows, _, _ := text.Rows(metaData) // will be zero if not available, which is OK
		textArea := &TextArea{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     text.Get(),
			Cols:     cols,
			Rows:     rows,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textArea)

	case *model.Url:
		value := s.Get()
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Size:     80,
			Disabled: form.IsFieldDisabled(metaData),
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.GeoLocation:
		value := s.String()
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, HTML(value))
	}

	panic(fmt.Sprintf("Unsupported model.Value type %T", modelValue))
}

func (self *VerticalFormLayout) addLabelToEditorWidget(form *Form, modelValue model.Value, metaData *model.MetaData, errors []*model.ValidationError, editorView View, extraLabels ...View) View {
	views := make(Views, 0, 2+len(errors)*2+1)
	var labelContent View = Escape(form.FieldLabel(metaData))
	if form.IsFieldRequired(metaData) {
		labelContent = Views{labelContent, SPAN("required", HTML("*"))}
	}
	views = append(views, &Label{Class: "vertical", Content: labelContent, For: editorView})
	views = append(views, extraLabels...)
	for _, error := range errors {
		views = append(views, SPAN(form.GetErrorMessageClass(), Escape(error.WrappedError.Error())))
	}
	views = append(views, editorView)
	return P(views)
}

func (self *VerticalFormLayout) BeforeFormContent(form *Form) View {
	return nil
}

func (self *VerticalFormLayout) AfterFormContent(form *Form) View {
	return nil
}

func (self *VerticalFormLayout) BeforeStruct(form *Form, strct reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) StructField(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) AfterStruct(form *Form, strct reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) BeforeArray(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) ArrayField(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) AfterArray(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) BeforeSlice(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) SliceField(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) AfterSlice(form *Form, field reflect.Value, metaData *model.MetaData) View {
	return nil
}
