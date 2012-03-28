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

func (self *MultipleChoice) FixValue(metaData MetaData) {
}

func (self *MultipleChoice) Validate(metaData MetaData) []*ValidationError {
	// todo
	return NoValidationErrors
}
