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

// FieldTypeVisitor calls callback for every struct field whose type
// is asignable to the type of the callback's first argument.
// callback can have a second, optional argument of type *MetaData.
// callback can return one error value or no result.
func FieldTypeVisitor(callback interface{}) Visitor {
	cv := reflect.ValueOf(callback)
	ct := cv.Type()
	if ct.NumIn() != 1 && ct.NumIn() != 2 {
		panic(fmt.Errorf("model.FieldTypeVisitor callback must have one or two arguments, got %d", ct.NumIn()))
	}
	if ct.NumIn() != 2 && ct.In(1) != reflect.TypeOf((*MetaData)(nil)) {
		panic(fmt.Errorf("model.FieldTypeVisitor callback's second argument must be of type *MetaData, got %s", ct.In(1)))
	}
	if ct.NumOut() != 0 && ct.NumOut() != 1 {
		panic(fmt.Errorf("model.FieldTypeVisitor callback must have zero or one results, got %d", ct.NumIn()))
	}
	if ct.NumOut() == 1 && ct.Out(0) != reflection.TypeOfError {
		panic(fmt.Errorf("model.FieldTypeVisitor callback result must be of type error, got %s", ct.Out(0)))
	}
	argT := ct.In(0)
	return FieldOnlyVisitor(
		func(field *MetaData) error {
			var args []reflect.Value
			if field.Value.Type().AssignableTo(argT) {
				args = []reflect.Value{field.Value}
			} else if field.Value.CanAddr() && field.Value.Addr().Type().AssignableTo(argT) {
				args = []reflect.Value{field.Value.Addr()}
			} else {
				return nil
			}
			if ct.NumIn() == 2 {
				args = append(args, reflect.ValueOf(field))
			}
			results := cv.Call(args)
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
