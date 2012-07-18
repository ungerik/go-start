package view

import (
	"bytes"
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/utils"
	"io/ioutil"
	"strconv"
	"strings"
	// "mime/multipart"
)

const MultipartFormData = "multipart/form-data"
const FormIDName = "gostart_form_id"

type GetFormModelFunc func(form *Form, context *Context) (model interface{}, err error)
type OnSubmitFormFunc func(form *Form, formModel interface{}, context *Context) (message string, redirect URL, err error)

type FileUploadFormModel struct {
	File model.File
}

// GetFileUploadFormModel when set at Form.GetModel will create a
// file upload form. Form.Enctype will be set to "multipart/form-fata"
// *FileUploadFormModel will be passed as formModel at Form.OnSubmit
func GetFileUploadFormModel(form *Form, context *Context) (interface{}, error) {
	form.Enctype = MultipartFormData
	return &FileUploadFormModel{}, nil
}

func FormModel(model interface{}) GetFormModelFunc {
	debug.Nop()
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
	GetDefaultInputSize(metaData *model.MetaData) int

	BeginFormContent(form *Form, context *Context, formContent *Views) error
	// SubmitSuccess will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// not returned an error.
	SubmitSuccess(message string, form *Form, context *Context, formContent *Views) error
	// SubmitError will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// returned an error.
	SubmitError(message string, form *Form, context *Context, formContent *Views) error
	EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, context *Context, formContent *Views) error

	BeginStruct(strct *model.MetaData, form *Form, context *Context, formContent *Views) error
	StructField(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error
	EndStruct(strct *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error

	BeginArray(array *model.MetaData, form *Form, context *Context, formContent *Views) error
	ArrayField(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error
	EndArray(array *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error

	BeginSlice(slice *model.MetaData, form *Form, context *Context, formContent *Views) error
	SliceField(field *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error
	EndSlice(slice *model.MetaData, validationErr error, form *Form, context *Context, formContent *Views) error
}

type FormFieldFactory interface {
	CanCreateInput(metaData *model.MetaData, form *Form) bool
	NewInput(withLabel bool, metaData *model.MetaData, form *Form) (View, error)
	NewHiddenInput(metaData *model.MetaData, form *Form) (View, error)
	NewTableHeader(metaData *model.MetaData, form *Form) (View, error)
	NewFieldDescrtiption(description string, form *Form) View
	NewFieldErrorMessage(message string, metaData *model.MetaData, form *Form) View
	NewGeneralErrorMessage(message string, form *Form) View
	NewSuccessMessage(message string, form *Form) View
	NewSubmitButton(text, confirmationMessage string, form *Form) View
	NewAddButton(onclick string, form *Form) View
	NewRemoveButton(onclick string, form *Form) View
	NewUpButton(disabled bool, onclick string, form *Form) View
	NewDownButton(disabled bool, onclick string, form *Form) View
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
	Layout        FormLayout       // Config.Form.DefaultLayout will be used if nil
	FieldFactory  FormFieldFactory // Config.Form.DefaultFieldFactory will be used if nil
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
	InputSizes            map[string]int
	ErrorMessageClass     string // If empty, Config.Form.DefaultErrorMessageClass will be used
	SuccessMessageClass   string // If empty, Config.Form.DefaultSuccessMessageClass will be used
	FieldDescriptionClass string
	RequiredMarker        View // If nil, Config.Form.DefaultRequiredMarker will be used
	SuccessMessage        string
	SubmitButtonText      string
	SubmitButtonClass     string
	SubmitButtonConfirm   string // Will add a confirmation dialog for onclick
	Redirect              URL    // 302 redirect after successful OnSubmit()
	ShowRefIDs            bool
	Enctype               string
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

func (self *Form) GetInputSize(metaData *model.MetaData) int {
	selector := metaData.Selector()
	if size, ok := self.InputSizes[selector]; ok {
		return size
	}
	wildcardSelector := metaData.WildcardSelector()
	if size, ok := self.InputSizes[wildcardSelector]; ok {
		return size
	}
	if sizeStr, ok := metaData.Attrib("size"); ok {
		size, _ := strconv.Atoi(sizeStr)
		return size
	}
	layout := self.GetLayout()
	if layout != nil {
		return layout.GetDefaultInputSize(metaData)
	}
	return 0
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

func (self *Form) IsPost(context *Context) bool {
	return context.Request.Method == "POST" &&
		context.Request.FormValue(FormIDName) == self.FormID
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

	var formModel interface{}
	var hasErrors bool
	content := Views{&HiddenInput{Name: FormIDName, Value: self.FormID}}
	isPost := self.IsPost(context)

	if self.GetModel == nil {
		submitButton := self.GetFieldFactory().NewSubmitButton(self.GetSubmitButtonText(), self.SubmitButtonConfirm, self)
		content = append(content, submitButton)
	} else {
		formModel, err = self.GetModel(self, context)
		if err != nil {
			return err
		}
		if isPost {
			setPostValues := &setPostValuesStructVisitor{
				form:      self,
				formModel: formModel,
				context:   context,
			}
			err = model.Visit(formModel, setPostValues)
			if err != nil {
				return err
			}
		}
		validateAndFormLayout := &validateAndFormLayoutStructVisitor{
			form:        self,
			formLayout:  layout,
			formModel:   formModel,
			formContent: &content,
			context:     context,
			isPost:      isPost,
		}
		err = model.Visit(formModel, validateAndFormLayout)
		if err != nil {
			return err
		}
		hasErrors = len(validateAndFormLayout.fieldValidationErrors) > 0 || len(validateAndFormLayout.generalValidationErrors) > 0
	}

	if isPost && !hasErrors {
		message, redirect, err := self.OnSubmit(self, formModel, context)
		if err == nil {
			if redirect == nil {
				redirect = self.Redirect
			}
			if redirect != nil {
				return Redirect(redirect.URL(context.PathArgs...))
			}
			if message == "" {
				message = self.SuccessMessage
			}
			if message != "" {
				self.GetLayout().SubmitSuccess(message, self, context, &content)
			}
		} else {
			if message == "" {
				message = err.Error()
			}
			self.GetLayout().SubmitError(message, self, context, &content)
		}
	}

	// Render HTML form element
	method := strings.ToUpper(self.Method)
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
	// Hack: Value of hidden input FormIDName is not available when
	// enctype is multipart/form-data (bug?), so pass the form id as
	// URL parameter
	if self.Enctype == MultipartFormData && context.Request.Method != "POST" {
		action = utils.AddUrlParam(action, FormIDName, self.FormID)
	}

	writer.OpenTag("form").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("method", method)
	writer.Attrib("action", action)
	writer.AttribIfNotDefault("enctype", self.Enctype)
	if len(content) > 0 {
		content.Init(content)
		err = content.Render(context, writer)
		if err != nil {
			return err
		}
	}
	writer.ForceCloseTag() // form
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// setPostValuesStructVisitor

type setPostValuesStructVisitor struct {
	form      *Form
	formModel interface{}
	context   *Context
}

func (self *setPostValuesStructVisitor) trySetFieldValue(field *model.MetaData) error {
	if self.form.IsFieldDisabled(field) || self.form.IsFieldExcluded(field) {
		return nil
	}

	switch s := field.Value.Addr().Interface().(type) {
	case model.Reference:
		// we don't handle references with forms yet
		return nil

	case *model.Bool:
		s.Set(self.context.Request.FormValue(field.Selector()) != "")
		return nil

	case *model.Blob:
		file, _, err := self.context.Request.FormFile(field.Selector())
		if err != nil {
			return err
		}
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		s.Set(bytes)
		return nil

	case *model.File:
		file, header, err := self.context.Request.FormFile(field.Selector())
		if err != nil {
			return err
		}
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		s.Name = header.Filename
		s.Data = bytes
		return nil

	case model.Value:
		s.SetString(self.context.Request.FormValue(field.Selector()))
		return nil
	}
	return nil
}

func (self *setPostValuesStructVisitor) BeginStruct(strct *model.MetaData) error {
	return nil
}

func (self *setPostValuesStructVisitor) StructField(field *model.MetaData) error {
	return self.trySetFieldValue(field)
}

func (self *setPostValuesStructVisitor) EndStruct(strct *model.MetaData) error {
	return nil
}

func (self *setPostValuesStructVisitor) BeginSlice(slice *model.MetaData) error {
	if slice.Value.CanSet() && !self.form.IsFieldExcluded(slice) {
		if lengthStr := self.context.Request.FormValue(slice.Selector() + ".length"); lengthStr != "" {
			length, err := strconv.Atoi(lengthStr)
			if err != nil {
				panic(err.Error())
			}
			if length != slice.Value.Len() {
				slice.Value.Set(utils.SetSliceLengh(slice.Value, length))
				mongo.InitRefs(self.formModel)
			}
		}
	}
	return nil
}

func (self *setPostValuesStructVisitor) SliceField(field *model.MetaData) error {
	return self.trySetFieldValue(field)
}

func (self *setPostValuesStructVisitor) EndSlice(slice *model.MetaData) error {
	if slice.Value.CanSet() && !self.form.IsFieldExcluded(slice) {
		slice.Value.Set(utils.DeleteEmptySliceElementsVal(slice.Value))
	}
	return nil
}

func (self *setPostValuesStructVisitor) BeginArray(array *model.MetaData) error {
	return nil
}

func (self *setPostValuesStructVisitor) ArrayField(field *model.MetaData) error {
	return self.trySetFieldValue(field)
}

func (self *setPostValuesStructVisitor) EndArray(array *model.MetaData) error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// validateAndFormLayoutStructVisitor

// validateAndFormLayoutStructVisitor creates the views for the form layout
type validateAndFormLayoutStructVisitor struct {
	// Input 
	form        *Form
	formLayout  FormLayout
	formModel   interface{}
	formContent *Views
	context     *Context
	isPost      bool

	// Output
	formFields              Views
	fieldValidationErrors   []error
	generalValidationErrors []error
}

func (self *validateAndFormLayoutStructVisitor) validateField(field *model.MetaData) error {
	if !self.isPost {
		return nil
	}
	err := field.Validate()
	if err != nil {
		self.fieldValidationErrors = append(self.fieldValidationErrors, err)
	}
	return err
}

func (self *validateAndFormLayoutStructVisitor) validateGeneral(data *model.MetaData) error {
	if !self.isPost {
		return nil
	}
	err := data.Validate()
	if err != nil {
		self.generalValidationErrors = append(self.generalValidationErrors, err)
	}
	return err
}

func (self *validateAndFormLayoutStructVisitor) endForm(data *model.MetaData) (err error) {
	return self.formLayout.EndFormContent(self.fieldValidationErrors, self.generalValidationErrors, self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) BeginStruct(strct *model.MetaData) error {
	if strct.Parent == nil {
		err := self.formLayout.BeginFormContent(self.form, self.context, self.formContent)
		if err != nil {
			return err
		}
	}
	return self.formLayout.BeginStruct(strct, self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) StructField(field *model.MetaData) error {
	return self.formLayout.StructField(field, self.validateField(field), self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) EndStruct(strct *model.MetaData) error {
	err := self.formLayout.EndStruct(strct, self.validateGeneral(strct), self.form, self.context, self.formContent)
	if strct.Parent == nil && err == nil {
		return self.endForm(strct)
	}
	return err
}

func (self *validateAndFormLayoutStructVisitor) BeginSlice(slice *model.MetaData) error {
	if slice.Parent == nil {
		err := self.formLayout.BeginFormContent(self.form, self.context, self.formContent)
		if err != nil {
			return err
		}
	}
	// Add an empty slice field to generate one extra input row for slices
	if !self.isPost && slice.Value.CanSet() && !self.form.IsFieldExcluded(slice) {
		slice.Value.Set(utils.AppendEmptySliceField(slice.Value))
		mongo.InitRefs(self.formModel)
	}
	return self.formLayout.BeginSlice(slice, self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) SliceField(field *model.MetaData) error {
	return self.formLayout.SliceField(field, self.validateField(field), self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) EndSlice(slice *model.MetaData) error {
	err := self.formLayout.EndSlice(slice, self.validateGeneral(slice), self.form, self.context, self.formContent)
	if slice.Parent == nil && err == nil {
		return self.endForm(slice)
	}
	return err
}

func (self *validateAndFormLayoutStructVisitor) BeginArray(array *model.MetaData) error {
	if array.Parent == nil {
		err := self.formLayout.BeginFormContent(self.form, self.context, self.formContent)
		if err != nil {
			return err
		}
	}
	return self.formLayout.BeginArray(array, self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) ArrayField(field *model.MetaData) error {
	return self.formLayout.ArrayField(field, self.validateField(field), self.form, self.context, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) EndArray(array *model.MetaData) error {
	err := self.formLayout.EndArray(array, self.validateGeneral(array), self.form, self.context, self.formContent)
	if array.Parent == nil && err == nil {
		return self.endForm(array)
	}
	return err
}
