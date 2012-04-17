package model

import (
	"github.com/ungerik/go-start/utils"
	"strings"
)

type Phone string

func (self *Phone) Get() string {
	return string(*self)
}

func (self *Phone) Set(value string) {
	*self = Phone(NormalizePhoneNumber(value))
}

func (self *Phone) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Phone) String() string {
	return self.Get()
}

func (self *Phone) SetString(str string) error {
	self.Set(str)
	return nil
}

func (self *Phone) FixValue(metaData *MetaData) {
}

func (self *Phone) Validate(metaData *MetaData) []*ValidationError {
	return NoValidationErrors
}

func NormalizePhoneNumber(number string) string {
	if number == "" {
		return ""
	}
	number = strings.Replace(number, " ", "", -1)
	number = strings.Replace(number, "/", "", -1)
	number = strings.Replace(number, "-", "", -1)
	number = strings.Replace(number, ".", "", -1)
	number = strings.Replace(number, "(0)", "", -1)
	number = strings.Replace(number, ")", "", -1)
	number = strings.Replace(number, "(", "", -1)
	if utils.StringStartsWith(number, "++") {
		return "00" + number[2:]
	}
	if utils.StringStartsWith(number, "+") {
		return "00" + number[1:]
	}
	return number
}
