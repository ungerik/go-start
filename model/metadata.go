package model

import (
	"reflect"
	"strings"
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// MetaData

type MetaData []FieldMetaData

func (self MetaData) TopField() *FieldMetaData {
	return &self[len(self)-1]
}

func (self MetaData) Selector() string {
	var sb utils.StringBuilder
	for i := range self {
		if i != 0 {
			sb.Byte('.')
		}
		sb.Write(self[i].Name)
	}
	return sb.String()
}

func (self MetaData) ArrayWildcardSelector() string {
	var sb utils.StringBuilder
	for i := range self {
		if i != 0 {
			sb.Byte('.')
		}
		m := &self[i]
		if m.IsIndex() {
			sb.Byte('$')
		} else {
			sb.Write(m.Name)
		}
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////////////////
// FieldMetaData

type FieldMetaData struct {
	ParentStruct reflect.Value
	Name         string
	Index        int
	tag          string
	attributes   map[string]string
}

func (self *FieldMetaData) IsIndex() bool {
	return self.Index != -1
}

func (self *FieldMetaData) Attrib(name string) (value string, ok bool) {
	if self.attributes == nil {
		self.attributes = map[string]string{}
		for _, s := range strings.Split(self.tag, "|") {
			pos := strings.Index(s, "=")
			if pos == -1 {
				self.attributes[s] = "true"
			} else {
				self.attributes[s[:pos]] = s[pos+1:]
			}
		}
	}
	value, ok = self.attributes[name]
	return value, ok
}

func (self *FieldMetaData) BoolAttrib(name string) bool {
	value, ok := self.Attrib(name)
	return ok && value == "true"
}


