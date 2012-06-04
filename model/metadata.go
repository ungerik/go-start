package model

import (
	// "github.com/ungerik/go-start/debug"
	"bytes"
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
	Kind    MetaDataKind
	Value   reflect.Value
	Depth   int    // number of steps up to the root parent
	Name    string // empty for array and slice fields
	Index   int    // will also be set for struct fields
	Parent  *MetaData
	tag     string
	attribs map[string]string // cached tag attributes
	path    []*MetaData
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

// NameOrIndex returns self.Name if not empty or else self.Index.
func (self *MetaData) NameOrIndex() string {
	if self.Name != "" {
		return self.Name
	}
	return strconv.Itoa(self.Index)
}

// NameOrWildcard returns self.Name if not empty or else the wildcard "$".
func (self *MetaData) NameOrWildcard() string {
	if self.Name != "" {
		return self.Name
	}
	return strconv.Itoa(self.Index)
}

// Path returns a slice of *MetaData that holds all parents from
// the root parent up to (and including) self.
func (self *MetaData) Path() []*MetaData {
	if self.path == nil {
		self.path = make([]*MetaData, self.Depth+1)
		for i, m := self.Depth, self; i >= 0; i-- {
			self.path[i] = m
			m = m.Parent
		}
	}
	return self.path
}

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
	if self.attribs == nil {
		structField := self
		for !structField.IsStructField() {
			structField = structField.Parent
			if structField == nil {
				return "", false
			}
		}
		self.attribs = ParseTagAttribs(structField.tag)
	}
	value, ok = self.attribs[name]
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
	var buf bytes.Buffer
	for _, m := range self.Path()[1:] {
		if buf.Len() > 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(m.NameOrIndex())
	}
	return buf.String()
}

// WildcardSelector returns the field names or the wildcard '$' for array/slice
// fields from the root parent up the the current data item
// concatenated with with '.'.
func (self *MetaData) WildcardSelector() string {
	var buf bytes.Buffer
	for _, m := range self.Path()[1:] {
		if buf.Len() > 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(m.NameOrWildcard())
	}
	return buf.String()
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
