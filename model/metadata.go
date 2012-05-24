package model

import (
	// "github.com/ungerik/go-start/debug"
	"reflect"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// MetaData

// MetaDataKind is the kind of a data item described by MetaData.
type MetaDataKind int

func (self MetaDataKind) String() string {
	switch self {
	case StructKind:
		return "Struct"
	case ArrayKind:
		return "Array"
	case SliceKind:
		return "Slice"
	case ValueKind:
		return "Value"
	}
	panic("Unknown MetaDataKind!")
}

const (
	StructKind MetaDataKind = iota
	ArrayKind
	SliceKind
	ValueKind // everything else, no children
)

func GetMetaDataKind(v reflect.Value) MetaDataKind {
	switch v.Kind() {
	case reflect.Struct:
		return StructKind
	case reflect.Array:
		return ArrayKind
	case reflect.Slice:
		return SliceKind
	}
	return ValueKind
}

// MetaData holds meta data about an model data item.
type MetaData struct {
	Value   reflect.Value
	Kind    MetaDataKind
	Parent  *MetaData
	Depth   int    // number of steps up to the root parent
	Name    string // empty for array and slice fields
	Index   int    // will also be set for struct fields
	tag     string
	attribs map[string]string // cached tag attributes
}

// ParentKind returns the parent's MetaDataKind.
// It returns StructKind if Parent is nil because the root parent will
// always be a struct.
func (self *MetaData) ParentKind() MetaDataKind {
	if self.Parent == nil {
		return StructKind
	}
	return self.Parent.Kind
}

func (self *MetaData) IsArrayOrSliceField() bool {
	return self.Kind == ValueKind && self.Name == ""
}

func (self *MetaData) IsStructField() bool {
	return self.Name != ""
}

// // GetDepth returns self.Depth or -1 if self is nil.
// func (self *MetaData) GetDepth() int {
// 	if self == nil {
// 		return -1
// 	}
// 	return self.Depth
// }

/*
Attrib returns the value of a tag attribute if available.
Array and slice fields inherit the attributes of their named
parent fields.
The meaning of attributes is interpreted by the package that reads them.
Attributes are defined in a struct tag named "gostart" and written
as name=value. Multiple attributes are separated by '|'.
Example:

	type Struct {
		X int `gostart:"min=0|max=10"`
		S []int `gostart:"maxlen=3|min=0|max=10"
		hidden int
		Ignore int `gostart:"-"`
	}
*/
func (self *MetaData) Attrib(name string) (value string, ok bool) {
	m := self
	for !m.IsStructField() {
		m = m.Parent
		if m == nil {
			return "", false
		}
	}
	if m.attribs == nil {
		m.attribs = ParseTagAttribs(m.tag)
	}
	value, ok = m.attribs[name]
	return value, ok
}

// BoolAttrib uses Attrib() to check if the value of an attribute is "true".
// A non existing attribibute is considered to be false.
func (self *MetaData) BoolAttrib(name string) bool {
	value, ok := self.Attrib(name)
	return ok && value == "true"
}

// Selector returns the field names or indices for array/slice fields
// from the root parent up the the current data item concatenated with with '.'.
func (self *MetaData) Selector() string {
	names := make([]string, self.Depth-1)
	for i, m := self.Depth-2, self; i >= 0; i-- {
		if m.IsStructField() {
			names[i] = m.Name
		} else {
			names[i] = strconv.Itoa(m.Index)
		}
		m = m.Parent
	}
	return strings.Join(names, ".")
}

// WildcardSelector returns the field names or the wildcard '$' for array/slice
// fields from the root parent up the the current data item
// concatenated with with '.'.
func (self *MetaData) WildcardSelector() string {
	names := make([]string, self.Depth-1)
	for i, m := self.Depth-2, self; i >= 0; i-- {
		if m.IsStructField() {
			names[i] = m.Name
		} else {
			names[i] = "$"
		}
		m = m.Parent
	}
	return strings.Join(names, ".")
}

func ParseTagAttribs(tag string) map[string]string {
	attribs := make(map[string]string)
	for _, s := range strings.Split(tag, "|") {
		pos := strings.Index(s, "=")
		if pos == -1 {
			attribs[s] = "true"
		} else {
			attribs[s[:pos]] = s[pos+1:]
		}
	}
	return attribs
}
