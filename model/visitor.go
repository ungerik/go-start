package model

import (
	"fmt"
	"reflect"

	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/reflection"
)

type Visitor interface {
	BeginNamedFields(namedFields *MetaData) error
	NamedField(field *MetaData) error
	EndNamedFields(namedFields *MetaData) error

	BeginIndexedFields(indexedFields *MetaData) error
	IndexedField(field *MetaData) error
	EndIndexedFields(indexedFields *MetaData) error
}

func Visit(model interface{}, visitor Visitor) error {
	debug.Nop()
	return reflection.VisitStruct(model, &structVisitorWrapper{visitor: visitor})
}

func VisitMaxDepth(model interface{}, maxDepth int, visitor Visitor) error {
	return reflection.VisitStructDepth(model, &structVisitorWrapper{visitor: visitor}, maxDepth)
}

///////////////////////////////////////////////////////////////////////////////

// todo fix special case of empty structs
type structVisitorWrapper struct {
	visitor  Visitor
	metaData *MetaData
}

func (self *structVisitorWrapper) begin(depth int, v reflect.Value, kind MetaDataKind) (metaData *MetaData) {
	// debug.Print("begin")
	if depth == 0 {
		// no parent
		if self.metaData != nil {
			panic(fmt.Sprintf("Begin%s at depth 0 must not have a parent (self.metaData)", kind))
		}
		return &MetaData{Value: v, Kind: kind}
	}

	if self.metaData.Depth != depth {
		panic(fmt.Sprintf("Begin%s: If not the root, there must be some self.MetaData from NamedFields before", kind))
	}
	return self.metaData
}

func (self *structVisitorWrapper) end(depth int, kind MetaDataKind) (metaData *MetaData) {
	// debug.Print("end")
	if depth == self.metaData.Depth {
		metaData = self.metaData
	} else {
		// Non empty parent, current self.metaData.Depth must be depth+1
		if depth+1 != self.metaData.Depth {
			panic(fmt.Sprintf("End%s: self.metaData.Depth (%d) must be depth or depth+1 (%d or %d)", kind, self.metaData.Depth, depth, depth+1))
		}
		metaData = self.metaData.Parent
	}
	if metaData.Kind != kind {
		panic(fmt.Sprintf("End%s called for %s", kind, metaData.Kind))
	}
	return metaData
}

func (self *structVisitorWrapper) indexedFieldMetaData(depth int, v reflect.Value, index int, parentKind MetaDataKind) *MetaData {
	// debug.Print("onArrayOrSliceField")
	var parent *MetaData
	if index == 0 {
		// first field of array or struct
		if depth != self.metaData.Depth+1 {
			panic(fmt.Sprintf("Depth of first field of a %s must be its parent %s's depth plus one", parentKind, parentKind))
		}
		// no previous sibling available, parent ist current self.metaData
		parent = self.metaData
	} else {
		if depth != self.metaData.Depth {
			panic(fmt.Sprintf("If not the first field of a %s, there must already be MetaData of the same depth from the previous sibling", parentKind))
		}
		// parent is the same of previous sibling which is current self.metaData
		parent = self.metaData.Parent
	}
	if parent.Kind != parentKind {
		panic(fmt.Sprintf("%sField called for %s parent", parentKind, parent.Kind))
	}
	return &MetaData{
		Value:  v,
		Kind:   GetMetaDataKind(v),
		Parent: parent,
		Depth:  depth,
		Index:  index,
	}
}

func (self *structVisitorWrapper) namedFieldMetaData(depth int, v reflect.Value, name string, index int) *MetaData {
	// debug.Print("onNamedField")
	var parent *MetaData
	if index == 0 {
		// first field of struct
		if depth != self.metaData.Depth+1 {
			panic("Depth of first field of NamedFields must be its parent NamedFields' depth plus one")
		}
		// no previous sibling available, parent ist current self.metaData
		parent = self.metaData
	} else {
		if depth != self.metaData.Depth {
			panic("If not the first field of NamedFields, there must already be MetaData of the same depth from the previous sibling")
		}
		// parent is the same of previous sibling which is current self.metaData
		parent = self.metaData.Parent
	}
	if !parent.Kind.HasNamedFields() {
		panic(fmt.Sprintf("namedFieldMetaData called for %s parent", parent.Kind))
	}
	return &MetaData{
		Value:  v,
		Kind:   GetMetaDataKind(v),
		Parent: parent,
		Depth:  depth,
		Name:   name,
		Index:  index,
	}
}

func (self *structVisitorWrapper) BeginStruct(depth int, v reflect.Value) error {
	self.metaData = self.begin(depth, v, StructKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore struct fields of Value implementations
		return nil
	}
	return self.visitor.BeginNamedFields(self.metaData)
}

