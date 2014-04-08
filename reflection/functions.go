package reflection

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode"

	// "github.com/ungerik/go-start/debug"
)

// Built-in types
var (
	// TypeOfError is the built-in error type
	TypeOfError = reflect.TypeOf((*error)(nil)).Elem()

	// TypeOfInterface is the type of an empty interface{}
	TypeOfInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

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

// Implements returns if the type of v or a pointer to the type of v
// implements an interfaceType.
func Implements(v reflect.Value, interfaceType reflect.Type) bool {
	// Having a pointer to a type is the most common case for methods
	if v.CanAddr() && v.Addr().Type().Implements(interfaceType) {
		return true
	}
	// Less common is not using a pointer but the type by value for methods
	if v.Type().Implements(interfaceType) {
		return true
	}
	// Type implements method by value but we have a pointer to it
	if v.Kind() == reflect.Ptr && v.Elem().Type().Implements(interfaceType) {
		return true
	}
	return false
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

/*
ExportedStructFields returns a map from exported struct field names to values,
inlining anonymous sub-structs so that their field names are available
at the base level.
Example:
	type A struct {
		X int
	}
	type B Struct {
		A
		Y int
	}
	// Yields X and Y instead of A and Y:
	InlineAnonymousStructFields(reflect.ValueOf(B{}))
*/
func ExportedStructFields(v reflect.Value) map[string]reflect.Value {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("Expected a struct, got %s", t))
	}
	result := make(map[string]reflect.Value)
	exportedStructFields(v, t, result)
	return result
}

func exportedStructFields(v reflect.Value, t reflect.Type, result map[string]reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		if IsExportedField(structField) {
			if structField.Anonymous && structField.Type.Kind() == reflect.Struct {
				exportedStructFields(v.Field(i), structField.Type, result)
			} else {
				result[structField.Name] = v.Field(i)
			}
		}
	}
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
	case reflect.String:
		return v.Len() == 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Bool:
		return v.Bool() == false

	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Slice, reflect.Map:
		return v.IsNil()

	case reflect.Struct:
		return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
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

// GetStruct returns reflect.Value for the struct found in s.
// s can be a struct, a struct pointer or a reflect.Value of
// a struct or struct pointer.
func GetStruct(s interface{}) reflect.Value {
	v, ok := s.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(s)
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
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
	return CopyExportedStructFieldsVal(vsrc, vdst)
}

