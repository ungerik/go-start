package view

import (
	"fmt"
	"github.com/ungerik/go-start/model"
)

type StandardFormFieldFactory struct {
}

func (self *StandardFormFieldFactory) NewInput(form *Form, data interface{}, metaData *model.MetaData) View {
	if modelValue, ok := data.(model.Value); ok {
		switch s := modelValue.(type) {
		case *model.Bool:
			value := s.Get()
			return &Checkbox{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Label:    getLabel(metaData),
				Disabled: form.IsFieldDisabled(metaData),
				Checked:  value,
			}

		case *model.Choice:
			choice := s
			selectModel := &StringsSelectModel{choice.Options(metaData), choice.Get()}
			return &Select{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Model:    selectModel,
				Disabled: form.IsFieldDisabled(metaData),
				Size:     1,
			}

		case *model.DynamicChoice:
			dynamicChoice := s
			selectModel := &IndexedStringsSelectModel{dynamicChoice.Options(), dynamicChoice.Index()}
			return &Select{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Model:    selectModel,
				Disabled: form.IsFieldDisabled(metaData),
				Size:     1,
			}

		case *model.Country:
			// todo
			value := modelValue.(fmt.Stringer).String()
			if value == "" {
				value = "[empty]"
			}
			return Escape(value)

		case *model.Date:
			date := s
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     date.Get(),
				Size:     len(model.DateFormat),
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.DateTime:
			dateTime := s
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     dateTime.Get(),
				Size:     len(model.DateTimeFormat),
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.Email:
			value := s.Get()
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Type:     EmailTextField,
				Text:     value,
				Size:     40,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.Float:
			str := s
			if str.Hidden(metaData) {
				return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
			}
			value := modelValue.(fmt.Stringer).String()
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     value,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.Int:
			str := s
			if str.Hidden(metaData) {
				return &HiddenInput{Name: metaData.Selector(), Value: str.String()}
			}
			value := modelValue.(fmt.Stringer).String()
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     value,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.Language:
			// todo
			value := modelValue.(fmt.Stringer).String()
			if value == "" {
				value = "[empty]"
			}
			return HTML(value)

		case *model.Password:
			value := s.Get()
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Type:     PasswordTextField,
				Text:     value,
				Size:     40,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.Phone:
			value := s.Get()
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     value,
				Size:     20,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case model.Reference:
			if !form.ShowRefIDs {
				return nil
			}
			value := modelValue.(fmt.Stringer).String()
			if value == "" {
				value = "[empty]"
			}
			return HTML(value)

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
				Disabled: form.IsFieldDisabled(metaData),
			}
			if maxlen, ok, _ := str.Maxlen(metaData); ok {
				textField.Size = maxlen
				textField.MaxLength = maxlen
			}
			return textField

		case *model.Text:
			text := s
			cols, _, _ := text.Cols(metaData) // will be zero if not available, which is OK
			rows, _, _ := text.Rows(metaData) // will be zero if not available, which is OK
			return &TextArea{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     text.Get(),
				Cols:     cols,
				Rows:     rows,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.Url:
			value := s.Get()
			return &TextField{
				Class:    getClass(metaData),
				Name:     metaData.Selector(),
				Text:     value,
				Size:     80,
				Disabled: form.IsFieldDisabled(metaData),
			}

		case *model.GeoLocation:
			value := s.String()
			return HTML(value)
		}

		panic(fmt.Sprintf("Unsupported model.Value type %T", modelValue))
	}
	return nil
}

func (self *StandardFormFieldFactory) NewLabel(form *Form, forView View, data interface{}, metaData *model.MetaData) View {
	// todo add extra label for date/time
	// HTML("(Format: "+model.DateFormat+")<br/>")
	// HTML("(Format: "+model.DateTimeFormat+")<br/>")
	var labelContent View = Escape(getLabel(metaData))
	if form.IsFieldRequired(metaData) {
		labelContent = Views{labelContent, form.GetRequiredMarker()}
	}
	return &Label{For: forView, Content: labelContent}
}

func (self *StandardFormFieldFactory) NewFieldErrorMessage(form *Form, message string, metaData *model.MetaData) View {
	return nil
}

func (self *StandardFormFieldFactory) NewFormErrorMessage(form *Form, message string) View {
	return nil
}

func (self *StandardFormFieldFactory) NewSuccessMessage(form *Form, message string) View {
	return nil
}
