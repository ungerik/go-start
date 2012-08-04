package view

import (
	"fmt"
	"strconv"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
)

func AddStandardLabel(form *Form, forView View, metaData *model.MetaData) Views {
	var labelContent View = Escape(form.FieldLabel(metaData))
	if form.IsFieldRequired(metaData) {
		labelContent = Views{labelContent, form.GetRequiredMarker()}
	}
	label := &Label{For: forView, Content: labelContent}
	return Views{label, forView}
}

type StandardFormFieldFactory struct {
}

func (self *StandardFormFieldFactory) CanCreateInput(metaData *model.MetaData, form *Form) bool {
	switch metaData.Value.Addr().Interface().(type) {
	case model.Reference:
		return false
	case model.Value:
		return true
	}
	return false
}

func (self *StandardFormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	switch s := metaData.Value.Addr().Interface().(type) {
	case *model.Bool:
		checkbox := &Checkbox{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Disabled: form.IsFieldDisabled(metaData),
			Checked:  s.Get(),
		}
		if withLabel {
			checkbox.Label = form.FieldLabel(metaData)
		}
		return checkbox, nil

	case *model.Choice:
		options := s.Options(metaData)
		if len(options) == 0 || options[0] != "" {
			options = append([]string{""}, options...)
		}
		input = &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    &StringsSelectModel{options, s.Get()},
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}

	case *model.DynamicChoice:
		options := s.Options()
		index := s.Index()
		if len(options) == 0 || options[0] != "" {
			options = append([]string{""}, options...)
			index++
		}
		input = &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    &IndexedStringsSelectModel{options, index},
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}

	case *model.Country:
		// todo
		value := s.String()
		if value == "" {
			value = "[empty]"
		}
		input = Escape(value)

	case *model.Date:
		input = Views{
			HTML("(Format: " + model.DateFormat + ")<br/>"),
			&TextField{
				Class:       form.FieldInputClass(metaData),
				Name:        metaData.Selector(),
				Text:        s.Get(),
				Size:        len(model.DateFormat),
				Disabled:    form.IsFieldDisabled(metaData),
				Placeholder: form.InputFieldPlaceholder(metaData),
			},
		}

	case *model.DateTime:
		input = Views{
			HTML("(Format: " + model.DateTimeFormat + ")<br/>"),
			&TextField{
				Class:       form.FieldInputClass(metaData),
				Name:        metaData.Selector(),
				Text:        s.Get(),
				Size:        len(model.DateTimeFormat),
				Disabled:    form.IsFieldDisabled(metaData),
				Placeholder: form.InputFieldPlaceholder(metaData),
			},
		}

	case *model.Email:
		input = &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Type:        EmailTextField,
			Text:        s.Get(),
			Size:        form.GetInputSize(metaData),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case *model.Float:
		input = &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        s.String(),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case *model.Int:
		input = &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        s.String(),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case *model.Language:
		// todo
		value := s.String()
		if value == "" {
			value = "[empty]"
		}
		input = HTML(value)

	case *model.Password:
		input = &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Type:        PasswordTextField,
			Text:        s.Get(),
			Size:        form.GetInputSize(metaData),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case *model.Phone:
		input = &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        s.Get(),
			Size:        form.GetInputSize(metaData),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case model.Reference:
		if !form.ShowRefIDs {
			return nil, nil
		}
		value := s.String()
		if value == "" {
			value = "[empty]"
		}
		input = HTML(value)

	case *model.String:
		textField := &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        s.Get(),
			Size:        form.GetInputSize(metaData),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}
		if maxlen, ok, _ := s.Maxlen(metaData); ok {
			textField.Size = maxlen
			textField.MaxLength = maxlen
		}
		input = textField

	case *model.Text:
		var cols int // will be zero if not available, which is OK
		if str, ok := metaData.Attrib(StructTagKey, "cols"); ok {
			cols, err = strconv.Atoi(str)
			if err != nil {
				panic("Error in StandardFormFieldFactory.NewInput(): " + err.Error())
			}
		}
		var rows int // will be zero if not available, which is OK
		if str, ok := metaData.Attrib(StructTagKey, "rows"); ok {
			rows, err = strconv.Atoi(str)
			if err != nil {
				panic("Error in StandardFormFieldFactory.NewInput(): " + err.Error())
			}
		}
		input = &TextArea{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        s.Get(),
			Cols:        cols,
			Rows:        rows,
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case *model.Url:
		input = &TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        s.Get(),
			Size:        form.GetInputSize(metaData),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
		}

	case *model.GeoLocation:
		input = HTML(s.String())

	case *model.Blob, *model.File:
		input = &FileInput{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Disabled: form.IsFieldDisabled(metaData),
		}

	default:
		return nil, fmt.Errorf("Unsupported model.Value type %T", metaData.Value.Addr().Interface())
	}

	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

func (self *StandardFormFieldFactory) NewHiddenInput(metaData *model.MetaData, form *Form) (View, error) {
	return &HiddenInput{
		Name:  metaData.Selector(),
		Value: fmt.Sprintf("%v", metaData.Value.Interface()),
	}, nil
}

func (self *StandardFormFieldFactory) NewTableHeader(metaData *model.MetaData, form *Form) (View, error) {
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

func (self *StandardFormFieldFactory) NewFieldDescrtiption(description string, form *Form) View {
	return SPAN(form.GetFieldDescriptionClass(), Escape(description))
}

func (self *StandardFormFieldFactory) NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View {
	return SPAN(form.GetErrorMessageClass(), Escape(message))
}

func (self *StandardFormFieldFactory) NewGeneralErrorMessage(message string, form *Form) View {
	return SPAN(form.GetErrorMessageClass(), Escape(message))
}

func (self *StandardFormFieldFactory) NewSuccessMessage(message string, form *Form) View {
	return SPAN(form.GetSuccessMessageClass(), Escape(message))
}

func (self *StandardFormFieldFactory) NewSubmitButton(text, confirmationMessage string, form *Form) View {
	return &SubmitButton{
		Class:          form.GetSubmitButtonClass(),
		Value:          text,
		OnClickConfirm: confirmationMessage,
	}
}

func (self *StandardFormFieldFactory) NewAddButton(onclick string, form *Form) View {
	return &Button{
		Content: HTML("+"),
		OnClick: onclick,
		Class:   "form-table-button",
	}
}

func (self *StandardFormFieldFactory) NewRemoveButton(onclick string, form *Form) View {
	return &Button{
		Content: HTML("X"),
		OnClick: onclick,
		Class:   "form-table-button",
	}
}

func (self *StandardFormFieldFactory) NewUpButton(disabled bool, onclick string, form *Form) View {
	return &Button{
		Content:  HTML("&uarr;"),
		Disabled: disabled,
		OnClick:  onclick,
		Class:    "form-table-button",
	}
}

func (self *StandardFormFieldFactory) NewDownButton(disabled bool, onclick string, form *Form) View {
	return &Button{
		Content:  HTML("&darr;"),
		Disabled: disabled,
		OnClick:  onclick,
		Class:    "form-table-button",
	}
}
