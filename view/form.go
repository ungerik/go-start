package view

import (
	"bytes"
	"fmt"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
	"strings"
	//	"github.com/ungerik/go-start/debug"
)

type GetFormModelFunc func(form *Form, context *Context) (model interface{}, err error)
type OnSubmitFormFunc func(form *Form, formModel interface{}, context *Context) error

func FormModel(model interface{}) GetFormModelFunc {
	return func(form *Form, context *Context) (interface{}, error) {
		return model, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// Form

type Form struct {
	ViewBaseWithId
	Class               string
	Action              string
	Content             View
	Method              string
	FormID              string
	GetModel            GetFormModelFunc
	ModelMaxDepth       int      // if zero, no depth limit
	HideFields          []string // Use point notation for nested fields
	DisableFields       []string // Use point notation for nested fields
	OnSubmit            OnSubmitFormFunc
	ErrorMessageClass   string // If empty, Config.FormErrorMessageClass will be used
	SuccessMessageClass string // If empty, Config.FormSuccessMessageClass will be used
	SuccessMessage      string
	ButtonText          string
	ButtonClass         string
	Redirect            URL // 302 redirect after successful Save()
	ShowRefIDs          bool
}

func (self *Form) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Form) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.OnSubmit != nil {
		// Determine if it's a POST request for this form:
		isPOST := false
		if context.Request.Method == "POST" {
			// Every HTML form gets an ID to allow more than one form per page:
			id, ok := context.Params["form_id"]
			if ok && id == self.FormID {
				isPOST = true
			}
		}

		// Create views for form fields:

		// Set a view before and after the form fields
		// if there is an error or success message
		// (won't be rendered if nil)
		// Also add a hidden form field with the form id
		// and a submit button
		views := make(Views, 1, 32)

		numValueErrors := 0
		generalErrors := []*model.ValidationError{}

		var formModel interface{}

		if self.GetModel != nil {
			formModel, err = self.GetModel(self, context)
			if err != nil {
				return err
			}
			model.WalkStructure(formModel, self.ModelMaxDepth, func(data interface{}, metaData model.MetaData) {
				if modelValue, ok := data.(model.Value); ok {
					selector := metaData.Selector()
					arraySelector := metaData.ArrayWildcardSelector()

					if utils.StringIn(selector, self.HideFields) || utils.StringIn(arraySelector, self.HideFields) {
						return
					}

					var valueErrors []*model.ValidationError
					if isPOST {
						formValue, ok := context.Params[selector]
						if b, isBool := modelValue.(*model.Bool); isBool {
							b.Set(formValue != "")
						} else if ok {
							err = modelValue.SetString(formValue)
							if err == nil {
								valueErrors = modelValue.Validate(metaData)
							} else {
								valueErrors = model.NewValidationErrors(err, metaData)
							}
							numValueErrors += len(valueErrors)
						}
					}

					disable := utils.StringIn(selector, self.DisableFields)

					views = append(views, self.newFormField(modelValue, metaData, disable, valueErrors))

				} else if validator, ok := data.(model.Validator); ok {

					generalErrors = append(generalErrors, validator.Validate(metaData)...)

				}
			})
		}

		hasErrors := numValueErrors > 0 || len(generalErrors) > 0

		if isPOST {
			var message string
			if hasErrors {
				message = "Form not saved because of invalid values! "
				for _, err := range generalErrors {
					message = message + err.WrappedError.Error() + ". "
				}
			} else {
				// Try to save the new form field values
				err = self.OnSubmit(self, formModel, context)
				if err == nil {
					message = self.SuccessMessage
				} else {
					message = err.Error()
					hasErrors = true
				}

				// Redirect if saved without errors and redirect URL is set
				if !hasErrors && self.Redirect != nil {
					return Redirect(self.Redirect.URL(context))
				}
			}

			if hasErrors {
				views[0] = newFormMessage(self.GetErrorMessageClass(), message)
				if len(views)-1 > Config.NumFieldRepeatFormMessage {
					views = append(views, newFormMessage(self.GetErrorMessageClass(), message))
				}
			} else {
				views[0] = newFormMessage(self.GetSuccessMessageClass(), message)
				if len(views)-1 > Config.NumFieldRepeatFormMessage {
					views = append(views, newFormMessage(self.GetSuccessMessageClass(), message))
				}
			}

		}

		// Add submit button and form ID:
		buttonText := self.ButtonText
		if buttonText == "" {
			buttonText = "Save"
		}
		formId := &HiddenInput{Name: "form_id", Value: self.FormID}
		submitButton := &Button{Submit: true, Value: buttonText, Class: self.ButtonClass}
		views = append(views, formId, submitButton)

		// Replace old form fields with the new ones:
		if self.Content != nil {
			self.Content.OnRemove()
		}
		self.Content = views
		views.Init(views)
	}

	// Render HTML form element
	method := self.Method
	if method == "" {
		method = "POST"
	}
	action := self.Action
	if action == "" {
		action = "."
	}

	writer.OpenTag("form").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("method", method)
	writer.Attrib("action", action)
	err = RenderChildViewsHTML(self, context, writer)
	writer.ExtraCloseTag() // form
	return err
}

