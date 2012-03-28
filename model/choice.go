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

func (self *Choice) FixValue(metaData MetaData) {
}

func (self *Choice) Validate(metaData MetaData) []*ValidationError {
	str := string(*self)
	options := self.Options(metaData)
	if options == nil {
		return NewValidationErrors(errs.Format("model.Choice needs options"), metaData)
	}
	if len(options) == 0 {
		return NewValidationErrors(errs.Format("model.Choice needs options"), metaData)
	}
	found := false
	for _, option := range options {
		if str == option {
			found = true
			break
		}
	}
	if !found {
		return NewValidationErrors(&InvalidChoice{str, options}, metaData)
	}
	return NoValidationErrors
}

func (self *Choice) Options(metaData MetaData) []string {
	options, ok := metaData.TopField().Attrib("options")
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
