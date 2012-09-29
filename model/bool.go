package model

import "strconv"

func NewBool(value bool) *Bool {
	return (*Bool)(&value)
}

/*
Bool model value.
Struct tag attributes:
	required
*/
type Bool bool

func (self *Bool) Get() bool {
	return bool(*self)
}

func (self *Bool) Set(value bool) {
	*self = Bool(value)
}

func (self *Bool) String() string {
	return strconv.FormatBool(self.Get())
}

func (self *Bool) SetString(str string) error {
	value, err := strconv.ParseBool(str)
	if err == nil {
		self.Set(value)
	}
	return err
}

func (self *Bool) IsEmpty() bool {
	return false
}

func (self *Bool) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Bool) Validate(metaData *MetaData) error {
	if self.Required(metaData) && *self == false {
		return NewRequiredError(metaData)
	}
	return nil
}
