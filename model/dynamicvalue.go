package model

import (
	"reflect"

	"labix.org/v2/mgo/bson"
)

var DynamicValueType = reflect.TypeOf(DynamicValue{})

// DynamicValue can be used to build a model dynamically
// instead of using a static typed struct.
// See DynamicValues
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

func (self *DynamicValue) GetBSON() (interface{}, error) {
	if self.Value == nil {
		return nil, nil
	}
	v := reflect.ValueOf(self.Value).Elem()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface(), nil
}

func (self *DynamicValue) SetBSON(raw bson.Raw) error {
	if self.Value == nil {
		panic("Can't SetBSON for nil model.Value")
	}
	return raw.Unmarshal(reflect.ValueOf(self.Value).Elem().Interface())
}

///////////////////////////////////////////////////////////////////////////////

// DynamicValues is used to create a model dynamically instead
// of using a static typed struct.
type DynamicValues []DynamicValue

func (self DynamicValues) Get(name string) (value *DynamicValue, found bool) {
	for i := range self {
		if self[i].Name == name {
			return &self[i], true
		}
	}
	return nil, false
}

func (self DynamicValues) GetValue(name string) (value Value, found bool) {
	if v, found := self.Get(name); found {
		return v.Value, true
	}
	return nil, false
}
