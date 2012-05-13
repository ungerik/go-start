package utils

import (
	"log"
	"os"
	"reflect"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// StructVisitor

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

/*
VisitStruct visits recursively all exported fields of a struct
and reports them via StructVisitor methods.
Pointers and interfaces are dereferenced silently until a non nil value
is found.
Structs that are embedded anonymously are inlined so that their fields
are reported as fields of the embedding struct at the same depth.
Anonymous struct fields that are not structs themselves are omitted.
Struct fields with the tag gostart:"-" are ignored.
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
			if f.PkgPath == "" && f.Tag.Get("gostart") != "-" { // Only exported fields
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

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			visitStructRecursive(v.Elem(), visitor, maxDepth, depth)
		}

	case reflect.Struct:
		visitor.BeginStruct(depth, v)
		depth1 := depth + 1
		if maxDepth == -1 || depth1 <= maxDepth {
			t := v.Type()
			n := t.NumField()
			index := 0
			for i := 0; i < n; i++ {
				f := t.Field(i)
				if f.PkgPath == "" && f.Tag.Get("gostart") != "-" { // Only exported fields
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
		}
		visitor.EndStruct(depth, v)

	case reflect.Slice:
		visitor.BeginSlice(depth, v)
		depth1 := depth + 1
		if maxDepth == -1 || depth1 <= maxDepth {
			n := v.Len()
			for i := 0; i < n; i++ {
				if vi, ok := DereferenceValue(v.Index(i)); ok {
					visitor.SliceField(depth1, vi, i)
					visitStructRecursive(vi, visitor, maxDepth, depth1)
				}
			}
		}
		visitor.EndSlice(depth, v)

	case reflect.Array:
		visitor.BeginArray(depth, v)
		depth1 := depth + 1
		if maxDepth == -1 || depth1 <= maxDepth {
			n := v.Len()
			for i := 0; i < n; i++ {
				if vi, ok := DereferenceValue(v.Index(i)); ok {
					visitor.ArrayField(depth1, vi, i)
					visitStructRecursive(vi, visitor, maxDepth, depth1)
				}
			}
		}
		visitor.EndArray(depth, v)
	}
}

////////////////////////////////////////////////////////////////////////////////
// LogStructVisitor

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
	self.Logger.Printf("%BeginArray(%T)", indent, v.Interface())
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
	self.Logger.Printf("%EndArray(%T)", indent, v.Interface())
}
