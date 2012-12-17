package model

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// MetaData

// MetaDataKind is the kind of a data item described by MetaData.
type MetaDataKind int

func (self MetaDataKind) HasNamedFields() bool {
	return self == StructKind || self == MapKind || self == DynamicKind
}

func (self MetaDataKind) HasIndexedFields() bool {
	return self == ArrayKind || self == SliceKind
}

func (self MetaDataKind) String() string {
	switch self {
	case StructKind:
		return "Struct"
	case MapKind:
		return "Map"
	case DynamicKind:
		return "Dynamic"
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
	MapKind
	DynamicKind
	ArrayKind
	SliceKind
	ValueKind // everything else, no children
)

func GetMetaDataKind(v reflect.Value) MetaDataKind {
	switch v.Kind() {
	case reflect.Struct:
		return StructKind

	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String {
			return MapKind
		}

	case reflect.Array:
		if v.Type().Elem() == DynamicValueType {
			return DynamicKind
		}
		return ArrayKind

	case reflect.Slice:
		if v.Type().Elem() == DynamicValueType {
			return DynamicKind
		}
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
	Dynamic bool
	tag     reflect.StructTag
	// cached tag attributes, by package
	attribs map[string]map[string]string
	path    []*MetaData
}

func (self *MetaData) ModelValue() (val Value, ok bool) {
	if self == nil {
		return nil, false
	}
	val, ok = self.Value.Addr().Interface().(Value)
	return val, ok
}

func (self *MetaData) IsModelValue() bool {
	if self == nil {
		return false
	}
	_, ok := self.ModelValue()
	return ok
}

func (self *MetaData) IsModelValueOrChild() bool {
	if self == nil {
		return false
	}
	if self.IsModelValue() {
		return true
	}
	return self.Parent.IsModelValueOrChild()
}

func (self *MetaData) ModelValidator() (val Validator, ok bool) {
	val, ok = self.Value.Interface().(Validator)
	if !ok && self.Value.CanAddr() {
		val, ok = self.Value.Addr().Interface().(Validator)
	}
	return val, ok
}

func (self *MetaData) RootParent() *MetaData {
	root := self
	for root.Parent != nil {
		root = root.Parent
	}
	return root
}

// // ParentKind returns the parent's MetaDataKind.
// // It returns StructKind if Parent is nil because the root parent will
// // always be a struct.
// func (self *MetaData) ParentKind() MetaDataKind {
// 	if self.Parent == nil {
// 		return NamedFieldsKind
// 	}
// 	return self.Parent.Kind
// }

func (self *MetaData) IsSelfOrParentIndexed() bool {
	if self.Name == "" {
		return true
	}
	if self.Parent != nil {
		return self.Parent.IsSelfOrParentIndexed()
	}
	return false
}

func (self *MetaData) IsIndexedValue() bool {
	return self.Kind == ValueKind && self.Name == ""
}

func (self *MetaData) IsNamedValue() bool {
	return self.Kind == ValueKind && self.Name != ""
}

func (self *MetaData) IsNamed() bool {
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
	return "$"
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

func ParseTagAttribs(tag string) map[string]string {
	attribs := make(map[string]string)
	for _, s := range strings.Split(tag, "|") {
		pos := strings.IndexRune(s, '=')
		if pos == -1 {
			attribs[s] = "true"
		} else {
			attribs[s[:pos]] = s[pos+1:]
		}
	}
	return attribs
}

/*
Attrib returns the value of a tag attribute if available.
Array and slice fields inherit the attributes of their named
parent fields.
The meaning of attributes is interpreted by the package that reads them.
Attributes are defined in a struct tag named "gostart", "model" or "view"
and written as name=value. Multiple attributes are separated by '|'.
Example:

	type Struct {
		X int `model:"min=0|max=10"`
		S []int `model:"maxlen=3|min=0|max=10"
		hidden int
		Ignore int `gostart:"-"`
		Z int `view:"lable=A longer label for display"`
	}
*/
func (self *MetaData) Attrib(tagKey, name string) (value string, ok bool) {
	if self.attribs == nil {
		// if self.Kind == DynamicValueKind {
		// 	self.attribs = self.Value.Interface().(DynamicValue).Attribs
		// 	if self.attribs == nil {
		// 		return "", false
		// 	}
		// } else {
		self.attribs = make(map[string]map[string]string)
		// }
	}
	keyAttribs, ok := self.attribs[tagKey]
	if !ok {
		// if self.Kind == DynamicValueKind {
		// 	return "", false
		// }
		// Find struct field with tag
		structField := self
		for !structField.IsNamed() {
			structField = structField.Parent
			if structField == nil {
				return "", false
			}
		}
		keyAttribs = ParseTagAttribs(structField.tag.Get(tagKey))
		self.attribs[tagKey] = keyAttribs
	}
	value, ok = keyAttribs[name]
	return value, ok
}

// BoolAttrib uses Attrib() to check if the value of an attribute is "true".
// A non existing attribibute is considered to be false.
func (self *MetaData) BoolAttrib(tagKey, name string) bool {
	value, ok := self.Attrib(tagKey, name)
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

func (self *MetaData) SelectorsMatch(list []string) bool {
	return utils.StringIn(self.Selector(), list) || utils.StringIn(self.WildcardSelector(), list)
}

func (self *MetaData) String() string {
	return fmt.Sprintf("Selector: %s, Kind: %s, Type: %T", self.Selector(), self.Kind, self.Value.Interface())
}
