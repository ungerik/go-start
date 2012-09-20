package view

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	// "mime/multipart"

	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/reflection"
	"github.com/ungerik/go-start/utils"
)

const MultipartFormData = "multipart/form-data"
const FormIDName = "gostart_form_id"

type GetFormModelFunc func(form *Form, ctx *Context) (model interface{}, err error)
type OnFormSubmitFunc func(form *Form, formModel interface{}, ctx *Context) (message string, redirect URL, err error)

type FileUploadFormModel struct {
	File model.File
}

// GetFileUploadFormModel when set at Form.GetModel will create a
// file upload form. Form.Enctype will be set to "multipart/form-fata"
// *FileUploadFormModel will be passed as formModel at Form.OnSubmit
func GetFileUploadFormModel(form *Form, ctx *Context) (interface{}, error) {
	form.Enctype = MultipartFormData
	return &FileUploadFormModel{}, nil
}

func FormModel(model interface{}) GetFormModelFunc {
	debug.Nop()
	return func(form *Form, ctx *Context) (interface{}, error) {
		return model, nil
	}
}

func OnFormSubmitSaveModel(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
	return "", nil, formModel.(mongo.Document).Save()
}

// OnFormSubmit wraps onSubmitFunc as a OnFormSubmitFunc with the
// arguments and results of OnFormSubmitFunc optional for onSubmitFunc.
// The order of the arguments and results is also free to choose.
func OnFormSubmit(onSubmitFunc interface{}) OnFormSubmitFunc {
	v := reflect.ValueOf(onSubmitFunc)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("OnFormSubmitBindURLArgs: onSubmitFunc must be a function, got %T", onSubmitFunc))
	}
	iForm := -1
	for i := 0; i < t.NumIn(); i++ {
		if t.In(i) == reflect.TypeOf((*Form)(nil)) {
			iForm = i
			break
		}
	}
	iFormModel := -1
	for i := 0; i < t.NumIn(); i++ {
		if t.In(i) == reflect.TypeOf((*interface{})(nil)).Elem() {
			iFormModel = i
			break
		}
	}
	iContext := -1
	for i := 0; i < t.NumIn(); i++ {
		if t.In(i) == reflect.TypeOf((*Context)(nil)) {
			iContext = i
			break
		}
	}
	iMessage := -1
	for i := 0; i < t.NumOut(); i++ {
		if t.Out(i) == reflect.TypeOf((*string)(nil)).Elem() {
			iMessage = i
			break
		}
	}
	iRedirect := -1
	for i := 0; i < t.NumOut(); i++ {
		if t.Out(i) == reflect.TypeOf((*URL)(nil)) {
			iRedirect = i
			break
		}
	}
	iError := -1
	for i := 0; i < t.NumOut(); i++ {
		if t.Out(i) == reflection.TypeOfError {
			iError = i
			break
		}
	}
	return func(form *Form, formModel interface{}, ctx *Context) (message string, redirect URL, err error) {
		args := make([]reflect.Value, t.NumIn())
		if iForm != -1 {
			args[iForm] = reflect.ValueOf(form)
		}
		if iFormModel != -1 {
			args[iFormModel] = reflect.ValueOf(formModel)
		}
		if iContext != -1 {
			args[iContext] = reflect.ValueOf(ctx)
		}

		results := v.Call(args)

		if iMessage != -1 {
			message = results[iMessage].Interface().(string)
		}
		if iRedirect != -1 {
			redirect, _ = results[iRedirect].Interface().(URL)
		}
		if iError != -1 {
			err, _ = results[iError].Interface().(error)
		}
		return message, redirect, err
	}
}

/*
OnFormSubmitBindURLArgs creates a View where onSubmitFunc is called from View.Render
with Context.URLArgs converted to function arguments.
The first onSubmitFunc argument can be of type *Context, but this is optional.
All further arguments will be converted from the corresponding Context.URLArgs
string to their actual type.
Conversion to integer, float, string and bool arguments are supported.
If there are more URLArgs than function args, then supernumerous URLArgs will
be ignored.
If there are less URLArgs than function args, then the function args will
receive their types zero value.
The first result value of onSubmitFunc has to be of type View,
the second result is optional and of type error.

Example:

	view.OnFormSubmitBindURLArgs(func(i int, b bool) view.View {
		return view.Printf("ctx.URLArgs[0]: %d, ctx.URLArgs[1]: %b", i, b)
	})
*/
func _OnFormSubmitBindURLArgs(onSubmitFunc interface{}) OnFormSubmitFunc {
	panic("todo")
}

