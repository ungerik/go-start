package view

import (
	"fmt"
	"github.com/ungerik/go-start/model"
)

type StandardFormFieldFactory struct {
}

func (self *StandardFormFieldFactory) NewInput(metaData *model.MetaData, form *Form) View {
	switch s := metaData.Value.Addr().Interface().(type) {
	case *model.Bool:
		return &Checkbox{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Label:    form.FieldLabel(metaData),
			Disabled: form.IsFieldDisabled(metaData),
			Checked:  s.Get(),
		}

	case *model.Choice:
		return &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    &StringsSelectModel{s.Options(metaData), s.Get()},
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}

	case *model.DynamicChoice:
		return &Select{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Model:    &IndexedStringsSelectModel{s.Options(), s.Index()},
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}

	case *model.Country:
		// todo
		value := s.String()
		if value == "" {
			value = "[empty]"
		}
		return Escape(value)

	case *model.Date:
		return &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     len(model.DateFormat),
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.DateTime:
		return &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     len(model.DateTimeFormat),
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Email:
		return &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Type:     EmailTextField,
			Text:     s.Get(),
			Size:     40,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Float:
		return &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.String(),
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Int:
		return &TextField{
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
		return HTML(value)

	case *model.Password:
		return &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Type:     PasswordTextField,
			Text:     s.Get(),
			Size:     40,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Phone:
		return &TextField{
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
		return HTML(value)

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
		return textField

	case *model.Text:
		cols, _, _ := s.Cols(metaData) // will be zero if not available, which is OK
		rows, _, _ := s.Rows(metaData) // will be zero if not available, which is OK
		return &TextArea{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Cols:     cols,
			Rows:     rows,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.Url:
		return &TextField{
			Class:    form.FieldInputClass(metaData),
			Name:     metaData.Selector(),
			Text:     s.Get(),
			Size:     80,
			Disabled: form.IsFieldDisabled(metaData),
		}

	case *model.GeoLocation:
		return HTML(s.String())
	}

	panic(fmt.Sprintf("Unsupported model.Value type %T", metaData.Value.Addr().Interface()))
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
	return DIV(form.GetFieldDescriptionClass(), Escape(description))
}

func (self *StandardFormFieldFactory) NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View {
	return DIV(form.GetErrorMessageClass(), Escape(message))
}

func (self *StandardFormFieldFactory) NewGeneralErrorMessage(message string, form *Form) View {
	return DIV(form.GetErrorMessageClass(), Escape(message))
}

func (self *StandardFormFieldFactory) NewSuccessMessage(message string, form *Form) View {
	return DIV(form.GetSuccessMessageClass(), Escape(message))
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
