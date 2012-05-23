package view

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
	// "reflect"
	"strings"
	//	"github.com/ungerik/go-start/debug"
)

type GetFormModelFunc func(form *Form, context *Context) (model interface{}, err error)

func FormModel(model interface{}) GetFormModelFunc {
	return func(form *Form, context *Context) (interface{}, error) {
		return model, nil
	}
}

/*
FormLayout is responsible for creating and structuring all dynamic content
of the form including the submit button.
It uses Form.GetFieldFactory() to create the field views.
*/
type FormLayout interface {
	BeginFormContent(form *Form, formFields Views) Views
	EndFormContent(fieldValidationErrs, generalValidationErrs []*model.ValidationError, form *Form, formFields Views) Views

	BeginStruct(strct *model.MetaData, form *Form, formFields Views) Views
	StructField(field *model.MetaData, validationErrs []*model.ValidationError, form *Form, formFields Views) Views
	EndStruct(strct *model.MetaData, validationErrs []*model.ValidationError, form *Form, formFields Views) Views

	BeginArray(array *model.MetaData, form *Form, formFields Views) Views
	ArrayField(field *model.MetaData, validationErrs []*model.ValidationError, form *Form, formFields Views) Views
	EndArray(array *model.MetaData, validationErrs []*model.ValidationError, form *Form, formFields Views) Views

	BeginSlice(slice *model.MetaData, form *Form, formFields Views) Views
	SliceField(field *model.MetaData, validationErrs []*model.ValidationError, form *Form, formFields Views) Views
	EndSlice(slice *model.MetaData, validationErrs []*model.ValidationError, form *Form, formFields Views) Views

	FieldNeedsLabel(field *model.MetaData, form *Form) bool
}

type FormFieldFactory interface {
	NewInput(metaData *model.MetaData, form *Form) View
	NewHiddenInput(metaData *model.MetaData, form *Form) View
	NewLabel(forView View, metaData *model.MetaData, form *Form) View
	NewFieldDescrtiption(description string, form *Form) View
	NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View
	NewGeneralErrorMessage(message string, form *Form) View
	NewSuccessMessage(message string, form *Form) View
	NewSubmitButton(text string, form *Form) View
	NewAddSliceItemButton(form *Form) View
	NewRemoveSliceItemButton(form *Form) View
}

///////////////////////////////////////////////////////////////////////////////
// Form

type Form struct {
	ViewBaseWithId
	Class         string
	Action        string // Default is "." plus any URL params
	Method        string
	FormID        string
	CSRFProtector CSRFProtector
	Layout        FormLayout       // Config.DefaultFormLayout will be used if nil
	FieldFactory  FormFieldFactory // Config.DefaultFormFieldFactory will be used if nil
	// Static content rendered before the dynamic form fields
	// that are generated via GetModel()
	StaticContent View
	GetModel      GetFormModelFunc
	// If redirect result is non nil, it will be used instead of Form.Redirect
	OnSubmit              func(form *Form, formModel interface{}, context *Context) (redirect URL, err error)
	ModelMaxDepth         int      // if zero, no depth limit
	ExcludedFields        []string // Use point notation for nested fields
	HiddenFields          []string // Use point notation for nested fields
	DisabledFields        []string // Use point notation for nested fields
	RequiredFields        []string // Also available as static struct field tag. Use point notation for nested fields
	FieldDescriptions     map[string]string
	ErrorMessageClass     string // If empty, Config.Form.DefaultErrorMessageClass will be used
	SuccessMessageClass   string // If empty, Config.Form.DefaultSuccessMessageClass will be used
	FieldDescriptionClass string
	RequiredMarker        View // If nil, Config.Form.DefaultRequiredMarker will be used
	SuccessMessage        string
	SubmitButtonText      string
	SubmitButtonClass     string
	Redirect              URL // 302 redirect after successful Save()
	ShowRefIDs            bool

	UseNewFormMode bool
}

