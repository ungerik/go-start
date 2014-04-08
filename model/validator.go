package model

import (
	"reflect"
)

type Validator interface {
	// Validate returns an error if metaData is not valid.
	// In case of multiple errors errs.ErrSlice is returned.
	Validate(metaData *MetaData) error
}

var ValidatorType = reflect.TypeOf((*Validator)(nil)).Elem()

func IsValidator(v reflect.Value) bool {
	return v.Type().Implements(ValidatorType) || (v.CanAddr() && v.Addr().Type().Implements(ValidatorType))
}
