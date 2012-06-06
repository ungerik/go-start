package model

import "fmt"

// todo change data type

type GeoLocation struct {
	Longitude Float
	Latitude  Float
}

func (self *GeoLocation) String() string {
	return fmt.Sprintf("%v / %v", self.Longitude, self.Latitude)
}

func (self *GeoLocation) SetString(str string) error {
	panic("not implemented")
	return nil
}

func (self *GeoLocation) IsEmpty() bool {
	return false
}

func (self *GeoLocation) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib("required")
}

func (self *GeoLocation) Validate(metaData *MetaData) error {
	return nil
}
