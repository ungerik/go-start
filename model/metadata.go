package model

import (
	"reflect"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// MetaData

type MetaData struct {
	ParentStruct reflect.Value
	Parent       *MetaData
	Depth        int
	Name         string
	Index        int
	tag          string
	attributes   map[string]string
}

func (self *MetaData) IsIndex() bool {
	return self.Index != -1
}

func (self *MetaData) Attrib(name string) (value string, ok bool) {
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

func (self *MetaData) BoolAttrib(name string) bool {
	value, ok := self.Attrib(name)
	return ok && value == "true"
}

func (self *MetaData) Selector() string {
	names := make([]string, self.Depth)
	for i, m := self.Depth-1, self; i >= 0; i-- {
		names[i] = m.Name
		m = m.Parent
	}
	return strings.Join(names, ".")
}

func (self *MetaData) ArrayWildcardSelector() string {
	names := make([]string, self.Depth)
	for i, m := self.Depth-1, self; i >= 0; i-- {
		if m.IsIndex() {
			names[i] = "$"
		} else {
			names[i] = m.Name
		}
		m = m.Parent
	}
	return strings.Join(names, ".")
}