func (self *structVisitorWrapper) StructField(depth int, v reflect.Value, f reflect.StructField, index int) error {
	self.metaData = self.namedFieldMetaData(depth, v, f.Name, index)
	self.metaData.tag = f.Tag
	if self.metaData.Parent.IsModelValueOrChild() {
		// Ignore struct fields of Value implementations
		return nil
	}
	return self.visitor.NamedField(self.metaData)
}

func (self *structVisitorWrapper) EndStruct(depth int, v reflect.Value) error {
	self.metaData = self.end(depth, StructKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore struct fields of Value implementations
		return nil
	}
	return self.visitor.EndNamedFields(self.metaData)
}

func (self *structVisitorWrapper) BeginMap(depth int, v reflect.Value) error {
	self.metaData = self.begin(depth, v, MapKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore map fields of Value implementations
		return nil
	}
	return self.visitor.BeginNamedFields(self.metaData)
}

func (self *structVisitorWrapper) MapField(depth int, v reflect.Value, key string, index int) error {
	self.metaData = self.namedFieldMetaData(depth, v, key, index)
	if self.metaData.Parent.IsModelValueOrChild() {
		// Ignore map fields of Value implementations
		return nil
	}
	return self.visitor.NamedField(self.metaData)
}

func (self *structVisitorWrapper) EndMap(depth int, v reflect.Value) error {
	self.metaData = self.end(depth, MapKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore map fields of Value implementations
		return nil
	}
	return self.visitor.EndNamedFields(self.metaData)
}

func (self *structVisitorWrapper) BeginSlice(depth int, v reflect.Value) error {
	if v.Type().Elem() == DynamicValueType {
		// A slice of DynamicValues is treated as NamedFields
		self.metaData = self.begin(depth, v, DynamicKind)
		return self.visitor.BeginNamedFields(self.metaData)
	}
	self.metaData = self.begin(depth, v, SliceKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore slice fields of Value implementations
		return nil
	}
	return self.visitor.BeginIndexedFields(self.metaData)
}

func (self *structVisitorWrapper) SliceField(depth int, v reflect.Value, index int) error {
	if dynamicValue, ok := v.Interface().(DynamicValue); ok {
		// A slice of DynamicValues is treated as NamedFields
		self.metaData = self.namedFieldMetaData(depth, reflect.ValueOf(dynamicValue.Value).Elem(), dynamicValue.Name, index)
		self.metaData.attribs = dynamicValue.Attribs
		return self.visitor.NamedField(self.metaData)
	}
	self.metaData = self.indexedFieldMetaData(depth, v, index, SliceKind)
	if self.metaData.Parent.IsModelValueOrChild() {
		// Ignore slice fields of Value implementations
		return nil
	}
	return self.visitor.IndexedField(self.metaData)
}

func (self *structVisitorWrapper) EndSlice(depth int, v reflect.Value) error {
	if v.Type().Elem() == DynamicValueType {
		// A slice of DynamicValues is treated as NamedFields
		self.metaData = self.end(depth, DynamicKind)
		return self.visitor.EndNamedFields(self.metaData)
	}
	self.metaData = self.end(depth, SliceKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore slice fields of Value implementations
		return nil
	}
	return self.visitor.EndIndexedFields(self.metaData)
}

func (self *structVisitorWrapper) BeginArray(depth int, v reflect.Value) error {
	if v.Type().Elem() == DynamicValueType {
		self.metaData = self.begin(depth, v, DynamicKind)
		return self.visitor.BeginNamedFields(self.metaData)
	}
	self.metaData = self.begin(depth, v, ArrayKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore array fields of Value implementations
		return nil
	}
	return self.visitor.BeginIndexedFields(self.metaData)
}

func (self *structVisitorWrapper) ArrayField(depth int, v reflect.Value, index int) error {
	if dynamicValue, ok := v.Interface().(DynamicValue); ok {
		self.metaData = self.namedFieldMetaData(depth, reflect.ValueOf(dynamicValue.Value), dynamicValue.Name, index)
		self.metaData.attribs = dynamicValue.Attribs
		return self.visitor.NamedField(self.metaData)
	}
	self.metaData = self.indexedFieldMetaData(depth, v, index, ArrayKind)
	if self.metaData.Parent.IsModelValueOrChild() {
		// Ignore array fields of Value implementations
		return nil
	}
	return self.visitor.IndexedField(self.metaData)
}

func (self *structVisitorWrapper) EndArray(depth int, v reflect.Value) error {
	if v.Type().Elem() == DynamicValueType {
		self.metaData = self.end(depth, DynamicKind)
		return self.visitor.EndNamedFields(self.metaData)
	}
	self.metaData = self.end(depth, ArrayKind)
	if self.metaData.IsModelValueOrChild() {
		// Ignore array fields of Value implementations
		return nil
	}
	return self.visitor.EndIndexedFields(self.metaData)
}
