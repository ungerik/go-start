package model

import (
	"fmt"
	"github.com/ungerik/go-start/i18n"
)

func NewLanguage(value string) *Language {
	return (*Language)(&value)
}

/*
Attributes:
	label
	required
*/
type Language string

func (self *Language) Get() string {
	return self.String()
}

func (self *Language) Set(value string) error {
	return self.SetString(value)
}

func (self *Language) IsEmpty() bool {
	return *self == ""
}

func (self *Language) String() string {
	return string(*self)
}

func (self *Language) SetString(str string) error {
	if str != "" {
		//		if _, ok := i18n.Languages()[str]; !ok {
		//			return &InvalidLanguageCode{str}
		//		}
	}
	*self = Language(str)
	return nil
}

func (self *Language) EnglishName() string {
	return i18n.EnglishLanguageName(self.Get())
}

func (self *Language) FixValue(metaData *MetaData) {
}

func (self *Language) Validate(metaData *MetaData) error {
	str := self.Get()
	if self.Required(metaData) || str != "" {
		//		if _, ok := i18n.Languages()[str]; !ok {
		//			errors = append(errors, &InvalidLanguageCode{str})
		//		}
	}
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	return nil
}

func (self *Language) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

type InvalidLanguageCode struct {
	Language string
}

func (self *InvalidLanguageCode) String() string {
	return fmt.Sprintf("Ivalid language code ''", self.Language)
}
