package model

import (
	"fmt"
	"math"
	"strconv"
)

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

func (self *Float) FixValue(metaData MetaData) {
}

func (self *Float) Validate(metaData MetaData) []*ValidationError {
	value := float64(*self)
	errors := NoValidationErrors
	min, ok, err := self.Min(metaData)
	if err != nil {
		errors = append(errors, &ValidationError{err, metaData})
	} else if ok && value < min {
		errors = append(errors, &ValidationError{&FloatBelowMin{value, min}, metaData})
	}
	max, ok, err := self.Max(metaData)
	if err != nil {
		errors = append(errors, &ValidationError{err, metaData})
	} else if ok && value > max {
		errors = append(errors, &ValidationError{&FloatAboveMax{value, max}, metaData})
	}
	if valid := self.Valid(metaData); valid && !self.IsValid() {
		errors = append(errors, &ValidationError{&FloatNotReal{value}, metaData})
	}
	return errors
}

func (self *Float) Min(metaData MetaData) (min float64, ok bool, err error) {
	str, ok := metaData.TopField().Attrib("min")
	if !ok {
		return 0, false, nil
	}
	value, err := strconv.ParseFloat(str, 64)
	return value, err == nil, err
}

func (self *Float) Max(metaData MetaData) (max float64, ok bool, err error) {
	str, ok := metaData.TopField().Attrib("max")
	if !ok {
		return 0, false, nil
	}
	value, err := strconv.ParseFloat(str, 64)
	return value, err == nil, err
}

func (self *Float) Valid(metaData MetaData) bool {
	return metaData.TopField().BoolAttrib("valid")
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

func (self *Float) Hidden(metaData MetaData) (hidden bool) {
	return metaData.TopField().BoolAttrib("hidden")
}
