package model

import (
	"fmt"
	"strings"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/utils"
)

func NewChoice(value string) *Choice {
	return (*Choice)(&value)
}

/*
Choice can hold one of several string options.
The options are defined by the struct tag attribute "options"
delimited by colons ",". 
A colon in an option can be escaped with "\,"
An empty string can be used as the first option.
Choice is required, when the first option in not an empty string.
Struct tag attributes:
	`model:"options=Red,Green,Blue"`
	`model:"options=,Second,Third"` // empty string is a valid value
*/
type Choice string

func (self *Choice) Get() string {
	return string(*self)
}

func (self *Choice) Set(value string) {
	*self = Choice(value)
}

func (self *Choice) String() string {
	return self.Get()
}

func (self *Choice) SetString(str string) error {
	self.Set(str)
	return nil
}

func choiceOptions(metaData *MetaData) []string {
	options, ok := metaData.Attrib(StructTagKey, "options")
	if !ok {
		return nil
	}
	options = strings.Replace(options, "\\,", "|", -1)
	options = strings.Replace(options, ",", "|", -1)
	return strings.Split(options, "|")
}

func (self *Choice) Options(metaData *MetaData) []string {
	return choiceOptions(metaData)
}

func (self *Choice) IsEmpty() bool {
	return *self == ""
}

func (self *Choice) Required(metaData *MetaData) bool {
	options := choiceOptions(metaData)
	return len(options) > 0 && options[0] != ""
}

func (self *Choice) Validate(metaData *MetaData) error {
	str := string(*self)
	options := self.Options(metaData)
	if len(options) == 0 {
		return errs.Format("model.Choice needs options")
	}
	if !utils.StringIn(str, options) {
		return &InvalidChoice{str, options}
	}
	return nil
}

type InvalidChoice struct {
	Value   string
	Options []string
}

func (self *InvalidChoice) Error() string {
	return fmt.Sprintf("Invalid choice %s (options: %s)", self.Value, strings.Join(self.Options, ","))
}