func CopyExportedStructFieldsVal(src, dst reflect.Value) (copied int) {
	if src.Kind() != reflect.Struct {
		panic(fmt.Errorf("CopyExportedStructFieldsVal: src must be struct, got %s", src.Type()))
	}
	if dst.Kind() != reflect.Struct {
		panic(fmt.Errorf("CopyExportedStructFieldsVal: dst must be struct, got %s", dst.Type()))
	}
	if !dst.CanSet() {
		panic(fmt.Errorf("CopyExportedStructFieldsVal: dst (%s) is not set-able", dst.Type()))
	}
	srcFields := ExportedStructFields(src)
	dstFields := ExportedStructFields(dst)
	for name, srcV := range srcFields {
		if dstV, ok := dstFields[name]; ok {
			if srcV.Type().AssignableTo(dstV.Type()) {
				dstV.Set(srcV)
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

// SetStructZero sets all elements of a struct to their zero values.
func SetStructZero(structVal reflect.Value) {
	t := structVal.Type()
	for i := 0; i < t.NumField(); i++ {
		if IsExportedField(t.Field(i)) {
			elem := structVal.Field(i)
			if elem.Kind() == reflect.Struct {
				SetStructZero(elem)
			} else {
				elem.Set(reflect.Zero(elem.Type()))
			}
		}
	}
}

// Reset sets all elements of the object pointed to
// by resultRef to their default or zero values.
// But it works different from simply zeroing out everything,
// here are the exceptions:
// If resultRef is a pointer to a pointer, then
// the pointed to pointer will be reset to a new instance
// If resultRef is a pointer to a map, then the map
// will be reset to a new empty one.
// All other types pointed to by resultRef will be set
// to their default zero values.
func Reset(resultRef interface{}) {
	ptr := reflect.ValueOf(resultRef)
	if ptr.Kind() != reflect.Ptr {
		panic(fmt.Errorf("reflection.Reset(): resultRef must be a pointer, got %T", resultRef))
	}
	val := ptr.Elem()
	switch val.Kind() {
	case reflect.Ptr:
		// If resultRef is a pointer to a pointer,
		// set the pointer to a new instance
		// of the pointed to type
		val.Set(reflect.New(val.Type().Elem()))

	case reflect.Map:
		// If resultRef is a pointer to a map,
		// set make an empty new map
		val.Set(reflect.MakeChan(val.Type(), 0))

	case reflect.Struct:
		SetStructZero(val)

	default:
		val.Set(reflect.Zero(val.Type()))
	}
}

// SmartCopy copies struct or map fields from source
// to equally named struct or map fields of resultRef,
// by dereferencing source and resultRef if necessary
// to find a matching assignable type.
// All fields of the object referenced by resultPtr will be
// set to their default values before copying from source
// by calling Reset().
// SmartCopy is typically used for iterators with
// a method Next(resultRef interface{}) bool.
func SmartCopy(source, resultRef interface{}) {
	resultRefVal := reflect.ValueOf(resultRef)
	if resultRefVal.Kind() != reflect.Ptr && resultRefVal.Kind() != reflect.Map {
		panic(fmt.Errorf("reflection.SmartCopy(): resultRef must be a pointer or a map, got %T", resultRef))
	}
	Reset(resultRef)

	var resultVal reflect.Value
	if resultRefVal.Kind() == reflect.Map {
		resultVal = resultRefVal
	} else {
		resultVal = resultRefVal.Elem()
	}
	if resultVal.Kind() == reflect.Ptr {
		// Pointer to a pointer
		resultVal = resultVal.Elem()
	}

	sourceVal := reflect.ValueOf(source)
	if sourceVal.Kind() == reflect.Ptr {
		sourceVal = sourceVal.Elem()
	}
	if !smartCopyVals(sourceVal, resultVal) {
		panic(fmt.Errorf("reflection.SmartCopy(): Can't copy %T to %T", source, resultRef))
	}
}

func smartCopyVals(sourceVal, resultVal reflect.Value) bool {
	switch {
	case sourceVal.Type().AssignableTo(resultVal.Type()):
		resultVal.Set(sourceVal)

	case sourceVal.Kind() == reflect.Struct && resultVal.Kind() == reflect.Struct:
		CopyExportedStructFieldsVal(sourceVal, resultVal)

	case sourceVal.Kind() == reflect.Ptr && sourceVal.Elem().Type().AssignableTo(resultVal.Type()):
		smartCopyVals(sourceVal.Elem(), resultVal)

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

	case sourceVal.Kind() == reflect.Struct &&
		resultVal.Kind() == reflect.Map && resultVal.Type().Key().Kind() == reflect.String:

		sourceMap := ExportedStructFields(sourceVal)
		for key, src := range sourceMap {
			dst := resultVal.MapIndex(reflect.ValueOf(key))
			if dst.IsValid() && src.Type().AssignableTo(dst.Type()) {
				dst.Set(src)
			}
		}

	case sourceVal.Kind() == reflect.Map && sourceVal.Type().Key().Kind() == reflect.String &&
		resultVal.Kind() == reflect.Map && resultVal.Type().Key().Kind() == reflect.String &&
		sourceVal.Type().Elem().AssignableTo(resultVal.Type().Elem()):

		for _, key := range sourceVal.MapKeys() {
			if dst := resultVal.MapIndex(key); dst.IsValid() {
				dst.Set(sourceVal.MapIndex(key))
			}
		}

	default:
		return false
	}

	return true
}

func checkFunctionSignatureNums(t reflect.Type, args, results int) error {
	if t.Kind() != reflect.Func {
		return fmt.Errorf("Expected a function but got a %s", t)
	}
	if t.NumIn() != args {
		return fmt.Errorf("Expected %d function arguments, got %d", args, t.NumIn())
	}
	if t.NumOut() != results {
		return fmt.Errorf("Expected %d function results, got %d", results, t.NumOut())
	}
	return nil
}

func CheckFunctionSignature(f interface{}, args, results []reflect.Type) error {
	t := reflect.TypeOf(f)
	err := checkFunctionSignatureNums(t, len(args), len(results))
	if err != nil {
		return err
	}
	for i := range args {
		if args[i] != t.In(i) {
			return fmt.Errorf("Function argument %d must be %s, got %s", i, args[i], t.In(i))
		}
	}
	for i := range results {
		if results[i] != t.Out(i) {
			return fmt.Errorf("Function result %d must be %s, got %s", i, results[i], t.Out(i))
		}
	}
	return nil
}

func CheckFunctionSignatureKind(f interface{}, args, results []reflect.Kind) error {
	t := reflect.TypeOf(f)
	err := checkFunctionSignatureNums(t, len(args), len(results))
	if err != nil {
		return err
	}
	for i := range args {
		if args[i] != t.In(i).Kind() {
			return fmt.Errorf("Function argument %d must be %s kind, got %s kind", i, args[i], t.In(i).Kind())
		}
	}
	for i := range results {
		if results[i] != t.Out(i).Kind() {
			return fmt.Errorf("Function result %d must be %s kind, got %s kind", i, results[i], t.Out(i).Kind())
		}
	}
	return nil
}
