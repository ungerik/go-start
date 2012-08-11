package utils

import (
	"fmt"
	"github.com/ungerik/go-start/errs"
	"reflect"
	"strconv"
	"unicode"
)

// IsErrorType checks if t is the built-in type error.
func IsErrorType(t reflect.Type) bool {
	return t == reflect.TypeOf(func(error) {}).In(0)
}

/*
DereferenceValue recursively dereferences v if it is a pointer or interface.
It will return ok == false if nil is encountered.
*/
func DereferenceValue(v reflect.Value) (result reflect.Value, ok bool) {
	k := v.Kind()
	if k == reflect.Ptr || k == reflect.Interface {
		if v.IsNil() {
			return v, false
		} else {
			return DereferenceValue(v.Elem())
		}
	}
	return v, true
}

type MatchStructFieldFunc func(field *reflect.StructField) bool

func IsExportedName(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}

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
	if value == nil {
		return true
	}
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
		return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())

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

func ExportedStructFields(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if IsExportedName(field.Name) {
			result[field.Name] = v.Field(i).Interface()
		}
	}
	return result
}

func CopyExportedStructFields(src, dst interface{}) (copied int) {
	vsrc := reflect.ValueOf(src)
	vdst := reflect.ValueOf(dst)
	tsrc := reflect.TypeOf(src)
	for i := 0; i < tsrc.NumField(); i++ {
		tsrcfield := tsrc.Field(i)
		if IsExportedName(tsrcfield.Name) {
			dstfield := vdst.FieldByName(tsrcfield.Name)
			if dstfield.IsValid() && dstfield.CanSet() && tsrcfield.Type.AssignableTo(dstfield.Type()) {
				dstfield.Set(vsrc.Field(i))
				copied++
			}
		}
	}
	return copied
}

func StringToValueOfType(s string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.String:
		return s, nil

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return b, nil

	case reflect.Float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		return float32(f), nil

	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return f, nil

	case reflect.Int:
		i, err := strconv.ParseInt(s, 0, 0)
		if err != nil {
			return nil, err
		}
		return i, nil

	case reflect.Int8:
		i, err := strconv.ParseInt(s, 0, 8)
		if err != nil {
			return nil, err
		}
		return int8(i), nil

	case reflect.Int16:
		i, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			return nil, err
		}
		return int16(i), nil

	case reflect.Int32:
		i, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return nil, err
		}
		return int32(i), nil

	case reflect.Int64:
		i, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return nil, err
		}
		return int64(i), nil

	case reflect.Uint:
		i, err := strconv.ParseUint(s, 0, 0)
		if err != nil {
			return nil, err
		}
		return i, nil

	case reflect.Uint8:
		i, err := strconv.ParseUint(s, 0, 8)
		if err != nil {
			return nil, err
		}
		return uint8(i), nil

	case reflect.Uint16:
		i, err := strconv.ParseUint(s, 0, 16)
		if err != nil {
			return nil, err
		}
		return uint16(i), nil

	case reflect.Uint32:
		i, err := strconv.ParseUint(s, 0, 32)
		if err != nil {
			return nil, err
		}
		return uint32(i), nil

	case reflect.Uint64:
		i, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return nil, err
		}
		return uint64(i), nil
	}

	return nil, fmt.Errorf("StringToValueOfType: can't convert string to type %s", t)
}
