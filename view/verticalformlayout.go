package view

import (
	"fmt"
	"github.com/ungerik/go-start/model"
)

/*
VerticalFormLayout puts labels above the input fields, except for checkboxes.

CSS needed form vertical form layout:

	form label:after {
		content: ":";
	}

	form input[type=checkbox] + label:after {
		content: "";
	}

	form label.vertical {
		display: block;
	}
*/
type VerticalFormLayout struct {
}

func (self *VerticalFormLayout) NewField(form *Form, modelValue model.Value, metaData *model.MetaData, disable bool, errors []*model.ValidationError) View {
	if metaData == nil {
		panic(fmt.Sprintf("model.Value must be a struct member to get a label and meta data for the form field. Passed as root model.Value: %T", modelValue))
	}

	switch s := modelValue.(type) {
	case *model.Bool:
		value := s.Get()
		checkbox := &Checkbox{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Label:    getLabel(metaData),
			Disabled: disable,
			Checked:  value,
		}
		return &Paragraph{Content: checkbox}

	case *model.Choice:
		choice := s
		selectModel := &StringsSelectModel{choice.Options(metaData), choice.Get()}
		sel := &Select{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Model:    selectModel,
			Disabled: disable,
			Size:     1,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, sel)

	case *model.DynamicChoice:
		dynamicChoice := s
		selectModel := &IndexedStringsSelectModel{dynamicChoice.Options(), dynamicChoice.Index()}
		sel := &Select{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Model:    selectModel,
			Disabled: disable,
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
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     date.Get(),
			Size:     len(model.DateFormat),
			Disabled: disable,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField, HTML("(Format: "+model.DateFormat+")<br/>"))

	case *model.DateTime:
		dateTime := s
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     dateTime.Get(),
			Size:     len(model.DateTimeFormat),
			Disabled: disable,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField, HTML("(Format: "+model.DateTimeFormat+")<br/>"))

	case *model.Email:
		value := s.Get()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Type:     EmailTextField,
			Text:     value,
			Size:     40,
			Disabled: disable,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Float:
		str := s
		if str.Hidden(metaData) {
			return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
		}
		value := modelValue.(fmt.Stringer).String()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Disabled: disable,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Int:
		str := s
		if str.Hidden(metaData) {
			return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
		}
		value := modelValue.(fmt.Stringer).String()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Disabled: disable,
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
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Type:     PasswordTextField,
			Text:     value,
			Size:     40,
			Disabled: disable,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textField)

	case *model.Phone:
		value := s.Get()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Size:     20,
			Disabled: disable,
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
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     str.Get(),
			Size:     80,
			Disabled: disable,
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
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     text.Get(),
			Cols:     cols,
			Rows:     rows,
			Disabled: disable,
		}
		return self.addLabelToEditorWidget(form, modelValue, metaData, errors, textArea)

	case *model.Url:
		value := s.Get()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Size:     80,
			Disabled: disable,
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
	label := Views{Escape(getLabel(metaData))}
	if form.IsFieldRequired(metaData) {
		label = append(label, SPAN("required", HTML("*")))
	}
	views = append(views, &Label{Class: "vertical", Content: label, For: editorView})
	views = append(views, extraLabels...)
	for _, error := range errors {
		views = append(
			views,
			SPAN(form.GetErrorMessageClass(), Escape(error.WrappedError.Error())),
			//BR(),
		)
	}
	views = append(views, editorView)

	return &Paragraph{Content: views}
}
