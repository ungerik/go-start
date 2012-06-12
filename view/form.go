package view

import (
	"bytes"
	"github.com/ungerik/go-start/debug"
	// "github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
	"strings"
	"strconv"
)

const FormIDName = "form_id"

type GetFormModelFunc func(form *Form, context *Context) (model interface{}, err error)
type OnSubmitFormFunc func(form *Form, formModel interface{}, context *Context) (message string, redirect URL, err error)

func FormModel(model interface{}) GetFormModelFunc {
	return func(form *Form, context *Context) (interface{}, error) {
		return model, nil
	}
}

func SaveModel(form *Form, formModel interface{}, context *Context) (string, URL, error) {
	return "", nil, formModel.(mongo.Document).Save()
}

/*
FormLayout is responsible for creating and structuring all dynamic content
of the form including the submit button.
It uses Form.GetFieldFactory() to create the field views.
*/
type FormLayout interface {
	BeginFormContent(form *Form, context *Context, formFields Views) Views
	// SubmitSuccess will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// not returned an error.
	SubmitSuccess(message string, form *Form, context *Context, formFields Views) Views
	// SubmitError will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// returned an error.
	SubmitError(message string, form *Form, context *Context, formFields Views) Views
	EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, context *Context, formFields Views) Views

	BeginStruct(strct *model.MetaData, form *Form, context *Context, formFields Views) Views
	StructField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views
	EndStruct(strct *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views

	BeginArray(array *model.MetaData, form *Form, context *Context, formFields Views) Views
	ArrayField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views
	EndArray(array *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views

	BeginSlice(slice *model.MetaData, form *Form, context *Context, formFields Views) Views
	SliceField(field *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views
	EndSlice(slice *model.MetaData, validationErr error, form *Form, context *Context, formFields Views) Views
}

type FormFieldFactory interface {
	CanCreateInput(metaData *model.MetaData, form *Form) bool
	NewInput(withLabel bool, metaData *model.MetaData, form *Form) View
	NewHiddenInput(metaData *model.MetaData, form *Form) View
	NewLabel(forView View, metaData *model.MetaData, form *Form) View
	NewFieldDescrtiption(description string, form *Form) View
	NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View
	NewGeneralErrorMessage(message string, form *Form) View
	NewSuccessMessage(message string, form *Form) View
	NewSubmitButton(text, confirmationMessage string, form *Form) View
	NewAddSliceItemButton(form *Form) View
	NewRemoveSliceItemButton(form *Form) View
}

///////////////////////////////////////////////////////////////////////////////
// Form

type Form struct {
	ViewBaseWithId
	Class  string
	Action string // Default is "." plus any URL params
	Method string
	// FormID must be unique on a page to identify the form for the case
	// that there are multiple forms on a single page.
	// Can be empty if there is only one form on a page, but it is good
	// practice to always use form ids.
	FormID        string
	CSRFProtector CSRFProtector
	Layout        FormLayout       // Config.DefaultFormLayout will be used if nil
	FieldFactory  FormFieldFactory // Config.DefaultFormFieldFactory will be used if nil
	GetModel      GetFormModelFunc

	// OnSubmit is called after the form was submitted and did not produce any
	// validation errors.
	// If message is non empty, it will be displayed instead of
	// Form.SuccessMessage or err.
	// If err is not nil, err.Error() will be displayed and not processed
	// any further (good idea?)
	// If redirect result is not nil, it will be used instead of Form.Redirect.
	// message or err will only be visible, if there is no redirect.
	OnSubmit OnSubmitFormFunc

	ModelMaxDepth         int      // if zero, no depth limit
	ExcludedFields        []string // Use point notation for nested fields. In case of arrays/slices use wildcards
	HiddenFields          []string // Use point notation for nested fields. In case of arrays/slices use wildcards
	DisabledFields        []string // Use point notation for nested fields
	RequiredFields        []string // Also available as static struct field tag. Use point notation for nested fields
	Labels                map[string]string
	FieldDescriptions     map[string]string
	ErrorMessageClass     string // If empty, Config.Form.DefaultErrorMessageClass will be used
	SuccessMessageClass   string // If empty, Config.Form.DefaultSuccessMessageClass will be used
	FieldDescriptionClass string
	RequiredMarker        View // If nil, Config.Form.DefaultRequiredMarker will be used
	SuccessMessage        string
	SubmitButtonText      string
	SubmitButtonClass     string
	SubmitButtonConfirm   string // Will add a confirmation dialog for onclick
	Redirect              URL    // 302 redirect after successful Save()
	ShowRefIDs            bool
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

func (self *Form) IsFieldRequired(metaData *model.MetaData) bool {
	if val, ok := metaData.ModelValue(); ok && val.Required(metaData) {
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

func (self *Form) IsPost(context *Context) bool {
	if context.Request.Method != "POST" {
		return false
	}
	formID, ok := context.Params[FormIDName]
	return ok && formID == self.FormID
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

// DirectFieldLabel returns a label for a form field generated from metaData.
// It creates the label only from the name or label tag of metaData,
// not including its parents.
func (self *Form) DirectFieldLabel(metaData *model.MetaData) string {
	if label, ok := metaData.Attrib("label"); ok {
		return label
	}
	return strings.Replace(metaData.NameOrIndex(), "_", " ", -1)
}

// FieldLabel returns a label for a form field generated from metaData.
// It creates the label from the names or label tags of metaData and
// all its parents, starting with the root parent, concanated with a space
// character.
func (self *Form) FieldLabel(metaData *model.MetaData) string {
	selector := metaData.Selector()
	if label, ok := self.Labels[selector]; ok {
		return label
	}
	wildcardSelector := metaData.WildcardSelector()
	if label, ok := self.Labels[wildcardSelector]; ok {
		return label
	}
	var buf bytes.Buffer
	for _, m := range metaData.Path()[1:] {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(self.DirectFieldLabel(m))
	}
	return buf.String()
}

func (self *Form) FieldInputClass(metaData *model.MetaData) string {
	class, _ := metaData.Attrib("class")
	return class
}

func (self *Form) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.OnSubmit == nil {
		panic("view.Form.OnSubmit must not be nil")
	}
	layout := self.GetLayout()
	if layout == nil {
		panic("view.Form.GetLayout() returned nil")
	}
	fieldFactory := self.GetFieldFactory()
	if fieldFactory == nil {
		panic("view.Form.GetFieldFactory() returned nil")
	}

	isPost := self.IsPost(context)
	var formFields Views
	var formModel interface{}

	if self.GetModel == nil {
		formId := &HiddenInput{Name: FormIDName, Value: self.FormID}
		submitButton := self.GetFieldFactory().NewSubmitButton(self.GetSubmitButtonText(), self.SubmitButtonConfirm, self)
		formFields = Views{formId, submitButton}
		if isPost {
			formFields, err = self.submit(nil, context, formFields)
			if err != nil {
				return err
			}
		}
	} else {
		formModel, err = self.GetModel(self, context)
		if err != nil {
			return err
		}
		visitor := &formLayoutWrappingStructVisitor{
			form:       self,
			formLayout: layout,
			formModel:  formModel,
			context:    context,
			isPost:     isPost,
		}
		err = model.Visit(formModel, visitor)
		if err != nil {
			return err
		}
		formFields = visitor.formFields
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
	if len(formFields) > 0 {
		formFields.Init(formFields)
		err = formFields.Render(context, writer)
		if err != nil {
			return err
		}
	}
	writer.ForceCloseTag() // form
	return nil
}

func (self *Form) submit(formModel interface{}, context *Context, formFields Views) (Views, error) {
	message, redirect, err := self.OnSubmit(self, formModel, context)
	if err == nil {
		if redirect == nil {
			redirect = self.Redirect
		}
		if redirect != nil {
			return nil, Redirect(redirect.URL(context))
		}
		if message == "" {
			message = self.SuccessMessage
		}
		if message != "" {
			formFields = self.GetLayout().SubmitSuccess(message, self, context, formFields)
		}
	} else {
		if message == "" {
			message = err.Error()
		}
		formFields = self.GetLayout().SubmitError(message, self, context, formFields)
	}
	return formFields, nil
}

///////////////////////////////////////////////////////////////////////////////
// formLayoutWrappingStructVisitor

type formLayoutWrappingStructVisitor struct {
	// Input 
	form       *Form
	formLayout FormLayout
	formModel  interface{}
	context    *Context
	isPost     bool

	// Output
	formFields              Views
	fieldValidationErrors   []error
	generalValidationErrors []error
}

func (self *formLayoutWrappingStructVisitor) postValue(metaData *model.MetaData) string {
	postValue, _ := self.context.Params[metaData.Selector()]
	return postValue
}

func (self *formLayoutWrappingStructVisitor) setFieldValue(field *model.MetaData) (err error) {
	switch s := field.Value.Addr().Interface().(type) {
	case *model.Bool:
		s.Set(self.postValue(field) != "")
		// model.Bool doesn't have validation
	case model.Reference:
		// do nothing
	case model.Value:
		err = s.SetString(self.postValue(field))
		if err == nil {
			err = s.Validate(field)
		}
		if err != nil {
			self.fieldValidationErrors = append(self.fieldValidationErrors, err)
		}
	}
	return err
}

func (self *formLayoutWrappingStructVisitor) validate(data *model.MetaData) (err error) {
	if validator, ok := data.ModelValidator(); ok {
		err = validator.Validate(data)
		if err != nil {
			self.generalValidationErrors = append(self.generalValidationErrors, err)
		}
	}
	return err
}

func (self *formLayoutWrappingStructVisitor) endForm(data *model.MetaData) (err error) {
	if self.isPost && len(self.fieldValidationErrors) == 0 && len(self.generalValidationErrors) == 0 {
		self.formFields, err = self.form.submit(self.formModel, self.context, self.formFields)
		if err != nil {
			return err
		}
	}
	self.formFields = self.formLayout.EndFormContent(self.fieldValidationErrors, self.generalValidationErrors, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) BeginStruct(strct *model.MetaData) error {
	if strct.Parent == nil {
		self.formFields = self.formLayout.BeginFormContent(self.form, self.context, self.formFields)
	}
	self.formFields = self.formLayout.BeginStruct(strct, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) StructField(field *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.setFieldValue(field)
	}
	self.formFields = self.formLayout.StructField(field, validationErr, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) EndStruct(strct *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.validate(strct)
	}
	self.formFields = self.formLayout.EndStruct(strct, validationErr, self.form, self.context, self.formFields)
	if strct.Parent == nil {
		return self.endForm(strct)
	}
	return nil
}

func (self *formLayoutWrappingStructVisitor) ModifySlice(slice *model.MetaData) error {
	debug.Nop()
	var length int
	if self.isPost {
		length, _ = strconv.Atoi(self.context.Params[slice.Selector()+".length"])
	} else {
		length = slice.Value.Len() + 1
	}
	if length != slice.Value.Len() {
		slice.Value = utils.SetSliceLengh(slice.Value, length)
		mongo.InitRefs(self.formModel)
	}
	return nil
}

func (self *formLayoutWrappingStructVisitor) BeginSlice(slice *model.MetaData) error {
	if slice.Parent == nil {
		self.formFields = self.formLayout.BeginFormContent(self.form, self.context, self.formFields)
	}
	self.formFields = self.formLayout.BeginSlice(slice, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) SliceField(field *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.setFieldValue(field)
	}
	self.formFields = self.formLayout.SliceField(field, validationErr, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) EndSlice(slice *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.validate(slice)
	}
	self.formFields = self.formLayout.EndSlice(slice, validationErr, self.form, self.context, self.formFields)
	if slice.Parent == nil {
		return self.endForm(slice)
	}
	return nil
}

func (self *formLayoutWrappingStructVisitor) BeginArray(array *model.MetaData) error {
	if array.Parent == nil {
		self.formFields = self.formLayout.BeginFormContent(self.form, self.context, self.formFields)
	}
	self.formFields = self.formLayout.BeginArray(array, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) ArrayField(field *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.setFieldValue(field)
	}
	self.formFields = self.formLayout.ArrayField(field, validationErr, self.form, self.context, self.formFields)
	return nil
}

func (self *formLayoutWrappingStructVisitor) EndArray(array *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.validate(array)
	}
	self.formFields = self.formLayout.EndArray(array, validationErr, self.form, self.context, self.formFields)
	if array.Parent == nil {
		return self.endForm(array)
	}
	return nil
}
