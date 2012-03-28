package utils

import (
	"fmt"
	"github.com/ungerik/go-start/errs"
	"reflect"
	"unicode"
)

type MatchStructFieldFunc func(field *reflect.StructField) bool

func FindFlattenedStructField(t reflect.Type, matchFunc MatchStructFieldFunc) *reflect.StructField {
	fieldCount := t.NumField()
	for i := 0; i < fieldCount; i++ {
		field := t.Field(i)
		if unicode.IsUpper(rune(field.Name[0])) {
			if field.Anonymous {
				if field.Type.Kind() == reflect.Struct {
					result := FindFlattenedStructField(field.Type, matchFunc)
					if result != nil {
						return result
					}
				}
			} else {
				if matchFunc(&field) {
					return &field
				}
			}
		}
	}
	return nil
}

// Creates a new zero valued instance of prototype
func NewInstance(prototype interface{}) interface{} {
	t := reflect.TypeOf(prototype)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}

func CallMethod(object interface{}, method string, args ...interface{}) (results []interface{}, err error) {
	m := reflect.ValueOf(object).MethodByName(method)
	if !m.IsValid() {
		return nil, fmt.Errorf("%T has no method %s", object, method)
	}

	a := make([]reflect.Value, len(args))
	for i, arg := range args {
		a[i] = reflect.ValueOf(arg)
	}

	defer func() {
		if r := recover(); r != nil {
			err = errs.Format("utils.CallMethod() recovered from: %v", r)
		}
	}()
	r := m.Call(a)

	results = make([]interface{}, len(r))
	for i, result := range r {
		results[i] = result.Interface()
	}
	return results, nil
}

func CallMethod1(object interface{}, method string, args ...interface{}) (result interface{}, err error) {
	results, err := CallMethod(object, method, args...)
	if err != nil {
		return
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("One result expected from method %s of %T, %d returned", method, object, len(results))
	}
	return results[0], nil
}

func IsDefaultValue(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsDefaultValue(v.Elem().Interface())

	case reflect.String:
		return v.String() == ""

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Bool:
		return v.Bool() == false

	case reflect.Struct:
		// todo when struct comparison works

	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()
	}

	return false
}

// Non nil interfaces can wrap nil values. Comparing the interface to nil, won't return if the wrapped value is nil.
func IsDeepNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil() || IsDeepNil(v.Elem().Interface())
	}
	return false
}
