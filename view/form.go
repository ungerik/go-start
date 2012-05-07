package view

import (
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

type FormLayout interface {
	NewField(form *Form, modelValue model.Value, metaData *model.MetaData, disable bool, errors []*model.ValidationError) View
}

///////////////////////////////////////////////////////////////////////////////
// Form

type Form struct {
	ViewBaseWithId
	Class               string
	Action              string // Default is "." plus any URL params
	Content             View
	Method              string
	FormID              string
	CSRFProtector       CSRFProtector
	Layout              FormLayout // Config.DefaultFormLayout will be used if nil
	GetModel            GetFormModelFunc
	ModelMaxDepth       int      // if zero, no depth limit
	HideFields          []string // Use point notation for nested fields
	DisableFields       []string // Use point notation for nested fields
	RequireFields       []string // Also available as static struct field tag. Use point notation for nested fields
	OnSubmit            OnSubmitFormFunc
	ErrorMessageClass   string // If empty, Config.FormErrorMessageClass will be used
	SuccessMessageClass string // If empty, Config.FormSuccessMessageClass will be used
	SuccessMessage      string
	ButtonText          string
	ButtonClass         string
	Redirect            URL // 302 redirect after successful Save()
	ShowRefIDs          bool
}

// GetLayout returns self.Layout if not nil,
// else Config.DefaultFormLayout will be used.
func (self *Form) GetLayout() FormLayout {
	if self.Layout == nil {
		return Config.DefaultFormLayout
	}
	return self.Layout
}

func (self *Form) GetCSRFProtector() CSRFProtector {
	if self.CSRFProtector == nil {
		return Config.DefaultCSRFProtector
	}
	return self.CSRFProtector
}

func (self *Form) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Form) IsFieldRequired(metaData *model.MetaData) bool {
	if metaData.BoolAttrib("required") {
		return true
	}
	selector := metaData.Selector()
	arraySelector := metaData.ArrayWildcardSelector()
	return utils.StringIn(selector, self.RequireFields) || utils.StringIn(arraySelector, self.RequireFields)
}

func (self *Form) isFieldRequiredSelectors(metaData *model.MetaData, selector, arraySelector string) bool {
	if metaData.BoolAttrib("required") {
		return true
	}
	return utils.StringIn(selector, self.RequireFields) || utils.StringIn(arraySelector, self.RequireFields)
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
			model.WalkStructure(formModel, self.ModelMaxDepth, func(data interface{}, metaData *model.MetaData) {
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
								if len(valueErrors) == 0 {
									if modelValue.IsEmpty() && self.isFieldRequiredSelectors(metaData, selector, arraySelector) {
										valueErrors = []*model.ValidationError{model.NewRequiredValidationError(metaData)}
									}
								}
							} else {
								valueErrors = model.NewValidationErrors(err, metaData)
							}
							numValueErrors += len(valueErrors)
						}
					}

					disable := utils.StringIn(selector, self.DisableFields)

					views = append(views, self.GetLayout().NewField(self, modelValue, metaData, disable, valueErrors))

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
		if i := strings.Index(context.Request.RequestURI, "?"); i != -1 {
			action += context.Request.RequestURI[i:]
		}
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

func getLabel(metaData *model.MetaData) string {
	names := make([]string, metaData.Depth)
	for i, m := metaData.Depth-1, metaData; i >= 0; i-- {
		label, ok := m.Attrib("label")
		if !ok {
			label = strings.Replace(m.Name, "_", " ", -1)
		}
		names[i] = label
		m = m.Parent
	}
	return strings.Join(names, " ")
}

func getClass(metaData *model.MetaData) string {
	class, _ := metaData.Attrib("class")
	return class
}
