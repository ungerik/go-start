package model

import (
	"strings"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/utils"
)

type MultipleChoice []string

func (self *MultipleChoice) Get() []string {
	return []string(*self)
}

func (self *MultipleChoice) Set(value ...string) {
	*self = MultipleChoice(value)
}

func (self *MultipleChoice) String() string {
	return strings.Join([]string(*self), "|")
}

func (self *MultipleChoice) SetString(str string) error {
	*self = MultipleChoice(strings.Split(str, "|"))
	return nil
}

func (self *MultipleChoice) IsSet(option string) bool {
	return utils.StringIn(option, self.Get())
}

func (self *MultipleChoice) IsEmpty() bool {
	return len(*self) == 0
}

func (self *MultipleChoice) Options(metaData *MetaData) []string {
	return choiceOptions(metaData)
}

func (self *MultipleChoice) Required(metaData *MetaData) bool {
	return false
}

func (self *MultipleChoice) Validate(metaData *MetaData) error {
	options := self.Options(metaData)
	if len(options) == 0 {
		return errs.Format("model.MultipleChoice needs options")
	}
	for _, option := range options {
		if option == "" {
			return errs.Format("model.MultipleChoice options must not be empty strings")
		}
	}
	for _, str := range *self {
		if !utils.StringIn(str, options) {
			return &InvalidChoice{str, options}
		}
	}
	return nil
}
