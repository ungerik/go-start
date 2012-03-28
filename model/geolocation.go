package model

import "fmt"

// todo change data type

type GeoLocation struct {
	Longitude, Latitude float64
}

func (self *GeoLocation) String() string {
	return fmt.Sprintf("%v / %v", self.Longitude, self.Latitude)
}

func (self *GeoLocation) SetString(str string) error {
	panic("not implemented")
	return nil
}

func (self *GeoLocation) FixValue(metaData MetaData) {
}

func (self *GeoLocation) Validate(metaData MetaData) []*ValidationError {
	return NoValidationErrors
}
