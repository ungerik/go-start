package model

import (
	"fmt"
	// "github.com/ungerik/go-start/debug"
)

// todo change data type

type GeoLocation struct {
	Longitude Float
	Latitude  Float
}

func (self *GeoLocation) String() string {
	return fmt.Sprintf("%v / %v", self.Longitude, self.Latitude)
}

func (self *GeoLocation) SetString(str string) error {
	// At the moment GeoLocation will be handled as a struct of Float,
	// so there is no value for GeoLocation itself when calling SetString
	// from Form, only two separate Floats for the struct.
	// That's why we ignore an empty string for the moment
	if str == "" {
		return nil
	}
	panic("not implemented")
}

func (self *GeoLocation) IsEmpty() bool {
	return false
}

func (self *GeoLocation) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *GeoLocation) Validate(metaData *MetaData) error {
	return nil
}
