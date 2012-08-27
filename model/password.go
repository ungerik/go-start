package model

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"crypto/rand"
	"io"
)

type Password string

func (self *Password) Get() string {
	return string(*self)
}

func (self *Password) Set(value string) {
	*self = Password(value)
}

// todo salt
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

func (self *Password) Required(metaData *MetaData) bool {
	if minlen, ok, _ := self.Minlen(metaData); ok {
		if minlen > 0 {
			return true
		}
	}
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Password) Validate(metaData *MetaData) error {
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
		err = &StringTooShort{value, minlen}
	}
	if err != nil {
		return err
	}

	maxlen, ok, err := self.Maxlen(metaData)
	if ok && len(value) > maxlen {
		err = &StringTooLong{value, maxlen}
	}
	if err != nil {
		return err
	}

	return nil
}

func (self *Password) Minlen(metaData *MetaData) (minlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "minlen"); ok {
		minlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return minlen, ok, err
}

func (self *Password) Maxlen(metaData *MetaData) (maxlen int, ok bool, err error) {
	var str string
	if str, ok = metaData.Attrib(StructTagKey, "maxlen"); ok {
		maxlen, err = strconv.Atoi(str)
		ok = err == nil
	}
	return maxlen, ok, err
}

func PasswordHash(password string) (string) {
	sha := sha256.New()
	sha.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}

func GenerateRandomPassword(length int) (string, error) {
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	
	b := make([]byte, length)
	r := make([]byte, length+(length/4))
	clen := byte(len(validChars))
	maxrb := byte(256 - (256 % len(validChars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, r); err != nil {
			return "", errors.New("error reading from random source: " + err.Error())
		}
		for _, c := range r {
			if c >= maxrb {
				continue
			}
			b[i] = validChars[c%clen]
			i++
			if i == length {
				return string(b), nil
			}
		}
	}
	return "",errors.New("password generation unreachable")

}

