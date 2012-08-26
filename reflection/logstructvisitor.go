package reflection

import (
	"log"
	"os"
	"reflect"
	"strings"
	// "github.com/ungerik/go-start/debug"
)

func NewStdLogStructVisitor() *LogStructVisitor {
	return &LogStructVisitor{Logger: log.New(os.Stdout, "", 0)}
}

// LogStructVisitor can be used for testing and debugging VisitStruct()
type LogStructVisitor struct {
	Logger *log.Logger
}

func (self *LogStructVisitor) BeginStruct(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginStruct(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) StructField(depth int, v reflect.Value, f reflect.StructField, index int) error {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sStructField(%d: %s %s)", indent, index, f.Name, v.Type())
	default:
		self.Logger.Printf("%sStructField(%d: %s %s = %#v)", indent, index, f.Name, v.Type(), v.Interface())
	}
	return nil
}

func (self *LogStructVisitor) EndStruct(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndStruct(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) BeginSlice(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginSlice(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) SliceField(depth int, v reflect.Value, index int) error {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sSliceField(%d: %s)", indent, index, v.Type())
	default:
		self.Logger.Printf("%sSliceField(%d: %s = %#v)", indent, index, v.Type(), v.Interface())
	}
	return nil
}

func (self *LogStructVisitor) EndSlice(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndSlice(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) BeginArray(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginArray(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) ArrayField(depth int, v reflect.Value, index int) error {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sArrayField(%d: %s)", indent, index, v.Type())
	default:
		self.Logger.Printf("%sArrayField(%d: %s = %#v)", indent, index, v.Type(), v.Interface())
	}
	return nil
}

func (self *LogStructVisitor) EndArray(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndArray(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) BeginMap(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sBeginMap(%T)", indent, v.Interface())
	return nil
}

func (self *LogStructVisitor) MapField(depth int, v reflect.Value, key string, index int) error {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array:
		self.Logger.Printf("%sMapField(\"%s\": %s)", indent, key, v.Type())
	default:
		self.Logger.Printf("%sMapField(\"%s\": %s = %#v)", indent, key, v.Type(), v.Interface())
	}
	return nil
}

func (self *LogStructVisitor) EndMap(depth int, v reflect.Value) error {
	indent := strings.Repeat("  ", depth)
	self.Logger.Printf("%sEndMap(%T)", indent, v.Interface())
	return nil
}
