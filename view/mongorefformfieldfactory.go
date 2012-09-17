package view

import (
	"fmt"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
)

type MongoRefFormFieldFactory struct {
	FormFieldFactoryWrapper

	DocIterator model.Iterator
	// GetName returns the name of a referenced document.
	// If GetName is nil, then the String() method of
	// the doc will be used.
	GetName func(doc interface{}) string
}

func (self *MongoRefFormFieldFactory) CanCreateInput(metaData *model.MetaData, form *Form) bool {
	if metaData.Value.CanAddr() {
		if _, ok := metaData.Value.Addr().Interface().(*mongo.Ref); ok {
			return true
		}
	}
	return self.FormFieldFactoryWrapper.Wrapped.CanCreateInput(metaData, form)
}

func (self *MongoRefFormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *Form) (View, error) {
	if mongoRef, ok := metaData.Value.Addr().Interface().(*mongo.Ref); ok {
		var options []string
		for doc := self.DocIterator.Next(); doc != nil; doc = self.DocIterator.Next() {
			if self.GetName != nil {
				options = append(options, self.GetName(doc))
			} else if s, ok := doc.(fmt.Stringer); ok {
				options = append(options, s.String())
			} else {
				panic(errs.Format("MongoRefFormFieldFactory: %T must implement fmt.String", doc))
			}
		}

		mongoRef.Get()
		if len(options) == 0 || options[0] != "" {
			options = append([]string{""}, options...)
		}
		input := &Select{
			Class: form.FieldInputClass(metaData),
			Name:  metaData.Selector(),
			// Model:    &StringsSelectModel{options, s.Get()},
			Disabled: form.IsFieldDisabled(metaData),
			Size:     1,
		}
		return input, nil
	}

	return self.FormFieldFactoryWrapper.Wrapped.NewInput(withLabel, metaData, form)
}
