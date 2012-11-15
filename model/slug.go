package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func NewSlug(value string) *Slug {
	s := new(Slug)
	s.Set(value)
	return s
}

type Slug string

func (self *Slug) Get() string {
	return string(*self)
}

func (self *Slug) Set(value string) {

	result := make([]byte, utf8.RuneCountInString(value))
	i := 0
	for _, c := range value {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' {
			result[i] = byte(c)
		} else if c >= 'A' && c <= 'Z' {
			result[i] = byte(unicode.ToLower(c))
		} else {
			result[i] = '-'
		}
		i++
	}

	*self = Slug(string(result))
}

func (self *Slug) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Slug) GetOrDefault(defaultSlug string) string {
	if self.IsEmpty() {
		return defaultSlug
	}
	return self.Get()
}

func (self *Slug) String() string {
	return self.Get()
}

func (self *Slug) SetSlug(str string) error {
	self.Set(str)
	return nil
}

func (self *Slug) FixValue(metaData *MetaData) {
}

func (self *Slug) Required(metaData *MetaData) bool {
	if minlen, ok, _ := self.Minlen(metaData); ok {
		if minlen > 0 {
			return true
		}
	}
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Slug) Validate(metaData *MetaData) error {
	value := string(*self)

	pos := strings.IndexAny(value, "\n\r")
	if pos != -1 {
		return errors.New("Line breaks not allowed")
	}

	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}

	minlen, ok, err := self.Minlen(metaData)
	if ok && len(value) < minlen {
		err = &SlugTooShort{value, minlen}
	}
	if err != nil {
		return err
	}

	maxlen, ok, err := self.Maxlen(metaData)
	if ok && len(value) > maxlen {
		err = &SlugTooLong{value, maxlen}
	}
	if err != nil {
		return err
	}

	return nil
}

func (self *Slug) Minlen(metaData *MetaData) (minlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "minlen"); ok {
		minlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return minlen, ok, err
}

func (self *Slug) Maxlen(metaData *MetaData) (maxlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "maxlen"); ok {
		maxlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return maxlen, ok, err
}

type SlugTooShort struct {
	Str    string
	Minlen int
}

func (self *SlugTooShort) Error() string {
	return fmt.Sprintf("Slug shorter than minimum of %d characters", self.Minlen)
}

type SlugTooLong struct {
	Str    string
	Maxlen int
}

func (self *SlugTooLong) Error() string {
	return fmt.Sprintf("String longer than maximum of %d characters", self.Maxlen)
}
