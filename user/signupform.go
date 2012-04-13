package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

// confirmationPage must have the confirmation code as first URL parameter
func NewSignupForm(class, errorMessageClass, successMessageClass string, confirmationURL, redirectURL view.URL) *view.Form {
	return &view.Form{
		Class:               class,
		ErrorMessageClass:   errorMessageClass,
		SuccessMessageClass: successMessageClass,
		SuccessMessage:      Config.ConfirmationSent,
		ButtonText:          "Signup",
		FormID:              "gostart_user_signup",
		GetModel: func(form *view.Form, context *view.Context) (interface{}, error) {
			return &SignupFormModel{}, nil
		},
		Redirect: redirectURL,
		OnSubmit: func(form *view.Form, formModel interface{}, context *view.Context) error {
			m := formModel.(*SignupFormModel)
			email := m.Email.Get()
			password := m.Password1.Get()
			var user *User
			doc, found, err := FindByEmail(email)
			if err != nil {
				return err
			}
			if found {
				user = From(doc)
				if user.EmailPasswordConfirmed() {
					return errors.New("A user with that email and a password already exists")
				}
				user.Password.SetHashed(password)
			} else {
				user, _, err = New(email, password)
				if err != nil {
					return err
				}
			}
			err = <-user.Email[0].SendConfirmationEmail(context, confirmationURL)
			if err != nil {
				return err
			}
			return user.Save()
		},
	}
}

///////////////////////////////////////////////////////////////////////////////
// SignupFormModel

type SignupFormModel struct {
	Email     model.Email    `gostart:"required"`
	Password1 model.Password `gostart:"required|label=Password|minlen=6"`
	Password2 model.Password `gostart:"label=Repeat password"`
}

func (self *SignupFormModel) Validate(metaData model.MetaData) []*model.ValidationError {
	if self.Password1 != self.Password2 {
		return model.NewValidationErrors(errors.New("Passwords don't match"), metaData)
	}
	return model.NoValidationErrors
}
