package model

import (
// "io/ioutil"
)

func NewBlob(value []byte) *Blob {
	return (*Blob)(&value)
}

/*
Blob is just a bunch of bytes.
Struct tag attributes:
	`model:"required"`
*/
type Blob []byte

func (self *Blob) Get() []byte {
	return []byte(*self)
}

func (self *Blob) Set(value []byte) {
	*self = Blob(value)
}

func (self *Blob) String() string {
	return string(*self)
}

func (self *Blob) SetString(str string) error {
	*self = Blob(str)
	return nil
}

func (self *Blob) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Blob) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Blob) Validate(metaData *MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	return nil
}
