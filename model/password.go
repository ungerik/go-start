package model

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

type Password string

func (self *Password) Get() string {
	return string(*self)
}

func (self *Password) Set(value string) {
	*self = Password(value)
}

func (self *Password) SetHashed(value string) {
	self.Set(PasswordHash(value))
}

func (self *Password) EqualsHashed(value string) bool {
	return self.Get() == PasswordHash(value)
}

func (self *Password) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Password) String() string {
	return self.Get()
}

func (self *Password) SetString(str string) error {
	self.Set(str)
	return nil
}

func (self *Password) FixValue(metaData *MetaData) {
}

func (self *Password) Validate(metaData *MetaData) error {
	value := string(*self)

	pos := strings.IndexAny(value, "\n\r")
	if pos != -1 {
		return errors.New("Line breaks not allowed")
	}

	minlen, ok, err := self.Minlen(metaData)
	if ok && len(value) < minlen {
		err = &StringTooShort{value, minlen}
	}
	if err != nil {
		return err
	}

	return nil
}

func (self *Password) Minlen(metaData *MetaData) (minlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib("minlen"); ok {
		minlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return minlen, ok, err
}

func PasswordHash(password string) (hash string) {
	sha := sha256.New()
	sha.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}
