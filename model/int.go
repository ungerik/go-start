package model

import (
	"fmt"
	"strconv"
)

func NewInt(value int64) *Int {
	return (*Int)(&value)
}

// Attributes:
// * min
// * max
// * label
type Int int64

func (self *Int) Get() int64 {
	return int64(*self)
}

func (self *Int) Set(value int64) {
	*self = Int(value)
}

func (self *Int) GetInt() int {
	return int(*self)
}

func (self *Int) SetInt(value int) {
	*self = Int(value)
}

func (self *Int) String() string {
	return strconv.FormatInt(self.Get(), 10)
}

func (self *Int) SetString(str string) error {
	value, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		self.Set(value)
	}
	return err
}

func (self *Int) IsEmpty() bool {
	return false
}

func (self *Int) Required(metaData *MetaData) bool {
	return true
}

func (self *Int) Validate(metaData *MetaData) error {
	value := int64(*self)
	min, ok, err := self.Min(metaData)
	if err != nil {
		return err
	} else if ok && value < min {
		return &IntBelowMin{value, min}
	}
	max, ok, err := self.Max(metaData)
	if err != nil {
		return err
	} else if ok && value > max {
		return &IntAboveMax{value, max}
	}
	return nil
}

func (self *Int) Min(metaData *MetaData) (min int64, ok bool, err error) {
	str, ok := metaData.Attrib(StructTagKey, "min")
	if !ok {
		return 0, false, nil
	}
	value, err := strconv.ParseInt(str, 10, 64)
	return value, err == nil, err
}

func (self *Int) Max(metaData *MetaData) (max int64, ok bool, err error) {
	str, ok := metaData.Attrib(StructTagKey, "max")
	if !ok {
		return 0, false, nil
	}
	value, err := strconv.ParseInt(str, 10, 64)
	return value, err == nil, err
}

type IntBelowMin struct {
	Value int64
	Min   int64
}

func (self *IntBelowMin) Error() string {
	return fmt.Sprintf("Int %d below minimum of %d", self.Value, self.Min)
}

type IntAboveMax struct {
	Value int64
	Max   int64
}

func (self *IntAboveMax) Error() string {
	return fmt.Sprintf("Int %d above maximum of %d", self.Value, self.Max)
}
