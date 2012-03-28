package model

import "strconv"

// Attributes:
// * label
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

func (self *Bool) FixValue(metaData MetaData) {
}

func (self *Bool) Validate(metaData MetaData) []*ValidationError {
	return NoValidationErrors
}
