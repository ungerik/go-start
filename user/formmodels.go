package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
)

type PasswordFormModel struct {
	Password1 model.Password `gostart:"label=Password|minlen=6"`
	Password2 model.Password `gostart:"label=Repeat password"`
}

func (self *PasswordFormModel) Validate(metaData *model.MetaData) error {
	if self.Password1 != self.Password2 {
		return errors.New("Passwords don't match")
	}
	return nil
}

type EmailPasswordFormModel struct {
	Email             model.Email `gostart:"required"`
	PasswordFormModel `bson:",inline"`
}

type LoginFormModel struct {
	Email    model.Email    `gostart:"required"`
	Password model.Password `gostart:"required"`
}