// GetLayout returns self.Layout if not nil,
// else Config.Form.DefaultLayout will be returned.
func (self *Form) GetLayout() FormLayout {
	if self.Layout == nil {
		return Config.Form.DefaultLayout
	}
	return self.Layout
}

// GetFieldFactory returns self.FieldFactory if not nil,
// else Config.Form.DefaultFieldFactory will be returned.
func (self *Form) GetFieldFactory() FormFieldFactory {
	if self.FieldFactory == nil {
		return Config.Form.DefaultFieldFactory
	}
	return self.FieldFactory
}

// GetCSRFProtector returns self.CSRFProtector if not nil,
// else Config.Form.DefaultCSRFProtector will be returned.
func (self *Form) GetCSRFProtector() CSRFProtector {
	if self.CSRFProtector == nil {
		return Config.Form.DefaultCSRFProtector
	}
	return self.CSRFProtector
}

func (self *Form) GetErrorMessageClass() string {
	if self.ErrorMessageClass == "" {
		return Config.Form.DefaultErrorMessageClass
	}
	return self.ErrorMessageClass
}

func (self *Form) GetSuccessMessageClass() string {
	if self.SuccessMessageClass == "" {
		return Config.Form.DefaultSuccessMessageClass
	}
	return self.SuccessMessageClass
}

func (self *Form) GetSubmitButtonClass() string {
	if self.SubmitButtonClass == "" {
		return Config.Form.DefaultSubmitButtonClass
	}
	return self.SubmitButtonClass
}

func (self *Form) GetFieldDescriptionClass() string {
	if self.FieldDescriptionClass == "" {
		return Config.Form.DefaultFieldDescriptionClass
	}
	return self.FieldDescriptionClass
}

func (self *Form) GetRequiredMarker() View {
	if self.RequiredMarker == nil {
		return Config.Form.DefaultRequiredMarker
	}
	return self.RequiredMarker
}

// Returns self.SubmitButtonText if not empty,
// else Config.Form.DefaultSubmitButtonText
func (self *Form) GetSubmitButtonText() string {
	if self.SubmitButtonText == "" {
		return Config.Form.DefaultSubmitButtonText
	}
	return self.SubmitButtonText
}

func (self *Form) IterateChildren(callback IterateChildrenCallback) {
	if self.StaticContent != nil {
		callback(self, self.StaticContent)
	}
}

func (self *Form) IsFieldRequired(metaData *model.MetaData) bool {
	if metaData.BoolAttrib("required") {
		return true
	}
	selector := metaData.Selector()
	wildcardSelector := metaData.WildcardSelector()
	return utils.StringIn(selector, self.RequiredFields) || utils.StringIn(wildcardSelector, self.RequiredFields)
}

func (self *Form) isFieldRequiredSelectors(metaData *model.MetaData, selector, wildcardSelector string) bool {
	if metaData.BoolAttrib("required") {
		return true
	}
	return utils.StringIn(selector, self.RequiredFields) || utils.StringIn(wildcardSelector, self.RequiredFields)
}

func (self *Form) IsFieldDisabled(metaData *model.MetaData) bool {
	if metaData.BoolAttrib("disabled") {
		return true
	}
	selector := metaData.Selector()
	wildcardSelector := metaData.WildcardSelector()
	return utils.StringIn(selector, self.DisabledFields) || utils.StringIn(wildcardSelector, self.DisabledFields)
}

func (self *Form) IsFieldHidden(metaData *model.MetaData) bool {
	if metaData.BoolAttrib("hidden") {
		return true
	}
	selector := metaData.Selector()
	wildcardSelector := metaData.WildcardSelector()
	return utils.StringIn(selector, self.HiddenFields) || utils.StringIn(wildcardSelector, self.HiddenFields)
}

func (self *Form) IsFieldExcluded(metaData *model.MetaData) bool {
	selector := metaData.Selector()
	wildcardSelector := metaData.WildcardSelector()
	return utils.StringIn(selector, self.ExcludedFields) || utils.StringIn(wildcardSelector, self.ExcludedFields)
}

