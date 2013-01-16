package view

import (
	"fmt"
	"io/ioutil"
	"strconv"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
)

///////////////////////////////////////////////////////////////////////////////
// FormFieldController

// FormFieldController is a MVC controller for form fields.
type FormFieldController interface {
	// Supports returns if this controller supports a model with the given metaData.
	Supports(metaData *model.MetaData, form *Form) bool

	// NewInput creates a new form field input view for the model metaData.
	NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error)

	// SetValue sets the value of the model from HTTP POST form data.
	SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error
}

///////////////////////////////////////////////////////////////////////////////
// FormFieldControllers

type FormFieldControllers []FormFieldController

func (self FormFieldControllers) Supports(metaData *model.MetaData, form *Form) bool {
	for _, c := range self {
		if c.Supports(metaData, form) {
			return true
		}
	}
	return false
}

func (self FormFieldControllers) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	for _, c := range self {
		if c.Supports(metaData, form) {
			return c.NewInput(withLabel, metaData, form)
		}
	}
	return nil, nil
}

func (self FormFieldControllers) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	for _, c := range self {
		if c.Supports(metaData, form) {
			return c.SetValue(value, ctx, metaData, form)
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ErrFormFieldTypeNotSupported

// type ErrFormFieldTypeNotSupported struct {
// 	*model.MetaData
// }

// func (self ErrFormFieldTypeNotSupported) Error() string {
// 	return fmt.Sprintf("Type %s of form field %s not supported", self.Value.Type(), self.Selector())
// }

///////////////////////////////////////////////////////////////////////////////
// SetModelValueControllerBase

type SetModelValueControllerBase struct{}

func (self SetModelValueControllerBase) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	mdl, ok := metaData.ModelValue()
	if !ok {
		panic("metaData must be model.Value")
	}
	mdl.SetString(value)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelStringController

type ModelStringController struct {
	SetModelValueControllerBase
}

func (self ModelStringController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.String)
	return ok
}

func (self ModelStringController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	str := metaData.Value.Addr().Interface().(*model.String)
	textField := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Text:        str.Get(),
		Size:        form.GetInputSize(metaData),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if maxlen, ok, _ := str.Maxlen(metaData); ok {
		textField.MaxLength = maxlen
		if maxlen < textField.Size {
			textField.Size = maxlen
		}
	}
	if metaData.BoolAttrib("view", "search") {
		textField.Type = SearchTextField
	}
	if metaData.BoolAttrib("view", "autofocus") {
		textField.Autofocus = true
	}
	if pattern, ok := metaData.Attrib("view", "pattern"); ok {
		textField.Pattern = pattern
	}
	if withLabel {
		return AddStandardLabel(form, textField, metaData), nil
	}
	return textField, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelTextController

type ModelTextController struct {
	SetModelValueControllerBase
}

func (self ModelTextController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Text)
	return ok
}

func (self ModelTextController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	text := metaData.Value.Addr().Interface().(*model.Text)
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
	textArea := &TextArea{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Text:        text.Get(),
		Cols:        cols,
		Rows:        rows,
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
	}
	if metaData.BoolAttrib("view", "autofocus") {
		textArea.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, textArea, metaData), nil
	}
	return textArea, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelRichTextController

type ModelRichTextController struct {
	SetModelValueControllerBase
}

func (self ModelRichTextController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.RichText)
	return ok
}

