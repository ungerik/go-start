package model

import (
	"fmt"
	"reflect"

	"github.com/ungerik/go-start/reflection"
)

// FieldOnlyVisitor calls its function for every struct, array and slice field.
type FieldOnlyVisitor func(field *MetaData) error

func (self FieldOnlyVisitor) BeginNamedFields(*MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) NamedField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndNamedFields(*MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) BeginIndexedFields(*MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) IndexedField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndIndexedFields(*MetaData) error {
	return nil
}

// VisitFieldType calls callback for every struct field whose type
// is asignable to the type of the callback's only argument.
// callback can return one error value or no result.
func FieldTypeVisitor(callback interface{}) Visitor {
	cv := reflect.ValueOf(callback)
	ct := cv.Type()
	if ct.NumIn() != 1 {
		panic(fmt.Errorf("model.VisitFieldType callback must have 1 argument, got %d", ct.NumIn()))
	}
	if ct.NumOut() != 0 && ct.NumOut() != 1 {
		panic(fmt.Errorf("model.VisitFieldType callback must have 0 or 1 result, got %d", ct.NumOut()))
	}
	if ct.NumOut() == 1 && ct.Out(0) != reflection.TypeOfError {
		panic(fmt.Errorf("model.VisitFieldType callback result must be of type error, got %s", ct.Out(0)))
	}
	argT := ct.In(0)
	return FieldOnlyVisitor(
		func(field *MetaData) error {
			var arg reflect.Value
			if field.Value.Type().AssignableTo(argT) {
				arg = field.Value
			} else if field.Value.CanAddr() && field.Value.Addr().Type().AssignableTo(argT) {
				arg = field.Value.Addr()
			} else {
				return nil
			}
			results := cv.Call([]reflect.Value{arg})
			if len(results) == 0 {
				return nil
			}
			err := results[0].Interface()
			if err != nil {
				return err.(error)
			}
			return nil
		},
	)
}
