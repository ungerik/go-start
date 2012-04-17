package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type String string

func (self *String) Get() string {
	return string(*self)
}

func (self *String) Set(value string) {
	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, "\r", "", -1)
	*self = String(value)
}

func (self *String) IsEmpty() bool {
	return len(*self) == 0
}

func (self *String) GetOrDefault(defaultString string) string {
	if self.IsEmpty() {
		return defaultString
	}
	return self.Get()
}

func (self *String) String() string {
	return self.Get()
}

func (self *String) SetString(str string) error {
	self.Set(str)
	return nil
}

func (self *String) FixValue(metaData *MetaData) {
}

func (self *String) Validate(metaData *MetaData) []*ValidationError {
	value := string(*self)
	e := NoValidationErrors

	pos := strings.IndexAny(value, "\n\r")
	if pos != -1 {
		e = append(e, &ValidationError{errors.New("model.String contains line breaks"), metaData})
	}

	minlen, ok, err := self.Minlen(metaData)
	if ok && len(value) < minlen {
		err = &StringTooShort{value, minlen}
	}
	if err != nil {
		e = append(e, &ValidationError{err, metaData})
	}

	maxlen, ok, err := self.Maxlen(metaData)
	if ok && len(value) > maxlen {
		err = &StringTooLong{value, maxlen}
	}
	if err != nil {
		e = append(e, &ValidationError{err, metaData})
	}

	return e
}

func (self *String) Minlen(metaData *MetaData) (minlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib("minlen"); ok {
		minlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return minlen, ok, err
}

func (self *String) Maxlen(metaData *MetaData) (maxlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib("maxlen"); ok {
		maxlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return maxlen, ok, err
}

func (self *String) Hidden(metaData *MetaData) (hidden bool) {
	return metaData.BoolAttrib("hidden")
}

type StringTooShort struct {
	Str    string
	Minlen int
}

func (self *StringTooShort) Error() string {
	return fmt.Sprintf("String shorter than minimum of %d characters ('%s')", self.Minlen, self.Str)
}

type StringTooLong struct {
	Str    string
	Maxlen int
}

func (self *StringTooLong) Error() string {
	return fmt.Sprintf("String longer than maximum of %d characters ('%s')", self.Maxlen, self.Str)
}