func (self ModelRichTextController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	text := metaData.Value.Addr().Interface().(*model.RichText)
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
	var stylesheetURL string
	if style, ok := metaData.Attrib(StructTagKey, "stylesheetURL"); ok {
		stylesheetURL = style
	}

	input = &RichTextArea{
		Class:                form.FieldInputClass(metaData),
		Name:                 metaData.Selector(),
		Text:                 text.Get(),
		Cols:                 cols,
		Rows:                 rows,
		Disabled:             form.IsFieldDisabled(metaData),
		Placeholder:          form.InputFieldPlaceholder(metaData),
		ToolbarStylesheetURL: stylesheetURL,
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelUrlController

type ModelUrlController struct {
	SetModelValueControllerBase
}

func (self ModelUrlController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Url)
	return ok
}

func (self ModelUrlController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (View, error) {
	url := metaData.Value.Addr().Interface().(*model.Url)
	input := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Text:        url.Get(),
		Size:        form.GetInputSize(metaData),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if metaData.BoolAttrib("view", "autofocus") {
		input.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelEmailController

type ModelEmailController struct {
	SetModelValueControllerBase
}

func (self ModelEmailController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Email)
	return ok
}

func (self ModelEmailController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (View, error) {
	email := metaData.Value.Addr().Interface().(*model.Email)
	input := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Type:        EmailTextField,
		Text:        email.Get(),
		Size:        form.GetInputSize(metaData),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if metaData.BoolAttrib("view", "autofocus") {
		input.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelPasswordController

type ModelPasswordController struct {
	SetModelValueControllerBase
}

func (self ModelPasswordController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Password)
	return ok
}

func (self ModelPasswordController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	password := metaData.Value.Addr().Interface().(*model.Password)
	textField := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Type:        PasswordTextField,
		Text:        password.Get(),
		Size:        form.GetInputSize(metaData),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if maxlen, ok, _ := password.Maxlen(metaData); ok {
		textField.MaxLength = maxlen
		if maxlen < textField.Size {
			textField.Size = maxlen
		}
	}
	if metaData.BoolAttrib("view", "autofocus") {
		textField.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, textField, metaData), nil
	}
	return textField, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelPhoneController

type ModelPhoneController struct {
	SetModelValueControllerBase
}

func (self ModelPhoneController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Phone)
	return ok
}

func (self ModelPhoneController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	phone := metaData.Value.Addr().Interface().(*model.Phone)
	textField := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Text:        phone.Get(),
		Size:        form.GetInputSize(metaData),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if metaData.BoolAttrib("view", "autofocus") {
		textField.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, textField, metaData), nil
	}
	return textField, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelBoolController

type ModelBoolController struct{}

func (self ModelBoolController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Bool)
	return ok
}

func (self ModelBoolController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	b := metaData.Value.Addr().Interface().(*model.Bool)
	checkbox := &Checkbox{
		Class:    form.FieldInputClass(metaData),
		Name:     metaData.Selector(),
		Disabled: form.IsFieldDisabled(metaData),
		Checked:  b.Get(),
	}
	if withLabel {
		labelContent := form.FieldLabel(metaData)
		if form.IsFieldRequired(metaData) {
			labelContent = "<span class='required'>*</span> " + labelContent
		}
		checkbox.Label = labelContent
	}
	return checkbox, nil
}

func (self ModelBoolController) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	b := metaData.Value.Addr().Interface().(*model.Bool)
	b.Set(value != "")
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelChoiceController

type ModelChoiceController struct {
	SetModelValueControllerBase
}

func (self ModelChoiceController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Choice)
	return ok
}

func (self ModelChoiceController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	choice := metaData.Value.Addr().Interface().(*model.Choice)
	options := choice.Options(metaData)
	if len(options) == 0 || options[0] != "" {
		options = append([]string{""}, options...)
	}
	input = &Select{
		Class:    form.FieldInputClass(metaData),
		Name:     metaData.Selector(),
		Model:    &StringsSelectModel{options, choice.Get()},
		Disabled: form.IsFieldDisabled(metaData),
		Size:     1,
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelMultipleChoiceController

type ModelMultipleChoiceController struct{}

func (self ModelMultipleChoiceController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.MultipleChoice)
	return ok
}

func (self ModelMultipleChoiceController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	m := metaData.Value.Addr().Interface().(*model.MultipleChoice)
	options := m.Options(metaData)
	checkboxes := make(Views, len(options))
	for i, option := range options {
		checkboxes[i] = &Checkbox{
			Label:    option,
			Class:    form.FieldInputClass(metaData),
			Name:     fmt.Sprintf("%s_%d", metaData.Selector(), i),
			Disabled: form.IsFieldDisabled(metaData),
			Checked:  m.IsSet(option),
		}
	}
	if withLabel {
		return AddStandardLabel(form, checkboxes, metaData), nil
	}
	return checkboxes, nil
}

func (self ModelMultipleChoiceController) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	m := metaData.Value.Addr().Interface().(*model.MultipleChoice)
	options := m.Options(metaData)
	*m = nil
	for i, option := range options {
		name := fmt.Sprintf("%s_%d", metaData.Selector(), i)
		if ctx.Request.FormValue(name) != "" {
			*m = append(*m, option)
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelDynamicChoiceController

type ModelDynamicChoiceController struct {
	SetModelValueControllerBase
}

func (self ModelDynamicChoiceController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.DynamicChoice)
	return ok
}

func (self ModelDynamicChoiceController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	choice := metaData.Value.Addr().Interface().(*model.DynamicChoice)
	options := choice.Options()
	index := choice.Index()
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
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelDateController

type ModelDateController struct {
	SetModelValueControllerBase
}

func (self ModelDateController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Date)
	return ok
}

func (self ModelDateController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	date := metaData.Value.Addr().Interface().(*model.Date)
	input = Views{
		HTML("(Format: " + model.DateFormat + ")<br/>"),
		&TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        date.Get(),
			Size:        len(model.DateFormat),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
			Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
			Title:       form.FieldLabel(metaData),
		},
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelDateTimeController

type ModelDateTimeController struct {
	SetModelValueControllerBase
}

func (self ModelDateTimeController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.DateTime)
	return ok
}

func (self ModelDateTimeController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	dateTime := metaData.Value.Addr().Interface().(*model.DateTime)
	input = Views{
		HTML("(Format: " + model.DateTimeFormat + ")<br/>"),
		&TextField{
			Class:       form.FieldInputClass(metaData),
			Name:        metaData.Selector(),
			Text:        dateTime.Get(),
			Size:        len(model.DateTimeFormat),
			Disabled:    form.IsFieldDisabled(metaData),
			Placeholder: form.InputFieldPlaceholder(metaData),
			Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
			Title:       form.FieldLabel(metaData),
		},
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelFloatController

type ModelFloatController struct {
	SetModelValueControllerBase
}

func (self ModelFloatController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Float)
	return ok
}

func (self ModelFloatController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	f := metaData.Value.Addr().Interface().(*model.Float)
	textField := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Text:        f.String(),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if metaData.BoolAttrib("view", "autofocus") {
		textField.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, textField, metaData), nil
	}
	return textField, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelIntController

type ModelIntController struct {
	SetModelValueControllerBase
}

func (self ModelIntController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Int)
	return ok
}

func (self ModelIntController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	i := metaData.Value.Addr().Interface().(*model.Int)
	textField := &TextField{
		Class:       form.FieldInputClass(metaData),
		Name:        metaData.Selector(),
		Text:        i.String(),
		Disabled:    form.IsFieldDisabled(metaData),
		Placeholder: form.InputFieldPlaceholder(metaData),
		Required:    form.IsFieldRequired(metaData) && !metaData.IsSelfOrParentIndexed(),
		Title:       form.FieldLabel(metaData),
	}
	if metaData.BoolAttrib("view", "autofocus") {
		textField.Autofocus = true
	}
	if withLabel {
		return AddStandardLabel(form, textField, metaData), nil
	}
	return textField, nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelFileController

type ModelFileController struct{}

func (self ModelFileController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.File)
	return ok
}

func (self ModelFileController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	input = &FileInput{
		Class:    form.FieldInputClass(metaData),
		Name:     metaData.Selector(),
		Disabled: form.IsFieldDisabled(metaData),
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

func (self ModelFileController) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	f := metaData.Value.Addr().Interface().(*model.File)
	file, header, err := ctx.Request.FormFile(metaData.Selector())
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	f.Name = header.Filename
	f.Data = bytes
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ModelBlobController

type ModelBlobController struct{}

func (self ModelBlobController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*model.Blob)
	return ok
}

func (self ModelBlobController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	input = &FileInput{
		Class:    form.FieldInputClass(metaData),
		Name:     metaData.Selector(),
		Disabled: form.IsFieldDisabled(metaData),
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

func (self ModelBlobController) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	b := metaData.Value.Addr().Interface().(*model.Blob)
	file, _, err := ctx.Request.FormFile(metaData.Selector())
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	b.Set(bytes)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// MongoRefController

type MongoRefController struct {
	// RefSelector selects for which mongo.Ref fields this
	// FormFieldController should be used.
	// An empty string selects all mongo.Ref in the form model.
	RefSelector string

	// LabelSelector is a selector into the referenced documents to retrieve
	// the labels for the documents.
	// If labelSelector is an empty string, then a struct tag `view:"optionsLabel="`
	// will be used or DocumentLabelSelector from the mongo.Ref's collection.
	LabelSelectors []string

	// OptionsIterator iterates the possible documents for the mongo.Ref.
	// If options is nil, then all documents of the mongo.Ref's collection are used
	OptionsIterator model.Iterator
}

func (self *MongoRefController) Supports(metaData *model.MetaData, form *Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*mongo.Ref)
	return ok && (self.RefSelector == "" || self.RefSelector == metaData.Selector() || self.RefSelector == metaData.WildcardSelector())
}

func (self *MongoRefController) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (input View, err error) {
	mongoRef := metaData.Value.Addr().Interface().(*mongo.Ref)
	valuesAndLabels := []string{"", ""}
	labelSelectors := self.LabelSelectors
	if len(labelSelectors) == 0 {
		if selector, ok := metaData.Attrib("view", "optionsLabel"); ok {
			labelSelectors = []string{selector}
		}
	}
	i := self.OptionsIterator
	if i == nil {
		i = mongoRef.Collection().Iterator()
	}
	var doc mongo.DocumentBase
	for i.Next(&doc) {
		label, err := mongoRef.Collection().DocumentLabel(doc.ID, labelSelectors...)
		if err != nil {
			return nil, err
		}
		valuesAndLabels = append(valuesAndLabels, doc.ID.Hex(), label)
	}
	if i.Err() != nil {
		return nil, i.Err()
	}

	input = &Select{
		Class: form.FieldInputClass(metaData),
		Name:  metaData.Selector(),
		Model: &ValueLabelSelectModel{
			ValuesAndLabels: valuesAndLabels,
			SelectedValue:   mongoRef.StringID(),
		},
		Size:     1,
		Disabled: form.IsFieldDisabled(metaData),
	}
	if withLabel {
		return AddStandardLabel(form, input, metaData), nil
	}
	return input, nil
}

func (self *MongoRefController) SetValue(value string, ctx *Context, metaData *model.MetaData, form *Form) error {
	mongoRef := metaData.Value.Addr().Interface().(*mongo.Ref)
	return mongoRef.SetStringID(value)
}
