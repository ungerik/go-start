package model

import (
	"fmt"
	"math"
	"strconv"
)

func NewFloat(value float64) *Float {
	return (*Float)(&value)
}

// Attributes:
// * label
// * min
// * max
// * valid
type Float float64

func (self *Float) Get() float64 {
	return float64(*self)
}

func (self *Float) Set(value float64) {
	*self = Float(value)
}

func (self *Float) IsValid() bool {
	value := self.Get()
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func (self *Float) String() string {
	return strconv.FormatFloat(self.Get(), 'f', -1, 64)
}

func (self *Float) SetString(str string) error {
	value, err := strconv.ParseFloat(str, 64)
	if err == nil {
		self.Set(value)
	}
	return err
}

func (self *Float) IsEmpty() bool {
	return false
}

func (self *Float) Required(metaData *MetaData) bool {
	return true
}

func (self *Float) Validate(metaData *MetaData) error {
	value := float64(*self)
	min, ok, err := self.Min(metaData)
	if err != nil {
		return err
	} else if ok && value < min {
		return &FloatBelowMin{value, min}
	}
	max, ok, err := self.Max(metaData)
	if err != nil {
		return err
	} else if ok && value > max {
		return &FloatAboveMax{value, max}
	}
	if valid := self.Valid(metaData); valid && !self.IsValid() {
		return &FloatNotReal{value}
	}
	return nil
}

func (self *Float) Min(metaData *MetaData) (min float64, ok bool, err error) {
	str, ok := metaData.Attrib(StructTagKey, "min")
	if !ok {
		return 0, false, nil
	}
	value, err := strconv.ParseFloat(str, 64)
	return value, err == nil, err
}

func (self *Float) Max(metaData *MetaData) (max float64, ok bool, err error) {
	str, ok := metaData.Attrib(StructTagKey, "max")
	if !ok {
		return 0, false, nil
	}
	value, err := strconv.ParseFloat(str, 64)
	return value, err == nil, err
}

func (self *Float) Valid(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "valid")
}

type FloatBelowMin struct {
	Value float64
	Min   float64
}

func (self *FloatBelowMin) Error() string {
	return fmt.Sprintf("Float %f below minimum of %f", self.Value, self.Min)
}

type FloatAboveMax struct {
	Value float64
	Max   float64
}

func (self *FloatAboveMax) Error() string {
	return fmt.Sprintf("Float %f above maximum of %f", self.Value, self.Max)
}

type FloatNotReal struct {
	Value float64
}

func (self *FloatNotReal) Error() string {
	return fmt.Sprintf("Float %v is not a real number", self.Value)
}
