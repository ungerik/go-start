package model

import (
	"fmt"
	"reflect"
)

type Visitor interface {
	BeginStruct(strct *MetaData)
	StructField(field *MetaData)
	EndStruct(strct *MetaData)

	BeginSlice(slice *MetaData)
	SliceField(field *MetaData)
	EndSlice(slice *MetaData)

	BeginArray(array *MetaData)
	ArrayField(field *MetaData)
	EndArray(array *MetaData)
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
		if depth != self.metaData.Depth+1 {
			panic(fmt.Sprintf("End%s: self.metaData.Depth(%d) can only be depth(%d) or depth + 1(%d)", kind, self.metaData.Depth, depth, depth+1))
		}
		self.metaData = self.metaData.Parent
	}
	if self.metaData.Kind != kind {
		panic(fmt.Sprintf("End%s called for %s", kind, self.metaData.Kind))
	}
}

func (self *structVisitorWrapper) onArrayOrSliceField(depth int, v reflect.Value, index int, parentKind MetaDataKind) {
	if index == 0 {
		// first field of array or struct
		if depth != self.metaData.Depth+1 {
			panic(fmt.Sprintf("Depth of first field of a %s must be its parent %s's depth plus one", parentKind, parentKind))
		}
		// create MetaData for this depth
		self.metaData = &MetaData{
			Value:  v,
			Kind:   GetMetaDataKind(v),
			Parent: self.metaData,
			Depth:  depth,
			Index:  index,
		}
	} else {
		if depth != self.metaData.Depth {
			panic(fmt.Sprintf("If not the first field of a %s, there must already be MetaData of the same depth from the previous sibling", parentKind))
		}
		// only have to change what's different for this field
		self.metaData.Value = v
		self.metaData.Index = index
	}
	if self.metaData.Parent.Kind != StructKind {
		panic(fmt.Sprintf("%sField called for %s parent", parentKind, self.metaData.Parent.Kind))
	}
}

func (self *structVisitorWrapper) BeginStruct(depth int, v reflect.Value) {
	self.onBegin(depth, v, StructKind)
	self.visitor.BeginStruct(self.metaData)
}

func (self *structVisitorWrapper) StructField(depth int, v reflect.Value, f reflect.StructField, index int) {
	if index == 0 {
		// first field of struct
		if depth != self.metaData.Depth+1 {
			panic("Depth of first field of a struct must be its parent struct's depth plus one")
		}
		// create MetaData for this depth
		self.metaData = &MetaData{
			Value:  v,
			Kind:   GetMetaDataKind(v),
			Parent: self.metaData,
			Depth:  depth,
			Name:   f.Name,
			Index:  index,
			tag:    f.Tag.Get("gostart"),
		}
	} else {
		if depth != self.metaData.Depth {
			panic("If not the first field of a struct, there must already be MetaData of the same depth from the previous sibling")
		}
		// only have to change what's different for this field
		self.metaData.Value = v
		self.metaData.Kind = GetMetaDataKind(v)
		self.metaData.Name = f.Name
		self.metaData.Index = index
		self.metaData.tag = f.Tag.Get("gostart")
	}
	if self.metaData.Parent.Kind != StructKind {
		panic(fmt.Sprintf("StructField called for %s parent", self.metaData.Parent.Kind))
	}
	self.visitor.StructField(self.metaData)
}

func (self *structVisitorWrapper) EndStruct(depth int, v reflect.Value) {
	self.onEnd(depth, StructKind)
	self.visitor.EndStruct(self.metaData)
}

func (self *structVisitorWrapper) BeginSlice(depth int, v reflect.Value) {
	self.onBegin(depth, v, SliceKind)
	self.visitor.BeginSlice(self.metaData)
}

func (self *structVisitorWrapper) SliceField(depth int, v reflect.Value, index int) {
	self.onArrayOrSliceField(depth, v, index, SliceKind)
	self.visitor.SliceField(self.metaData)
}

func (self *structVisitorWrapper) EndSlice(depth int, v reflect.Value) {
	self.onEnd(depth, SliceKind)
	self.visitor.EndSlice(self.metaData)
}

func (self *structVisitorWrapper) BeginArray(depth int, v reflect.Value) {
	self.onBegin(depth, v, ArrayKind)
	self.visitor.BeginArray(self.metaData)
}

func (self *structVisitorWrapper) ArrayField(depth int, v reflect.Value, index int) {
	self.onArrayOrSliceField(depth, v, index, ArrayKind)
	self.visitor.ArrayField(self.metaData)
}

func (self *structVisitorWrapper) EndArray(depth int, v reflect.Value) {
	self.onEnd(depth, ArrayKind)
	self.visitor.EndArray(self.metaData)
}
