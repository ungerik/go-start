package model

import (
	"fmt"
	"github.com/ungerik/go-start/utils"
	// "github.com/ungerik/go-start/debug"
	"reflect"
)

type Visitor interface {
	BeginStruct(strct *MetaData) error
	StructField(field *MetaData) error
	EndStruct(strct *MetaData) error

	BeginSlice(slice *MetaData) error
	SliceField(field *MetaData) error
	EndSlice(slice *MetaData) error

	BeginArray(array *MetaData) error
	ArrayField(field *MetaData) error
	EndArray(array *MetaData) error
}

func Visit(model interface{}, visitor Visitor) error {
	return utils.VisitStruct(model, &structVisitorWrapper{visitor: visitor})
}

func VisitMaxDepth(model interface{}, maxDepth int, visitor Visitor) error {
	return utils.VisitStructDepth(model, &structVisitorWrapper{visitor: visitor}, maxDepth)
}

type structVisitorWrapper struct {
	visitor  Visitor
	metaData *MetaData
}

func (self *structVisitorWrapper) onBegin(depth int, v reflect.Value, kind MetaDataKind) {
	if depth == 0 {
		// no parent
		if self.metaData != nil {
			panic(fmt.Sprintf("Begin%s at depth 0 must not have a parent (self.metaData)", kind))
		}
		self.metaData = &MetaData{Value: v, Kind: kind}
	} else {
		if self.metaData.Depth != depth {
			panic(fmt.Sprintf("Begin%s: If not the root, there must be some self.MetaData from StructField before", kind))
		}
	}
}

func (self *structVisitorWrapper) onEnd(depth int, kind MetaDataKind) {
	if depth != self.metaData.Depth {
		if depth+1 != self.metaData.Depth {
			panic(fmt.Sprintf("End%s: self.metaData.Depth (%d) must be depth or depth+1 (%d or %d)", kind, self.metaData.Depth, depth, depth+1))
		}
		self.metaData = self.metaData.Parent
	}
	if self.metaData.Kind != kind {
		panic(fmt.Sprintf("End%s called for %s", kind, self.metaData.Kind))
	}
}

func (self *structVisitorWrapper) onArrayOrSliceField(depth int, v reflect.Value, index int, parentKind MetaDataKind) {
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
	self.metaData = &MetaData{
		Value:  v,
		Kind:   GetMetaDataKind(v),
		Parent: parent,
		Depth:  depth,
		Index:  index,
	}
}

func (self *structVisitorWrapper) BeginStruct(depth int, v reflect.Value) error {
	self.onBegin(depth, v, StructKind)
	return self.visitor.BeginStruct(self.metaData)
}

func (self *structVisitorWrapper) StructField(depth int, v reflect.Value, f reflect.StructField, index int) error {
	var parent *MetaData
	if index == 0 {
		// first field of struct
		if depth != self.metaData.Depth+1 {
			panic("Depth of first field of a struct must be its parent struct's depth plus one")
		}
		// no previous sibling available, parent ist current self.metaData
		parent = self.metaData
	} else {
		if depth != self.metaData.Depth {
			panic("If not the first field of a struct, there must already be MetaData of the same depth from the previous sibling")
		}
		// parent is the same of previous sibling which is current self.metaData
		parent = self.metaData.Parent
	}
	if parent.Kind != StructKind {
		panic(fmt.Sprintf("StructField called for %s parent", parent.Kind))
	}
	self.metaData = &MetaData{
		Value:  v,
		Kind:   GetMetaDataKind(v),
		Parent: parent,
		Depth:  depth,
		Name:   f.Name,
		Index:  index,
		tag:    f.Tag,
	}
	return self.visitor.StructField(self.metaData)
}

func (self *structVisitorWrapper) EndStruct(depth int, v reflect.Value) error {
	self.onEnd(depth, StructKind)
	return self.visitor.EndStruct(self.metaData)
}

func (self *structVisitorWrapper) BeginSlice(depth int, v reflect.Value) error {
	self.onBegin(depth, v, SliceKind)
	return self.visitor.BeginSlice(self.metaData)
}

func (self *structVisitorWrapper) SliceField(depth int, v reflect.Value, index int) error {
	self.onArrayOrSliceField(depth, v, index, SliceKind)
	return self.visitor.SliceField(self.metaData)
}

func (self *structVisitorWrapper) EndSlice(depth int, v reflect.Value) error {
	self.onEnd(depth, SliceKind)
	return self.visitor.EndSlice(self.metaData)
}

func (self *structVisitorWrapper) BeginArray(depth int, v reflect.Value) error {
	self.onBegin(depth, v, ArrayKind)
	return self.visitor.BeginArray(self.metaData)
}

func (self *structVisitorWrapper) ArrayField(depth int, v reflect.Value, index int) error {
	self.onArrayOrSliceField(depth, v, index, ArrayKind)
	return self.visitor.ArrayField(self.metaData)
}

func (self *structVisitorWrapper) EndArray(depth int, v reflect.Value) error {
	self.onEnd(depth, ArrayKind)
	return self.visitor.EndArray(self.metaData)
}

// FieldOnlyVisitor calls its function for every struct, array and slice field.
type FieldOnlyVisitor func(field *MetaData) error

func (self FieldOnlyVisitor) BeginStruct(strct *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) StructField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndStruct(strct *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) BeginSlice(slice *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) SliceField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndSlice(slice *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) BeginArray(array *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) ArrayField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndArray(array *MetaData) error {
	return nil
}

// VisitorFunc calls its function for every Visitor method call,
// thus mapping all Visitor methods on a single function.
type VisitorFunc func(data *MetaData) error

func (self VisitorFunc) BeginStruct(strct *MetaData) error {
	return self(strct)
}

func (self VisitorFunc) StructField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndStruct(strct *MetaData) error {
	return self(strct)
}

func (self VisitorFunc) BeginSlice(slice *MetaData) error {
	return self(slice)
}

func (self VisitorFunc) SliceField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndSlice(slice *MetaData) error {
	return self(slice)
}

func (self VisitorFunc) BeginArray(array *MetaData) error {
	return self(array)
}

func (self VisitorFunc) ArrayField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndArray(array *MetaData) error {
	return self(array)
}
