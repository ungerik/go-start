package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
)

type PasswordFormModel struct {
	Password1 model.Password `model:"minlen=6" view:"label=Password|size=20"`
	Password2 model.Password `view:"label=Repeat password|size=20"`
}

func (self *PasswordFormModel) Validate(metaData *model.MetaData) error {
	if self.Password1 != self.Password2 {
		return errors.New("Passwords don't match")
	}
	return nil
}

type EmailPasswordFormModel struct {
	Email             model.Email `model:"required" view:"size=20"`
	PasswordFormModel `bson:",inline" view:"size=20"`
}

type LoginFormModel struct {
	Email    model.Email    `model:"required" view:"size=20"`
	Password model.Password `model:"required" view:"size=20"`
}