/*
FormLayout is responsible for creating and structuring all dynamic content
of the form including the submit button.
It uses Form.GetFieldFactory() to create the field views.
*/
type FormLayout interface {
	GetDefaultInputSize(metaData *model.MetaData) int

	BeginFormContent(form *Form, ctx *Context, formContent *Views) error
	// SubmitSuccess will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// not returned an error.
	SubmitSuccess(message string, form *Form, ctx *Context, formContent *Views) error
	// SubmitError will be called before EndFormContent if there were no
	// validation errors of the posted form data and Form.OnSubmit has
	// returned an error.
	SubmitError(message string, form *Form, ctx *Context, formContent *Views) error
	EndFormContent(fieldValidationErrs, generalValidationErrs []error, form *Form, ctx *Context, formContent *Views) error

	BeginNamedFields(namedFields *model.MetaData, form *Form, ctx *Context, formContent *Views) error
	NamedField(field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error
	EndNamedFields(namedFields *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error

	BeginIndexedFields(indexedFields *model.MetaData, form *Form, ctx *Context, formContent *Views) error
	IndexedField(field *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error
	EndIndexedFields(indexedFields *model.MetaData, validationErr error, form *Form, ctx *Context, formContent *Views) error
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

/*
Form creates an input form for a data model.

Rules for form fields:

	TODO
*/
type Form struct {
	ViewBaseWithId
	Class  string
	Style  string
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

	// GetModel returns the data-model used to create the form fields
	// and will receive changes from a form submit.
	// Thus, the model must be setable by reflection.
	// Be careful to not return the same model object at every call,
	// because the changed model from a submit will then be re-used
	// at a later form display which is usually not wanted.
	// Use the wrapper function FormModel() only when the form is
	// embedded in a DynamicView and the model argument is dynamically
	// created per View render.
	GetModel GetFormModelFunc

	// ModelFieldAuth maps a field selector to an Authenticator.
	// Fields that match the selector will only be displayed if the
	// Authenticator returns true.
	ModelFieldAuth map[string]Authenticator

	/*
		OnSubmit is called after the form was submitted by the user
		and did not contain any validation errors.
		If the OnSubmit result err is nil and the redirect result or Form.Redirect
		is not nil, then a HTTP redirect will be issued.
		Else if the OnSubmit result message is not empty, it will be displayed
		in the form regardless of the result err.
		If the results message and err are nil, then Form.SuccessMessage
		will be	displayed in the form.
		If the result err is not nil and message is empty, then err will be
		displayed in the form. Else message will be displayed.
		This can be used to return a readable error message for the user
		via message and an internal error for logging via err.
	*/
	OnSubmit OnFormSubmitFunc

	// NewFieldInput is a hook that will be called instead of FieldFactory.NewInput
	// if not nil. Can be used for form specific alterations to the form layout.
	NewFieldInput func(withLabel bool, metaData *model.MetaData, form *Form) (View, error)

	ModelMaxDepth            int      // if zero, no depth limit
	ExcludedFields           []string // Use point notation for nested fields. In case of arrays/slices use wildcards
	HiddenFields             []string // Use point notation for nested fields. In case of arrays/slices use wildcards
	DisabledFields           []string // Use point notation for nested fields
	RequiredFields           []string // Also available as static struct field tag. Use point notation for nested fields
	Labels                   map[string]string
	FieldDescriptions        map[string]string
	InputSizes               map[string]int
	ErrorMessageClass        string // If empty, Config.Form.DefaultErrorMessageClass will be used
	SuccessMessageClass      string // If empty, Config.Form.DefaultSuccessMessageClass will be used
	FieldDescriptionClass    string
	GeneralErrorOnFieldError bool
	RequiredMarker           View // If nil, Config.Form.DefaultRequiredMarker will be used
	SuccessMessage           string
	SubmitButtonText         string
	SubmitButtonClass        string
	SubmitButtonConfirm      string // Will add a confirmation dialog for onclick
	Redirect                 URL    // 302 redirect after successful OnSubmit()
	ShowRefIDs               bool
	Enctype                  string
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

func (self *Form) IsFieldRequired(field *model.MetaData) bool {
	if val, ok := field.ModelValue(); ok && val.Required(field) {
		return true
	}
	return field.SelectorsMatch(self.RequiredFields)
}

func (self *Form) IsFieldDisabled(field *model.MetaData) bool {
	return field.BoolAttrib(StructTagKey, "disabled") || field.SelectorsMatch(self.DisabledFields)
}

// IsFieldHidden returns if a hidden input element will be created for a form field.
// Hidden fields are not validated.
func (self *Form) IsFieldHidden(field *model.MetaData) bool {
	return field.BoolAttrib(StructTagKey, "hidden") || field.SelectorsMatch(self.HiddenFields)
}

// IsFieldExcluded returns weather a field will be excluded from the form.
// Fields will be excluded, if their selector matches one in Form.ExcludedFields
// or if a matching Authenticator from Form.ModelFieldAuth returns false.
// A field is also excluded when its parent field is excluded.
// This function is not restricted to model.Value, it works with all struct fields.
// This way a whole sub struct an be excluded by adding its selector to Form.ExcludedFields.
func (self *Form) IsFieldExcluded(field *model.MetaData, ctx *Context) bool {
	if field.Parent == nil {
		return false // can't exclude root
	}
	if self.IsFieldExcluded(field.Parent, ctx) || field.SelectorsMatch(self.ExcludedFields) {
		return true
	}
	if len(self.ModelFieldAuth) > 0 {
		auth, hasAuth := self.ModelFieldAuth[field.Selector()]
		if !hasAuth {
			auth, hasAuth = self.ModelFieldAuth[field.WildcardSelector()]
		}
		if hasAuth {
			ok, err := auth.Authenticate(ctx)
			if err != nil {
				fmt.Println("Error in view.Form.IsFieldExcluded(): " + err.Error())
			}
			if !ok {
				return true
			}
		}
	}
	if len(Config.NamedAuthenticators) > 0 {
		if authAttrib, ok := field.Attrib(StructTagKey, "auth"); ok {
			for _, name := range strings.Split(authAttrib, ",") {
				// if multi := strings.Split(name, "+"); len(multi) > 1 {
				// 	// Needs to pass all sub-Authenticators
				// 	for _, name := range multi {
				// 		if auth, ok := NamedAuthenticator(name); ok {
				// 			ok, err := auth.Authenticate(response)
				// 			if err != nil {
				// 				fmt.Println("Error in view.Form.IsFieldExcluded(): " + err.Error())
				// 			}
				// 			if !ok {
				// 				return true
				// 			}
				// 		}
				// 	}
				// } else {
				if auth, ok := NamedAuthenticator(name); ok {
					ok, err := auth.Authenticate(ctx)
					if ok {
						// Only needs to pass one Authenticator			
						return false
					}
					if err != nil {
						log.Println("Error in view.Form.IsFieldExcluded(): " + err.Error())
					}
				}
				// }
			}
			// No Authenticators passed, thus exclude
			return true
		}
	}
	return false
}

// IsFieldVisible returns if the FormFieldFactory can create an input widget
// for the field, and if the field neither hidden or excluded.
// Only visible fields are validated.
// IsFieldExcluded and IsFieldHidden have different semantics.
func (self *Form) IsFieldVisible(field *model.MetaData, ctx *Context) bool {
	return self.GetFieldFactory().CanCreateInput(field, self) &&
		!self.IsFieldExcluded(field, ctx) &&
		!self.IsFieldHidden(field)
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
	if desc, ok := metaData.Attrib(StructTagKey, "description"); ok {
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
	if sizeStr, ok := metaData.Attrib(StructTagKey, "size"); ok {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			panic("Error in Form.GetInputSize(): " + err.Error())
		}
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
	if label, ok := metaData.Attrib(StructTagKey, "label"); ok {
		return label
	}
	return strings.Replace(metaData.NameOrIndex(), "_", " ", -1)
}

func (self *Form) IsPost(request *Request) bool {
	return request.Method == "POST" && request.FormValue(FormIDName) == self.FormID
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

func (self *Form) InputFieldPlaceholder(metaData *model.MetaData) string {
	if placeholder, ok := metaData.Attrib(StructTagKey, "placeholder"); ok {
		return placeholder
	}
	return ""
}

func (self *Form) newFieldInput(withLabel bool, metaData *model.MetaData) (View, error) {
	if self.NewFieldInput != nil {
		return self.NewFieldInput(withLabel, metaData, self)
	}
	return self.GetFieldFactory().NewInput(withLabel, metaData, self)
}

func (self *Form) FieldInputClass(metaData *model.MetaData) string {
	class, _ := metaData.Attrib(StructTagKey, "class")
	return class
}

func (self *Form) Render(ctx *Context) (err error) {
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
	isPost := self.IsPost(ctx.Request)

	if self.GetModel == nil {
		submitButton := self.GetFieldFactory().NewSubmitButton(self.GetSubmitButtonText(), self.SubmitButtonConfirm, self)
		content = append(content, submitButton)
	} else {
		formModel, err = self.GetModel(self, ctx)
		if err != nil {
			return err
		}
		v := reflect.ValueOf(formModel)
		if !(v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) &&
			!(v.Kind() == reflect.Map && v.Type().Key().Kind() == reflect.String) &&
			!(v.Type() == model.DynamicValueType) {
			panic(fmt.Errorf("Invalid form model type: %T", formModel))
		}
		if isPost {
			setPostValues := &setPostValuesStructVisitor{
				form:      self,
				formModel: formModel,
				ctx:       ctx,
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
			ctx:         ctx,
			isPost:      isPost,
		}
		err = model.Visit(formModel, validateAndFormLayout)
		if err != nil {
			return err
		}
		hasErrors = len(validateAndFormLayout.fieldValidationErrors) > 0 || len(validateAndFormLayout.generalValidationErrors) > 0
	}

	if isPost && !hasErrors {
		message, redirect, err := self.OnSubmit(self, formModel, ctx)
		if err == nil {
			if redirect == nil {
				redirect = self.Redirect
			}
			if redirect != nil {
				return Redirect(redirect.URL(ctx))
			}
			if message == "" {
				message = self.SuccessMessage
			}
			if message != "" {
				self.GetLayout().SubmitSuccess(message, self, ctx, &content)
			}
		} else {
			if message == "" {
				message = err.Error()
			}
			self.GetLayout().SubmitError(message, self, ctx, &content)
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
		if i := strings.Index(ctx.Request.RequestURI, "?"); i != -1 {
			action += ctx.Request.RequestURI[i:]
		}
	}
	// Hack: Value of hidden input FormIDName is not available when
	// enctype is multipart/form-data (bug?), so pass the form id as
	// URL parameter
	if self.Enctype == MultipartFormData && ctx.Request.Method != "POST" {
		action = utils.AddUrlParam(action, FormIDName, self.FormID)
	}

	ctx.Response.XML.OpenTag("form")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.AttribIfNotDefault("style", self.Style)
	ctx.Response.XML.Attrib("method", method)
	ctx.Response.XML.Attrib("action", action)
	ctx.Response.XML.AttribIfNotDefault("enctype", self.Enctype)

	if len(content) > 0 {
		content.Init(content)
		err = content.Render(ctx)
		if err != nil {
			return err
		}
	}
	ctx.Response.XML.ForceCloseTag() // form
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// setPostValuesStructVisitor

type setPostValuesStructVisitor struct {
	form      *Form
	formModel interface{}
	ctx       *Context
}

func (self *setPostValuesStructVisitor) trySetFieldValue(field *model.MetaData) error {
	if self.form.IsFieldDisabled(field) || self.form.IsFieldExcluded(field, self.ctx) {
		return nil
	}

	switch s := field.Value.Addr().Interface().(type) {
	case model.Reference:
		// we don't handle references with forms yet
		return nil

	case *model.Bool:
		s.Set(self.ctx.Request.FormValue(field.Selector()) != "")
		return nil

	case *model.Blob:
		file, _, err := self.ctx.Request.FormFile(field.Selector())
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
		file, header, err := self.ctx.Request.FormFile(field.Selector())
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
		s.SetString(self.ctx.Request.FormValue(field.Selector()))
		return nil
	}
	return nil
}

func (self *setPostValuesStructVisitor) BeginNamedFields(namedFields *model.MetaData) error {
	return nil
}

func (self *setPostValuesStructVisitor) NamedField(field *model.MetaData) error {
	return self.trySetFieldValue(field)
}

func (self *setPostValuesStructVisitor) EndNamedFields(namedFields *model.MetaData) error {
	return nil
}

func (self *setPostValuesStructVisitor) BeginIndexedFields(indexedFields *model.MetaData) error {
	if indexedFields.Kind == model.SliceKind && indexedFields.Value.CanSet() && !self.form.IsFieldExcluded(indexedFields, self.ctx) {
		if lengthStr := self.ctx.Request.FormValue(indexedFields.Selector() + ".length"); lengthStr != "" {
			length, err := strconv.Atoi(lengthStr)
			if err != nil {
				panic(err.Error())
			}
			if length != indexedFields.Value.Len() {
				indexedFields.Value.Set(utils.SetSliceLengh(indexedFields.Value, length))
				mongo.InitRefs(self.formModel)
			}
		}
	}
	return nil
}

func (self *setPostValuesStructVisitor) IndexedField(field *model.MetaData) error {
	return self.trySetFieldValue(field)
}

func (self *setPostValuesStructVisitor) EndIndexedFields(indexedFields *model.MetaData) error {
	if indexedFields.Kind == model.SliceKind && indexedFields.Value.CanSet() && !self.form.IsFieldExcluded(indexedFields, self.ctx) {
		indexedFields.Value.Set(utils.DeleteEmptySliceElementsVal(indexedFields.Value))
	}
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
	ctx         *Context
	isPost      bool

	// Output
	fieldValidationErrors   []error
	generalValidationErrors []error
}

func (self *validateAndFormLayoutStructVisitor) validateField(field *model.MetaData) error {
	if value, ok := field.ModelValue(); ok {
		err := value.Validate(field)
		if err == nil && value.IsEmpty() && self.form.IsFieldRequired(field) {
			err = model.NewRequiredError(field)
		}
		if err != nil {
			self.fieldValidationErrors = append(self.fieldValidationErrors, err)
		}
		return err
	}
	return nil
}

func (self *validateAndFormLayoutStructVisitor) validateGeneral(data *model.MetaData) error {
	if validator, ok := data.ModelValidator(); ok {
		err := validator.Validate(data)
		if err != nil {
			self.generalValidationErrors = append(self.generalValidationErrors, err)
		}
		return err
	}
	return nil
}

func (self *validateAndFormLayoutStructVisitor) endForm(data *model.MetaData) (err error) {
	return self.formLayout.EndFormContent(self.fieldValidationErrors, self.generalValidationErrors, self.form, self.ctx, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) BeginNamedFields(namedFields *model.MetaData) error {
	if namedFields.Parent == nil {
		err := self.formLayout.BeginFormContent(self.form, self.ctx, self.formContent)
		if err != nil {
			return err
		}
	}
	return self.formLayout.BeginNamedFields(namedFields, self.form, self.ctx, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) NamedField(field *model.MetaData) error {
	if self.form.IsFieldExcluded(field, self.ctx) {
		return nil
	}
	var validationErr error
	if self.isPost && self.form.IsFieldVisible(field, self.ctx) {
		validationErr = self.validateField(field)
	}
	return self.formLayout.NamedField(field, validationErr, self.form, self.ctx, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) EndNamedFields(namedFields *model.MetaData) error {
	var validationErr error
	if self.isPost {
		validationErr = self.validateGeneral(namedFields)
	}
	err := self.formLayout.EndNamedFields(namedFields, validationErr, self.form, self.ctx, self.formContent)
	if namedFields.Parent == nil && err == nil {
		return self.endForm(namedFields)
	}
	return err
}

func (self *validateAndFormLayoutStructVisitor) BeginIndexedFields(indexedFields *model.MetaData) error {
	if indexedFields.Parent == nil {
		err := self.formLayout.BeginFormContent(self.form, self.ctx, self.formContent)
		if err != nil {
			return err
		}
	}
	if self.form.IsFieldExcluded(indexedFields, self.ctx) {
		return nil
	}
	// Add an empty indexedFields field to generate one extra input row for indexedFieldss
	if indexedFields.Kind == model.SliceKind && !self.isPost && indexedFields.Value.CanSet() {
		indexedFields.Value.Set(utils.AppendEmptySliceField(indexedFields.Value))
		mongo.InitRefs(self.formModel)
	}
	return self.formLayout.BeginIndexedFields(indexedFields, self.form, self.ctx, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) IndexedField(field *model.MetaData) error {
	if self.form.IsFieldExcluded(field, self.ctx) {
		return nil
	}
	var validationErr error
	if self.isPost && self.form.IsFieldVisible(field, self.ctx) {
		validationErr = self.validateField(field)
	}
	return self.formLayout.IndexedField(field, validationErr, self.form, self.ctx, self.formContent)
}

func (self *validateAndFormLayoutStructVisitor) EndIndexedFields(indexedFields *model.MetaData) error {
	if self.form.IsFieldExcluded(indexedFields, self.ctx) {
		return nil
	}
	var validationErr error
	if self.isPost {
		validationErr = self.validateGeneral(indexedFields)
	}
	err := self.formLayout.EndIndexedFields(indexedFields, validationErr, self.form, self.ctx, self.formContent)
	if indexedFields.Parent == nil && err == nil {
		return self.endForm(indexedFields)
	}
	return err
}
