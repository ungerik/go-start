package model

import (
	"fmt"
	"github.com/ungerik/go-start/i18n"
)

func NewCountry(value string) *Country {
	return (*Country)(&value)
}

type Country string

func (self *Country) Get() string {
	return string(*self)
}

func (self *Country) Set(value string) error {
	if value != "" {
		//		if _, ok := i18n.Countries()[value]; !ok {
		//			return &InvalidCountryCode{str}
		//		}
	}
	*self = Country(value)
	return nil
}

func (self *Country) IsEmpty() bool {
	return *self == ""
}

func (self *Country) String() string {
	return self.Get()
}

func (self *Country) SetString(str string) error {
	return self.Set(str)
}

func (self *Country) EnglishName() string {
	return i18n.EnglishCountryName(self.Get())
}

func (self *Country) FixValue(metaData *MetaData) {
}

func (self *Country) Validate(metaData *MetaData) error {
	str := self.Get()
	if self.Required(metaData) || str != "" {
		//		if _, ok := i18n.Countries()[value]; !ok {
		//			return NewValidationErrors(&InvalidCountryCode{str}, metaData)
		//		}
	}
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	return nil
}

func (self *Country) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

type InvalidCountryCode struct {
	CountryCode string
}

func (self *InvalidCountryCode) String() string {
	return fmt.Sprintf("Ivalid country code ''", self.CountryCode)
}
