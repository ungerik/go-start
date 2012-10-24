package view

import (
	"bytes"
	"strings"

	"github.com/ungerik/go-start/model"
)

func NewLabeledModelView(model interface{}) *LabeledModelView {
	return &LabeledModelView{Model: model}
}

// LabeledModelView displays a Model similar to Form,
// but with immutable strings instead of input fields.
type LabeledModelView struct {
	ViewBase
	Model          interface{}
	Labels         map[string]string
	ExcludedFields []string
	HideEmpty      bool
}

// DirectFieldLabel returns a label for a model field generated from metaData.
// It creates the label only from the name or label tag of metaData,
// not including its parents.
func (self *LabeledModelView) DirectFieldLabel(metaData *model.MetaData) string {
	if label, ok := metaData.Attrib(StructTagKey, "label"); ok {
		return label
	}
	return strings.Replace(metaData.NameOrIndex(), "_", " ", -1)
}

// FieldLabel returns a label for a model field generated from metaData.
// It creates the label from the names or label tags of metaData and
// all its parents, starting with the root parent, concanated with a space
// character.
func (self *LabeledModelView) FieldLabel(metaData *model.MetaData) string {
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

// IsFieldExcluded returns weather a field will be excluded.
// Fields will be excluded, if their selector matches one in LabeledModelView.ExcludedFields.
// A field is also excluded when its parent field is excluded.
// This function is not restricted to model.Value, it works with all struct fields.
// This way a whole sub struct an be excluded by adding its selector to LabeledModelView.ExcludedFields.
func (self *LabeledModelView) IsFieldExcluded(field *model.MetaData) bool {
	if field.Parent == nil {
		return false // can't exclude root
	}
	if self.IsFieldExcluded(field.Parent) || field.SelectorsMatch(self.ExcludedFields) {
		return true
	}
	return false
}

func (self *LabeledModelView) Render(ctx *Context) (err error) {
	return model.Visit(self.Model, model.FieldOnlyVisitor(
		func(field *model.MetaData) error {
			if !self.IsFieldExcluded(field) {
				if value, ok := field.ModelValue(); ok {
					if !value.IsEmpty() || !self.HideEmpty {
						v := DIV((Config.LabeledModelViewValueClass + " " + strings.ToLower(field.NameOrIndex())), HTML(value.String()))
						label := &Label{
							Class:   Config.LabeledModelViewLabelClass + " " + strings.ToLower(field.NameOrIndex()),
							Content: Escape(self.FieldLabel(field)),
							For:     v,
						}
						return Views{label, v}.Render(ctx)
					}
				}
			}
			return nil
		},
	))
}
