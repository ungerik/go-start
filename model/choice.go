package model

import (
	"fmt"
	"github.com/ungerik/go-start/errs"
	"strings"
)

/*
Separate options with ',' escape ',' in options with '\,'
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

func (self *Choice) IsEmpty() bool {
	return *self == ""
}

func (self *Choice) Validate(metaData *MetaData) error {
	str := string(*self)
	options := self.Options(metaData)
	if len(options) == 0 {
		return errs.Format("model.Choice needs options")
	}
	found := false
	for _, option := range options {
		if str == option {
			found = true
			break
		}
	}
	if !found {
		return &InvalidChoice{str, options}
	}
	return nil
}

func (self *Choice) Options(metaData *MetaData) []string {
	options, ok := metaData.Attrib("options")
	if !ok {
		return nil
	}
	options = strings.Replace(options, "\\,", "|", -1)
	options = strings.Replace(options, ",", "|", -1)
	return strings.Split(options, "|")
}

type InvalidChoice struct {
	Value   string
	Options []string
}

func (self *InvalidChoice) Error() string {
	return fmt.Sprintf("Invalid choice '%s'. Options: %#v", self.Value, self.Options)
}
