package reflection

import (
	"reflect"
)

type DefaultValuer interface {
	IsDefault() bool
	GetDefault() interface{}
}

var DefaultValuerType = reflect.TypeOf((*DefaultValuer)(nil)).Elem()

func SmartCastDefaultValuer(v reflect.Value) (DefaultValuer, bool) {
	// Having a pointer to a type is the most common case for methods
	if v.CanAddr() && v.Addr().Type().Implements(DefaultValuerType) {
		return v.Addr().Interface().(DefaultValuer), true
	}
	// Less common is not using a pointer but the type by value for methods
	if v.Type().Implements(DefaultValuerType) {
		return v.Interface().(DefaultValuer), true
	}
	// Or the type implements the methods by value but we have a pointer to it
	if v.Kind() == reflect.Ptr && v.Elem().Type().Implements(DefaultValuerType) {
		return v.Elem().Interface().(DefaultValuer), true
	}
	return nil, false
}

// IsDefaultValue calls SmartCastDefaultValuer to tedect if v in any kind implements
// DefaultValuer. If this is the case, then DefaultValuer.IsDefault will be returned.
// If v does not implement DefaultValuer, then the result of a comparison with the
// zero value of the type will be returned.
func IsDefaultValue(v reflect.Value) bool {
	if defaultValuer, ok := SmartCastDefaultValuer(v); ok {
		return defaultValuer.IsDefault()
	}

	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == complex(0, 0)

	case reflect.Bool:
		return v.Bool() == false

	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Slice, reflect.Map:
		return v.IsNil()

	case reflect.Ptr:
		return v.IsNil() || IsDefaultValue(v.Elem())

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !IsDefaultValue(v.Field(i)) {
				return false
			}
		}
		return true
	}

	panic("never reached")
}

// IsDefault return true if value is nil,
// or the result of IsDefaultValue for the reflect.Value of value.
func IsDefault(value interface{}) bool {
	if value == nil {
		return true
	}
	return IsDefaultValue(reflect.ValueOf(value))
}

// GetDefaultValue returns the result of the GetDefault() method
// if v or v pointer implements DefaultValuer.
// Else the zero value of v.Type() will be returned.
// Used to clone types that have internal state that represents
// a default value that is not deep zero.
func GetDefaultValue(v reflect.Value) reflect.Value {
	if defaultValuer, ok := v.Interface().(DefaultValuer); ok {
		return reflect.ValueOf(defaultValuer.GetDefault())
	}
	if v.CanAddr() {
		if defaultValuer, ok := v.Addr().Interface().(DefaultValuer); ok {
			return reflect.ValueOf(defaultValuer.GetDefault())
		}
	}
	return reflect.Zero(v.Type())
}
