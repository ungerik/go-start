package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
)

type PasswordFormModel struct {
	Password1 model.Password `gostart:"required|label=Password|minlen=6"`
	Password2 model.Password `gostart:"label=Repeat password"`
}

func (self *PasswordFormModel) Validate(metaData *model.MetaData) []*model.ValidationError {
	if self.Password1 != self.Password2 {
		return model.NewValidationErrors(errors.New("Passwords don't match"), metaData)
	}
	return model.NoValidationErrors
}

type EmailPasswordFormModel struct {
	Email             model.Email `gostart:"required"`
	PasswordFormModel `bson:",inline"`
}


