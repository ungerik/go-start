package model

import "strings"

type MultipleChoice []string

func (self *MultipleChoice) String() string {
	return strings.Join([]string(*self), ";")
}

func (self *MultipleChoice) SetString(str string) error {
	*self = MultipleChoice(strings.Split(str, ";"))
	return nil
}

func (self *MultipleChoice) IsEmpty() bool {
	return false
}

func (self *MultipleChoice) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib("required")
}

func (self *MultipleChoice) Validate(metaData *MetaData) error {
	// todo
	return nil
}
