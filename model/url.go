package model

import (
	"net/url"
	"strings"
)

func NewUrl(value string) *Url {
	return (*Url)(&value)
}

type Url string

func (self *Url) Get() string {
	return string(*self)
}

func (self *Url) Set(value string) error {
	if value == "" {
		*self = ""
		return nil
	}
	if !strings.HasPrefix(value, "http") {
		value = "http://" + value
	}
	_, err := url.Parse(value)
	if err == nil {
		*self = Url(value)
	}
	return err
}

func (self *Url) SetDataUrl(value string) error {
	if value == "" {
		*self = ""
		return nil
	}
	_, err := url.Parse(value)
	if err == nil {
		*self = Url(value)
	}
	return err
}

func (self *Url) IsEmpty() bool {
	return *self == ""
}

func (self *Url) GetOrDefault(defaultURL string) string {
	if self.IsEmpty() {
		return defaultURL
	}
	return self.Get()
}

func (self *Url) String() string {
	return self.Get()
}

func (self *Url) SetString(str string) error {
	return self.Set(str)
}

func (self *Url) FixValue(metaData *MetaData) {
}

func (self *Url) Validate(metaData *MetaData) error {
	if self.IsEmpty() {
		if self.Required(metaData) {
			return NewRequiredError(metaData)
		}
	} else {
		if _, err := url.Parse(self.Get()); err != nil {
			return err
		}
	}
	return nil
}

func (self *Url) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}
