package model

import "strconv"

func NewRichText(value string) *RichText {
	return (*RichText)(&value)
}

type RichText string

func (self *RichText) Get() string {
	return string(*self)
}

func (self *RichText) Set(value string) {
	*self = RichText(value)
}

func (self *RichText) IsEmpty() bool {
	return len(*self) == 0
}

func (self *RichText) GetOrDefault(defaultText string) string {
	if self.IsEmpty() {
		return defaultText
	}
	return self.Get()
}

func (self *RichText) String() string {
	return self.Get()
}

func (self *RichText) SetString(str string) error {
	self.Set(str)
	return nil
}

func (self *RichText) FixValue(metaData *MetaData) {
}

func (self *RichText) Required(metaData *MetaData) bool {
	if minlen, ok, _ := self.Minlen(metaData); ok {
		if minlen > 0 {
			return true
		}
	}
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *RichText) Validate(metaData *MetaData) error {
	value := string(*self)

	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}

	minlen, ok, err := self.Minlen(metaData)
	if ok && len(value) < minlen {
		err = &StringTooShort{value, minlen}
	}
	if err != nil {
		return err
	}

	maxlen, ok, err := self.Maxlen(metaData)
	if ok && len(value) > maxlen {
		err = &StringTooLong{value, maxlen}
	}
	if err != nil {
		return err
	}

	return nil
}

func (self *RichText) Minlen(metaData *MetaData) (minlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "minlen"); ok {
		minlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return minlen, ok, err
}

func (self *RichText) Maxlen(metaData *MetaData) (maxlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "maxlen"); ok {
		maxlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return maxlen, ok, err
}