func (self *Form) GetFieldDescription(metaData *model.MetaData) string {
	selector := metaData.Selector()
	if desc, ok := self.FieldDescriptions[selector]; ok {
		return desc
	}
	wildcardSelector := metaData.WildcardSelector()
	if desc, ok := self.FieldDescriptions[wildcardSelector]; ok {
		return desc
	}
	if desc, ok := metaData.Attrib("description"); ok {
		return desc
	}
	return ""
}

func (self *Form) Render(context *Context, writer *utils.XMLWriter) (err error) {
	layout := self.GetLayout()
	if layout == nil {
		panic("view.Form.GetLayout() returned nil")
	}
	fieldFactory := self.GetFieldFactory()
	if fieldFactory == nil {
		panic("view.Form.GetFieldFactory() returned nil")
	}

	var dynamicFields Views

	if self.OnSubmit != nil {
		var formModel interface{}

		if self.GetModel != nil {
			formModel, err = self.GetModel(self, context)
			if err != nil {
				return err
			}

			visitor := &formLayoutWrappingStructVisitor{
				form:       self,
				formLayout: self.GetLayout(),
				context:    context,
			}

			model.Visit(formModel, visitor)

			// hasErrors := numValueErrors > 0 || len(generalErrors) > 0

			// if isPOST {
			// 	var message string
			// 	if hasErrors {
			// 		message = "Form not saved because of invalid values! "
			// 		for _, err := range generalErrors {
			// 			message = message + err.WrappedError.Error() + ". "
			// 		}
			// 	} else {
			// 		// Try to save the new form field values
			// 		redirect, err := self.OnSubmit(self, formModel, context)
			// 		if err == nil {
			// 			message = self.SuccessMessage
			// 			if redirect == nil {
			// 				redirect = self.Redirect
			// 			}
			// 		} else {
			// 			message = err.Error()
			// 			hasErrors = true
			// 		}

			// 		// Redirect if saved without errors and redirect URL is set
			// 		if !hasErrors && redirect != nil {
			// 			return Redirect(redirect.URL(context))
			// 		}
			// 	}

			// 	if hasErrors {
			// 		dynamicFields[0] = DIV(self.GetErrorMessageClass(), Escape(message))
			// 		if len(dynamicFields)-1 > Config.Form.NumFieldRepeatMessage {
			// 			dynamicFields = append(dynamicFields, DIV(self.GetErrorMessageClass(), Escape(message)))
			// 		}
			// 	} else {
			// 		dynamicFields[0] = DIV(self.GetSuccessMessageClass(), Escape(message))
			// 		if len(dynamicFields)-1 > Config.Form.NumFieldRepeatMessage {
			// 			dynamicFields = append(dynamicFields, DIV(self.GetSuccessMessageClass(), Escape(message)))
			// 		}
			// 	}

			// }

			// // Add submit button and form ID:
			// buttonText := self.SubmitButtonText
			// if buttonText == "" {
			// 	buttonText = "Save"
			// }
			// formId := &HiddenInput{Name: "form_id", Value: self.FormID}
			// submitButton := self.GetFieldFactory().NewSubmitButton(self, buttonText)
			// dynamicFields = append(dynamicFields, formId, submitButton)
		}
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
	if err != nil {
		return err
	}
	if dynamicFields != nil {
		dynamicFields.Init(dynamicFields)
		err = dynamicFields.Render(context, writer)
		if err != nil {
			return err
		}
	}
	writer.ExtraCloseTag() // form
	return nil
}

// func (self *Form) Render(context *Context, writer *utils.XMLWriter) (err error) {
// 	if self.UseNewFormMode {
// 		return self.newRender(context, writer)
// 	}

// 	var dynamicFields Views

// 	if self.OnSubmit != nil {
// 		// Determine if it's a POST request for this form:
// 		isPOST := false
// 		if context.Request.Method == "POST" {
// 			// Every HTML form gets an ID to allow more than one form per page:
// 			id, ok := context.Params["form_id"]
// 			if ok && id == self.FormID {
// 				isPOST = true
// 			}
// 		}

// 		// Create views for form fields:

// 		// Set a view before and after the form fields
// 		// if there is an error or success message
// 		// (won't be rendered if nil)
// 		// Also add a hidden form field with the form id
// 		// and a submit button
// 		dynamicFields = make(Views, 1, 32)

// 		numValueErrors := 0
// 		generalErrors := []*model.ValidationError{}

// 		var formModel interface{}

// 		if self.GetModel != nil {
// 			formModel, err = self.GetModel(self, context)
// 			if err != nil {
// 				return err
// 			}
// 			model.WalkStructure(formModel, self.ModelMaxDepth,
// 				func(data *model.MetaData) {
// 					if modelValue, ok := data.Value.Addr().Interface().(model.Value); ok {
// 						selector := data.Selector()
// 						arrayWildcardSelector := data.WildcardSelector()

// 						if utils.StringIn(selector, self.HideFields) || utils.StringIn(arrayWildcardSelector, self.HideFields) {
// 							return
// 						}

// 						var valueErrors []*model.ValidationError
// 						if isPOST {
// 							formValue, ok := context.Params[selector]
// 							if b, isBool := modelValue.(*model.Bool); isBool {
// 								b.Set(formValue != "")
// 							} else if ok {
// 								err = modelValue.SetString(formValue)
// 								if err == nil {
// 									valueErrors = modelValue.Validate(data)
// 									if len(valueErrors) == 0 {
// 										if modelValue.IsEmpty() && self.isFieldRequiredSelectors(data, selector, arrayWildcardSelector) {
// 											valueErrors = []*model.ValidationError{model.NewRequiredValidationError(data)}
// 										}
// 									}
// 								} else {
// 									valueErrors = model.NewValidationErrors(err, data)
// 								}
// 								numValueErrors += len(valueErrors)
// 							}
// 						}

// 						dynamicFields = append(dynamicFields, self.GetLayout().NewField_old(self, modelValue, data, valueErrors))

// 					} else if validator, ok := data.Value.Interface().(model.Validator); ok {

// 						generalErrors = append(generalErrors, validator.Validate(data)...)

// 					}
// 				},
// 			)
// 		}

// 		hasErrors := numValueErrors > 0 || len(generalErrors) > 0

// 		if isPOST {
// 			var message string
// 			if hasErrors {
// 				message = "Form not saved because of invalid values! "
// 				for _, err := range generalErrors {
// 					message = message + err.WrappedError.Error() + ". "
// 				}
// 			} else {
// 				// Try to save the new form field values
// 				redirect, err := self.OnSubmit(self, formModel, context)
// 				if err == nil {
// 					message = self.SuccessMessage
// 					if redirect == nil {
// 						redirect = self.Redirect
// 					}
// 				} else {
// 					message = err.Error()
// 					hasErrors = true
// 				}

// 				// Redirect if saved without errors and redirect URL is set
// 				if !hasErrors && redirect != nil {
// 					return Redirect(redirect.URL(context))
// 				}
// 			}

// 			if hasErrors {
// 				dynamicFields[0] = DIV(self.GetErrorMessageClass(), Escape(message))
// 				if len(dynamicFields)-1 > Config.Form.NumFieldRepeatMessage {
// 					dynamicFields = append(dynamicFields, DIV(self.GetErrorMessageClass(), Escape(message)))
// 				}
// 			} else {
// 				dynamicFields[0] = DIV(self.GetSuccessMessageClass(), Escape(message))
// 				if len(dynamicFields)-1 > Config.Form.NumFieldRepeatMessage {
// 					dynamicFields = append(dynamicFields, DIV(self.GetSuccessMessageClass(), Escape(message)))
// 				}
// 			}

// 		}

// 		// Add submit button and form ID:
// 		buttonText := self.SubmitButtonText
// 		if buttonText == "" {
// 			buttonText = "Save"
// 		}
// 		formId := &HiddenInput{Name: "form_id", Value: self.FormID}
// 		submitButton := self.GetFieldFactory().NewSubmitButton(self, buttonText)
// 		dynamicFields = append(dynamicFields, formId, submitButton)
// 	}

// 	// Render HTML form element
// 	method := self.Method
// 	if method == "" {
// 		method = "POST"
// 	}
// 	action := self.Action
// 	if action == "" {
// 		action = "."
// 		if i := strings.Index(context.Request.RequestURI, "?"); i != -1 {
// 			action += context.Request.RequestURI[i:]
// 		}
// 	}

// 	writer.OpenTag("form").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
// 	writer.Attrib("method", method)
// 	writer.Attrib("action", action)
// 	err = RenderChildViewsHTML(self, context, writer)
// 	if err != nil {
// 		return err
// 	}
// 	if dynamicFields != nil {
// 		dynamicFields.Init(dynamicFields)
// 		err = dynamicFields.Render(context, writer)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	writer.ExtraCloseTag() // form
// 	return nil
// }

func (self *Form) FieldLabel(metaData *model.MetaData) string {
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

func (self *Form) FieldInputClass(metaData *model.MetaData) string {
	class, _ := metaData.Attrib("class")
	return class
}

///////////////////////////////////////////////////////////////////////////////
// formLayoutWrappingStructVisitor

type formLayoutWrappingStructVisitor struct {
	form       *Form
	formLayout FormLayout
	context    *Context

	formFields              Views
	fieldValidationErrors   []*model.ValidationError
	generalValidationErrors []*model.ValidationError
}

func (self *formLayoutWrappingStructVisitor) setFieldValue(field *model.MetaData) (errs []*model.ValidationError) {
	if field.Kind != model.ValueKind || self.context.Request.Method != "POST" {
		return nil
	}
	postValue, _ := self.context.Params[field.Selector()]

	switch s := field.Value.Addr().Interface().(type) {
	case *model.Bool:
		s.Set(postValue != "")
		// model.Bool doesn't have validation
	case model.Value:
		err := s.SetString(postValue)
		if err == nil {
			errs = s.Validate(field)
		} else {
			errs = model.NewValidationErrors(err, field)
		}
	default:
		panic("Unknown form field type")
	}

	if len(errs) > 0 {
		self.generalValidationErrors = append(self.generalValidationErrors, errs...)
	}
	return errs
}

func (self *formLayoutWrappingStructVisitor) validate(data *model.MetaData) (errs []*model.ValidationError) {
	if validator, ok := data.Value.Addr().Interface().(model.Validator); ok {
		errs = validator.Validate(data)
		if len(errs) > 0 {
			self.generalValidationErrors = append(self.generalValidationErrors, errs...)
		}
	}
	return errs
}

func (self *formLayoutWrappingStructVisitor) BeginStruct(strct *model.MetaData) {
	self.formFields = self.formLayout.BeginStruct(strct, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) StructField(field *model.MetaData) {
	validationErrs := self.setFieldValue(field)
	self.formFields = self.formLayout.StructField(field, validationErrs, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) EndStruct(strct *model.MetaData) {
	validationErrs := self.validate(strct)
	self.formFields = self.formLayout.EndStruct(strct, validationErrs, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) BeginSlice(slice *model.MetaData) {
	self.formFields = self.formLayout.BeginSlice(slice, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) SliceField(field *model.MetaData) {
	validationErrs := self.setFieldValue(field)
	self.formFields = self.formLayout.SliceField(field, validationErrs, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) EndSlice(slice *model.MetaData) {
	validationErrs := self.validate(slice)
	self.formFields = self.formLayout.EndSlice(slice, validationErrs, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) BeginArray(array *model.MetaData) {
	self.formFields = self.formLayout.BeginArray(array, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) ArrayField(field *model.MetaData) {
	validationErrs := self.setFieldValue(field)
	self.formFields = self.formLayout.ArrayField(field, validationErrs, self.form, self.formFields)
}

func (self *formLayoutWrappingStructVisitor) EndArray(array *model.MetaData) {
	validationErrs := self.validate(array)
	self.formFields = self.formLayout.EndArray(array, validationErrs, self.form, self.formFields)
}