func (self *Form) GetErrorMessageClass() string {
	if self.ErrorMessageClass == "" {
		return Config.FormErrorMessageClass
	}
	return self.ErrorMessageClass
}

func (self *Form) GetSuccessMessageClass() string {
	if self.SuccessMessageClass == "" {
		return Config.FormSuccessMessageClass
	}
	return self.SuccessMessageClass
}

func newFormMessage(class, message string) View {
	return &Div{
		Class:   class,
		Content: HTML(message),
	}
}

func getLabel(metaData model.MetaData) string {
	var buf bytes.Buffer
	for i := range metaData {
		if i > 0 {
			buf.WriteString(" ")
		}
		meta := &metaData[i]
		label, ok := meta.Attrib("label")
		if !ok {
			label = strings.Replace(meta.Name, "_", " ", -1)
		}
		buf.WriteString(label)
	}
	return buf.String()
}

func getClass(metaData model.MetaData) string {
	class, _ := metaData.TopField().Attrib("class")
	return class
}

func (self *Form) newVerticalFormField(modelValue model.Value, metaData model.MetaData, errors []*model.ValidationError, editorView View, extraLabels ...View) View {
	views := make(Views, 0, 2+len(errors)*2+1)
	views = append(views, &Label{Class: "vertical", Content: Escape(getLabel(metaData)), For: editorView})
	views = append(views, extraLabels...)
	for _, error := range errors {
		views = append(
			views,
			&Span{
				Class:   self.GetErrorMessageClass(),
				Content: Escape(error.WrappedError.Error()),
			},
			Br(),
		)
	}
	views = append(views, editorView)

	return &Paragraph{Content: views}
}

func (self *Form) newFormField(modelValue model.Value, metaData model.MetaData, disable bool, errors []*model.ValidationError) View {
	if len(metaData) == 0 {
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
		return self.newVerticalFormField(modelValue, metaData, errors, sel)

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
		return self.newVerticalFormField(modelValue, metaData, errors, sel)

	case *model.Country:
		// todo
		value := modelValue.(fmt.Stringer).String()
		if value == "" {
			value = "[empty]"
		}
		return self.newVerticalFormField(modelValue, metaData, errors, Escape(value))

	case *model.Date:
		date := s
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     date.Get(),
			Size:     len(model.DateFormat),
			Disabled: disable,
		}
		return self.newVerticalFormField(modelValue, metaData, errors, textField, HTML("(Format: "+model.DateFormat+")<br/>"))

	case *model.DateTime:
		dateTime := s
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     dateTime.Get(),
			Size:     len(model.DateTimeFormat),
			Disabled: disable,
		}
		return self.newVerticalFormField(modelValue, metaData, errors, textField, HTML("(Format: "+model.DateTimeFormat+")<br/>"))

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
		return self.newVerticalFormField(modelValue, metaData, errors, textField)

	case *model.Float, *model.Int:
		value := modelValue.(fmt.Stringer).String()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Disabled: disable,
		}
		return self.newVerticalFormField(modelValue, metaData, errors, textField)

	case *model.Language:
		// todo
		value := modelValue.(fmt.Stringer).String()
		if value == "" {
			value = "[empty]"
		}
		return self.newVerticalFormField(modelValue, metaData, errors, HTML(value))

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
		return self.newVerticalFormField(modelValue, metaData, errors, textField)

	case *model.Phone:
		value := s.Get()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Size:     20,
			Disabled: disable,
		}
		return self.newVerticalFormField(modelValue, metaData, errors, textField)

	case model.Reference:
		if !self.ShowRefIDs {
			return nil
		}
		value := modelValue.(fmt.Stringer).String()
		if value == "" {
			value = "[empty]"
		}
		return self.newVerticalFormField(modelValue, metaData, errors, HTML(value))

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
		return self.newVerticalFormField(modelValue, metaData, errors, textField)

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
		return self.newVerticalFormField(modelValue, metaData, errors, textArea)

	case *model.Url:
		value := s.Get()
		textField := &TextField{
			Class:    getClass(metaData),
			Name:     metaData.Selector(),
			Text:     value,
			Size:     80,
			Disabled: disable,
		}
		return self.newVerticalFormField(modelValue, metaData, errors, textField)

	case *model.GeoLocation:
		value := s.String()
		return self.newVerticalFormField(modelValue, metaData, errors, HTML(value))
	}

	panic(fmt.Sprintf("Unsupported model.Value type %T", modelValue))
}
