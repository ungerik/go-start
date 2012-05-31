package model

import (
	"crypto/sha256"
	"encoding/base64"
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
	return nil
}

func PasswordHash(password string) (hash string) {
	sha := sha256.New()
	sha.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}
