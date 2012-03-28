package model

import "net/url"

type Url string

func (self *Url) Get() string {
	return string(*self)
}

func (self *Url) Set(value string) error {
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

func (self *Url) FixValue(metaData MetaData) {
}

func (self *Url) Validate(metaData MetaData) []*ValidationError {
	errors := NoValidationErrors
	if self.IsEmpty() {
		if self.Required(metaData) {
			errors = append(errors, NewRequiredValidationError(metaData))
		}
	} else {
		if _, err := url.Parse(self.Get()); err != nil {
			errors = append(errors, &ValidationError{err, metaData})
		}
	}
	return errors
}

func (self *Url) Required(metaData MetaData) bool {
	return metaData.TopField().BoolAttrib("required")
}
