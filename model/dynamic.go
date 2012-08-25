package model

import (
	"reflect"
)

type Dynamic []DynamicValue

func (self Dynamic) Get(name string) (value *DynamicValue, found bool) {
	for i := range self {
		if self[i].Name == name {
			return &self[i], true
		}
	}
	return nil, false
}

func (self Dynamic) GetValue(name string) (value Value, found bool) {
	if v, found := self.Get(name); found {
		return v.Value, true
	}
	return nil, false
}

var typeOfDynamicValue = reflect.TypeOf(DynamicValue{})

type DynamicValue struct {
	Name string
	Value
	Attribs map[string]map[string]string
}

func (self *DynamicValue) Attrib(tagKey, name string) string {
	if attribs, ok := self.Attribs[tagKey]; ok {
		if attrib, ok := attribs[name]; ok {
			return attrib
		}
	}
	return ""
}

func (self *DynamicValue) SetAttrib(tagKey, name, value string) {
	if self.Attribs == nil {
		self.Attribs = make(map[string]map[string]string)
	}
	self.Attribs[tagKey] = map[string]string{name: value}
}

// func (self *DynamicValue) String() *String {
// 	return self.Value.(*String)
// }
