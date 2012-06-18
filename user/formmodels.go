package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
)

type PasswordFormModel struct {
	Password1 model.Password `gostart:"label=Password|minlen=6|size=20"`
	Password2 model.Password `gostart:"label=Repeat password|size=20"`
}

func (self *PasswordFormModel) Validate(metaData *model.MetaData) error {
	if self.Password1 != self.Password2 {
		return errors.New("Passwords don't match")
	}
	return nil
}

type EmailPasswordFormModel struct {
	Email             model.Email `gostart:"required|size=20"`
	PasswordFormModel `bson:",inline" gostart:"size=20"`
}

type LoginFormModel struct {
	Email    model.Email    `gostart:"required|size=20"`
	Password model.Password `gostart:"required|size=20"`
}
