package model

import (
	"fmt"
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// DynamicChoice

type DynamicChoice struct {
	index   int
	options []string
}

func (self *DynamicChoice) Index() int {
	return self.index
}

func (self *DynamicChoice) SetIndex(index int) error {
	err := self.CheckIndex(index)
	if err == nil {
		self.index = index
	}
	return err
}

func (self *DynamicChoice) Options() []string {
	return self.options
}

func (self *DynamicChoice) SetOptions(options []string) {
	self.options = options
}

func (self *DynamicChoice) String() string {
	if self.index < 0 || self.index >= len(self.options) {
		return ""
	}
	return self.options[self.index]
}

func (self *DynamicChoice) SetString(str string) error {
	for i, option := range self.options {
		if str == option {
			self.index = i
			return nil
		}
	}
	// Temporary Hot Fix
	return nil
	// return &InvalidChoice{str, self.options}
}

func (self *DynamicChoice) GetBSON() (interface{}, error) {
	return self.String(), nil
}

func (self *DynamicChoice) SetBSON(raw bson.Raw) (err error) {
	var s string
	err = raw.Unmarshal(&s)
	if err != nil {
		return err
	}
	return self.SetString(s)
}

func (self *DynamicChoice) IsEmpty() bool {
	return self.String() == ""
}

func (self *DynamicChoice) Required(metaData *MetaData) bool {
	return len(self.options) > 0 && self.options[0] != ""
}

func (self *DynamicChoice) Validate(metaData *MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	return self.CheckIndex(self.index)
}

func (self *DynamicChoice) CheckIndex(index int) error {
	if index < 0 || index >= len(self.options) {
		return fmt.Errorf("Choice index %d out of range [0..%d)", index, len(self.options))
	}
	return nil
}
