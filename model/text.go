package model

import "strconv"

func NewText(value string) *Text {
	return (*Text)(&value)
}

type Text string

func (self *Text) Get() string {
	return string(*self)
}

func (self *Text) Set(value string) {
	*self = Text(value)
}

func (self *Text) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Text) GetOrDefault(defaultText string) string {
	if self.IsEmpty() {
		return defaultText
	}
	return self.Get()
}

func (self *Text) String() string {
	return self.Get()
}

func (self *Text) SetString(str string) error {
	self.Set(str)
	return nil
}

func (self *Text) FixValue(metaData *MetaData) {
}

func (self *Text) Required(metaData *MetaData) bool {
	if minlen, ok, _ := self.Minlen(metaData); ok {
		if minlen > 0 {
			return true
		}
	}
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Text) Validate(metaData *MetaData) error {
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

func (self *Text) Minlen(metaData *MetaData) (minlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "minlen"); ok {
		minlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return minlen, ok, err
}

func (self *Text) Maxlen(metaData *MetaData) (maxlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "maxlen"); ok {
		maxlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return maxlen, ok, err
}
