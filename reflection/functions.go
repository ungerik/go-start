package reflection

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode"

	// "github.com/ungerik/go-start/errs"
)

// TypeOfError is the built-in error type
var TypeOfError = reflect.TypeOf(func(error) {}).In(0)

func GenericSlice(sliceOrArray interface{}) []interface{} {
	v := reflect.ValueOf(sliceOrArray)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic(fmt.Errorf("Expected slice or array, got %T", sliceOrArray))
	}
	l := v.Len()
	result := make([]interface{}, l)
	for i := 0; i < l; i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
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

func FindFlattenedStructField(t reflect.Type, matchFunc MatchStructFieldFunc) *reflect.StructField {
	fieldCount := t.NumField()
	for i := 0; i < fieldCount; i++ {
		field := t.Field(i)
		if IsExportedField(field) {
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

// func CallMethod(object interface{}, method string, args ...interface{}) (results []interface{}, err error) {
// 	m := reflect.ValueOf(object).MethodByName(method)
// 	if !m.IsValid() {
// 		return nil, fmt.Errorf("%T has no method %s", object, method)
// 	}

// 	a := make([]reflect.Value, len(args))
// 	for i, arg := range args {
// 		a[i] = reflect.ValueOf(arg)
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = errs.Format("utils.CallMethod() recovered from: %v", r)
// 		}
// 	}()
// 	r := m.Call(a)

// 	results = make([]interface{}, len(r))
// 	for i, result := range r {
// 		results[i] = result.Interface()
// 	}
// 	return results, nil
// }

// func CallMethod1(object interface{}, method string, args ...interface{}) (result interface{}, err error) {
// 	results, err := CallMethod(object, method, args...)
// 	if err != nil {
// 		return
// 	}
// 	if len(results) != 1 {
// 		return nil, fmt.Errorf("One result expected from method %s of %T, %d returned", method, object, len(results))
// 	}
// 	return results[0], nil
// }

func IsDefaultValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := reflect.ValueOf(value); v.Kind() {
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

	panic(fmt.Errorf("Unknown value kind %T", value))
}

// IsNilOrWrappedNil returns if i is nil, or wraps a nil pointer
// in a non nil interface.
func IsNilOrWrappedNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil() || IsNilOrWrappedNil(v.Elem().Interface())
	}
	return false
}

// ExportedStructFields returns a map of exported struct fields
// with the field name as key and the reflect.Value as map value.
// s can be a struct, a struct pointer or a reflect.Value of
// a struct or struct pointer.
func ExportedStructFields(s interface{}) map[string]reflect.Value {
	result := make(map[string]reflect.Value)
	v, ok := s.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(s)
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if IsExportedField(field) {
			result[field.Name] = v.Field(i)
		}
	}
	return result
}

func IsExportedName(name string) bool {
	return name != "" && unicode.IsUpper(rune(name[0]))
}

func IsExportedField(structField reflect.StructField) bool {
	return structField.PkgPath == ""
}

// CopyExportedStructFields copies all exported struct fields from src
// that are assignable to their name siblings at dstPtr to dstPtr.
// src can be a struct or a pointer to a struct, dstPtr must be
// a pointer to a struct.
func CopyExportedStructFields(src, dstPtr interface{}) (copied int) {
	vsrc := reflect.ValueOf(src)
	if vsrc.Kind() == reflect.Ptr {
		vsrc = vsrc.Elem()
	}
	vdst := reflect.ValueOf(dstPtr).Elem()
	tsrc := reflect.TypeOf(src)
	for i := 0; i < tsrc.NumField(); i++ {
		tsrcfield := tsrc.Field(i)
		if IsExportedField(tsrcfield) {
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
		return int(i), nil

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
		return uint(i), nil

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

func CanStringToValueOfType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Float32,
		reflect.Float64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return true
	}
	return false
}

// SmartCopy copies the struct or map fields from source
// to equally named struct or map fields of resultRef,
// by dereferencing source and resultRef if necessary
// to find a matching assignable type.
// If resultRef is a pointer to a pointer type,
// then a new instance of that pointer type is created.
// SmartCopy is typically used for iterators with
// a method Next(resultRefRef interface{}) bool.
func SmartCopy(source, resultRef interface{}) {
	resultRefVal := reflect.ValueOf(resultRef)
	if resultRefVal.Kind() != reflect.Ptr && resultRefVal.Kind() != reflect.Map {
		panic(fmt.Errorf("reflection.SmartCopy(): resultRef must be a pointer or a map, got %T", resultRef))
	}
	var resultVal reflect.Value
	if resultRefVal.Kind() == reflect.Ptr {
		resultVal = resultRefVal.Elem()
	} else {
		resultVal = resultRefVal
	}
	if resultVal.Kind() == reflect.Ptr {
		// If resultRef points to a pointer,
		// create an instance of the pointer type
		resultVal.Set(reflect.New(resultVal.Type().Elem()))
		resultVal = resultVal.Elem()
	}
	if !smartCopyVals(reflect.ValueOf(source), resultVal) {
		panic(fmt.Errorf("reflection.SmartCopy(): Can't copy %T to %T", source, resultRef))
	}
}

func smartCopyVals(sourceVal, resultVal reflect.Value) bool {
	switch {
	case sourceVal.Type().AssignableTo(resultVal.Type()):
		resultVal.Set(sourceVal)
		return true

	case sourceVal.Kind() == reflect.Ptr && sourceVal.Elem().Type().AssignableTo(resultVal.Type()):
		smartCopyVals(sourceVal.Elem(), resultVal)
		return true

	case sourceVal.Kind() == reflect.Struct && resultVal.Kind() == reflect.Struct:
		CopyExportedStructFields(sourceVal.Interface(), resultVal.Interface())
		return true

	case sourceVal.Kind() == reflect.Map && sourceVal.Type().Key().Kind() == reflect.String &&
		resultVal.Kind() == reflect.Struct:

		resultMap := ExportedStructFields(resultVal)
		for _, key := range sourceVal.MapKeys() {
			if dst, ok := resultMap[key.String()]; ok {
				src := sourceVal.MapIndex(key)
				if src.Type().AssignableTo(dst.Type()) {
					dst.Set(src)
				}
			}
		}
		return true

	case sourceVal.Kind() == reflect.Struct &&
		resultVal.Kind() == reflect.Map && resultVal.Type().Key().Kind() == reflect.String:

		sourceMap := ExportedStructFields(sourceVal)
		for key, src := range sourceMap {
			dst := resultVal.MapIndex(reflect.ValueOf(key))
			if dst.IsValid() && src.Type().AssignableTo(dst.Type()) {
				dst.Set(src)
			}
		}
		return true

	case sourceVal.Kind() == reflect.Map && sourceVal.Type().Key().Kind() == reflect.String &&
		resultVal.Kind() == reflect.Map && resultVal.Type().Key().Kind() == reflect.String &&
		sourceVal.Type().Elem().AssignableTo(resultVal.Type().Elem()):

		for _, key := range sourceVal.MapKeys() {
			if dst := resultVal.MapIndex(key); dst.IsValid() {
				dst.Set(sourceVal.MapIndex(key))
			}
		}
		return true
	}

	return false
}
