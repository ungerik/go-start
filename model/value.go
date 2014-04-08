package model

import (
	"reflect"
)

type Value interface {
	String() string
	// SetString returns only error from converting str to the
	// underlying value type.
	// It does not return validation errors of the converted value.
	SetString(str string) (strconvErr error)
	IsEmpty() bool
	Required(metaData *MetaData) bool
	Validator
}

var ValueType = reflect.TypeOf((*Value)(nil)).Elem()

func IsValue(v reflect.Value) bool {
	return v.Type().Implements(ValueType) || (v.CanAddr() && v.Addr().Type().Implements(ValueType))
}
