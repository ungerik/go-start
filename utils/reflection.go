package utils

import (
	"fmt"
	"github.com/ungerik/go-start/errs"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"

	"github.com/ungerik/go-start/debug"
)

type StructVisitor interface {
	BeginStruct(depth int, v reflect.Value)
	StructField(depth int, v reflect.Value, f reflect.StructField, index int)
	EndStruct(depth int, v reflect.Value)

	BeginSlice(depth int, v reflect.Value)
	SliceField(depth int, v reflect.Value, index int)
	EndSlice(depth int, v reflect.Value)

	BeginArray(depth int, v reflect.Value)
	ArrayField(depth int, v reflect.Value, index int)
	EndArray(depth int, v reflect.Value)
}

func NewStdLogStructVisitor() *LogStructVisitor {
	return &LogStructVisitor{Logger: log.New(os.Stdout, "", 0)}
}

// LogStructVisitor can be used for testing and debugging VisitStruct()
type LogStructVisitor struct {
	Logger *log.Logger
}

func (self *LogStructVisitor) BeginStruct(depth int, v reflect.Value) {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginStruct(%T)", indent, v.Interface())
}

func (self *LogStructVisitor) StructField(depth int, v reflect.Value, f reflect.StructField, index int) {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sStructField(%d, %s %s)", indent, index, f.Name, v.Type())
	default:
		self.Logger.Printf("%sStructField(%d, %s %s = %#v)", indent, index, f.Name, v.Type(), v.Interface())
	}
}

func (self *LogStructVisitor) EndStruct(depth int, v reflect.Value) {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndStruct(%T)", indent, v.Interface())
}

func (self *LogStructVisitor) BeginSlice(depth int, v reflect.Value) {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginSlice(%T)", indent, v.Interface())
}

func (self *LogStructVisitor) SliceField(depth int, v reflect.Value, index int) {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sSliceField(%d, %s)", indent, index, v.Type())
	default:
		self.Logger.Printf("%sSliceField(%d, %s = %#v)", indent, index, v.Type(), v.Interface())
	}
}

func (self *LogStructVisitor) EndSlice(depth int, v reflect.Value) {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndSlice(%T)", indent, v.Interface())
}

func (self *LogStructVisitor) BeginArray(depth int, v reflect.Value) {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginSlice(%T)", indent, v.Interface())
}

func (self *LogStructVisitor) ArrayField(depth int, v reflect.Value, index int) {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sArrayField(%d, %s)", indent, index, v.Type())
	default:
		self.Logger.Printf("%sArrayField(%d, %s = %#v)", indent, index, v.Type(), v.Interface())
	}

}

func (self *LogStructVisitor) EndArray(depth int, v reflect.Value) {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndSlice(%T)", indent, v.Interface())
}

/*
VisitStruct visits recursively all exported fields of a struct
and reports them via StructVisitor methods.
Pointers and interfaces are dereferenced silently until a non nil value
is found.
Structs that are embedded anonymously are inlined so that their fields
are reported as fields of the embedding struct at the same depth.
Anonymous struct fields that are not structs themselves are omitted.
*/
func VisitStruct(strct interface{}, visitor StructVisitor) {
	VisitStructDepth(strct, visitor, -1)
}

/*
VisitStructDepth is identical to VisitStruct except that its recursive
depth is limited to maxDepth with the first depth level being zero.
If maxDepth is -1, then the recursive depth is unlimited (VisitStruct).
*/
func VisitStructDepth(strct interface{}, visitor StructVisitor, maxDepth int) {
	visitStructRecursive(reflect.ValueOf(strct), visitor, maxDepth, 0)
}

func visitAnonymousStructFieldRecursive(visitor StructVisitor, v reflect.Value, index *int, depth int) {
	if v.Kind() == reflect.Struct {
		t := v.Type()
		n := t.NumField()
		for i := 0; i < n; i++ {
			f := t.Field(i)
			if f.PkgPath == "" { // Only exported fields
				if vi, ok := DereferenceValue(v.Field(i)); ok {
					if f.Anonymous {
						visitAnonymousStructFieldRecursive(visitor, vi, index, depth)
					} else {
						visitor.StructField(depth, vi, f, *index)
						*index++
					}
				}
			}
		}
	}
}

func visitStructRecursive(v reflect.Value, visitor StructVisitor, maxDepth, depth int) {
	if (maxDepth != -1 && depth > maxDepth) || !v.IsValid() {
		return
	}

	debug.Nop()
	// debug.Dump(v.Interface())

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			visitStructRecursive(v.Elem(), visitor, maxDepth, depth)
		}

	case reflect.Struct:
		visitor.BeginStruct(depth, v)
		t := v.Type()
		n := t.NumField()
		depth1 := depth + 1
		index := 0
		for i := 0; i < n; i++ {
			f := t.Field(i)
			if f.PkgPath == "" { // Only exported fields
				if vi, ok := DereferenceValue(v.Field(i)); ok {
					if f.Anonymous {
						visitAnonymousStructFieldRecursive(visitor, vi, &index, depth1)
					} else {
						visitor.StructField(depth1, vi, f, index)
						visitStructRecursive(vi, visitor, maxDepth, depth1)
						index++
					}
				}
			}
		}
		visitor.EndStruct(depth, v)

	case reflect.Slice:
		visitor.BeginSlice(depth, v)
		n := v.Len()
		depth1 := depth + 1
		for i := 0; i < n; i++ {
			if vi, ok := DereferenceValue(v.Index(i)); ok {
				visitor.SliceField(depth1, vi, i)
				visitStructRecursive(vi, visitor, maxDepth, depth1)
			}
		}
		visitor.EndSlice(depth, v)

	case reflect.Array:
		visitor.BeginArray(depth, v)
		n := v.Len()
		depth1 := depth + 1
		for i := 0; i < n; i++ {
			if vi, ok := DereferenceValue(v.Index(i)); ok {
				visitor.ArrayField(depth1, vi, i)
				visitStructRecursive(vi, visitor, maxDepth, depth1)
			}
		}
		visitor.EndArray(depth, v)
	}
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
