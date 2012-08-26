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

	BeginSlice(slice *MetaData) error
	SliceField(field *MetaData) error
	EndSlice(slice *MetaData) error

	BeginArray(array *MetaData) error
	ArrayField(field *MetaData) error
	EndArray(array *MetaData) error
}

func Visit(model interface{}, visitor Visitor) error {
	return reflection.VisitStruct(model, &structVisitorWrapper{visitor: visitor})
}

func VisitMaxDepth(model interface{}, maxDepth int, visitor Visitor) error {
	return reflection.VisitStructDepth(model, &structVisitorWrapper{visitor: visitor}, maxDepth)
}

///////////////////////////////////////////////////////////////////////////////

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

func (self *structVisitorWrapper) arrayOrSliceFieldMetaData(depth int, v reflect.Value, index int, parentKind MetaDataKind) *MetaData {
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
	if parent.Kind != NamedFieldsKind {
		panic(fmt.Sprintf("StructField/MapField called for %s parent", parent.Kind))
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
	self.metaData = self.begin(depth, v, NamedFieldsKind)
	return self.visitor.BeginNamedFields(self.metaData)
}

func (self *structVisitorWrapper) StructField(depth int, v reflect.Value, f reflect.StructField, index int) error {
	self.metaData = self.namedFieldMetaData(depth, v, f.Name, index)
	self.metaData.tag = f.Tag
	return self.visitor.NamedField(self.metaData)
}

func (self *structVisitorWrapper) EndStruct(depth int, v reflect.Value) error {
	self.metaData = self.end(depth, NamedFieldsKind)
	return self.visitor.EndNamedFields(self.metaData)
}

func (self *structVisitorWrapper) BeginMap(depth int, v reflect.Value) error {
	debug.Print("BeginMap")
	self.metaData = self.begin(depth, v, NamedFieldsKind)
	return self.visitor.BeginNamedFields(self.metaData)
}

func (self *structVisitorWrapper) MapField(depth int, v reflect.Value, key string, index int) error {
	debug.Print("MapField")
	self.metaData = self.namedFieldMetaData(depth, v, key, index)
	return self.visitor.NamedField(self.metaData)
}

func (self *structVisitorWrapper) EndMap(depth int, v reflect.Value) error {
	debug.Print("EndMap")
	self.metaData = self.end(depth, NamedFieldsKind)
	return self.visitor.EndNamedFields(self.metaData)
}

func (self *structVisitorWrapper) BeginSlice(depth int, v reflect.Value) error {
	// if v.Type().Elem() == typeOfDynamicValue {

	// }
	self.metaData = self.begin(depth, v, SliceKind)
	return self.visitor.BeginSlice(self.metaData)
}

func (self *structVisitorWrapper) SliceField(depth int, v reflect.Value, index int) error {
	self.metaData = self.arrayOrSliceFieldMetaData(depth, v, index, SliceKind)
	return self.visitor.SliceField(self.metaData)
}

func (self *structVisitorWrapper) EndSlice(depth int, v reflect.Value) error {
	self.metaData = self.end(depth, SliceKind)
	return self.visitor.EndSlice(self.metaData)
}

func (self *structVisitorWrapper) BeginArray(depth int, v reflect.Value) error {
	self.metaData = self.begin(depth, v, ArrayKind)
	return self.visitor.BeginArray(self.metaData)
}

func (self *structVisitorWrapper) ArrayField(depth int, v reflect.Value, index int) error {
	self.metaData = self.arrayOrSliceFieldMetaData(depth, v, index, ArrayKind)
	return self.visitor.ArrayField(self.metaData)
}

func (self *structVisitorWrapper) EndArray(depth int, v reflect.Value) error {
	self.metaData = self.end(depth, ArrayKind)
	return self.visitor.EndArray(self.metaData)
}
