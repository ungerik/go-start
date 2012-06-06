package view

import (
	"fmt"
	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
)

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

func (self *StandardFormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View) {
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
		return checkbox

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
		if len(options) == 0 || options[0] != "" {
			options = append([]string{""}, options...)
		}
		input = &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    &IndexedStringsSelectModel{options, s.Index()},
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
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     len(model.DateFormat),
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.DateTime:
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     len(model.DateTimeFormat),
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Email:
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Type:     EmailTextField,
			Text:     s.Get(),
			Size:     40,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Float:
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.String(),
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Int:
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.String(),
			Disabled: form.IsFieldDisabled(metaData),
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
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Type:     PasswordTextField,
			Text:     s.Get(),
			Size:     40,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Phone:
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     20,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case model.Reference:
		if !form.ShowRefIDs {
			return nil
		}
		value := s.String()
		if value == "" {
			value = "[empty]"
		}
		input = HTML(value)

	case *model.String:
		textField := &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     80,
			Disabled: form.IsFieldDisabled(metaData),
		}
		if maxlen, ok, _ := s.Maxlen(metaData); ok {
			textField.Size = maxlen
			textField.MaxLength = maxlen
		}
		input = textField

	case *model.Text:
		cols, _, _ := s.Cols(metaData) // will be zero if not available, which is OK
		rows, _, _ := s.Rows(metaData) // will be zero if not available, which is OK
		input = &TextArea{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Cols:     cols,
			Rows:     rows,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Url:
		input = &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     80,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.GeoLocation:
		input = HTML(s.String())

	default:
		panic(fmt.Sprintf("Unsupported model.Value type %T", metaData.Value.Addr().Interface()))
	}

	if withLabel {
		return Views{self.NewLabel(input, metaData, form), input}
	}
	return input
}

func (self *StandardFormFieldFactory) NewHiddenInput(metaData *model.MetaData, form *Form) View {
	return &HiddenInput{
		Name:  metaData.Selector(),
		Value: fmt.Sprintf("%v", metaData.Value.Interface()),
	}
}

func (self *StandardFormFieldFactory) NewLabel(forView View, metaData *model.MetaData, form *Form) View {
	// todo add extra label for date/time
	// HTML("(Format: "+model.DateFormat+")<br/>")
	// HTML("(Format: "+model.DateTimeFormat+")<br/>")
	var labelContent View = Escape(form.FieldLabel(metaData))
	if form.IsFieldRequired(metaData) {
		labelContent = Views{labelContent, form.GetRequiredMarker()}
	}
	return &Label{For: forView, Content: labelContent}
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

func (self *StandardFormFieldFactory) NewSubmitButton(text string, form *Form) View {
	return &Button{Class: form.GetSubmitButtonClass(), Submit: true, Value: text}
}

func (self *StandardFormFieldFactory) NewAddSliceItemButton(form *Form) View {
	return &Button{Value: "+"} // todo
}

func (self *StandardFormFieldFactory) NewRemoveSliceItemButton(form *Form) View {
	return &Button{Value: "-"} // todo
}
