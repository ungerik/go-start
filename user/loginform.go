package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
	//	"net/url"

	// "github.com/ungerik/go-start/debug"
)

func NewLoginForm(class, errorMessageClass, successMessageClass string, redirectURL view.URL) view.View {
	return view.NewDynamicView(
		func(context *view.Context) (v view.View, err error) {
			var r view.URL
			if from, ok := context.Params["from"]; ok {
				r = view.StringURL(from)
			} else {
				r = redirectURL
			}
			model := &LoginFormModel{}
			if email, ok := context.Params["email"]; ok {
				model.Email.Set(email)
			}
			form := &view.Form{
				Class:               class,
				ErrorMessageClass:   errorMessageClass,
				SuccessMessageClass: successMessageClass,
				SuccessMessage:      "Login successful",
				ButtonText:          "Login",
				FormID:              "gostart_user_login",
				GetModel:            view.FormModel(model),
				Redirect:            r,
				OnSubmit: func(form *view.Form, formModel interface{}, context *view.Context) (err error) {
					m := formModel.(*LoginFormModel)
					ok, err := LoginEmailPassword(context, m.Email.Get(), m.Password.Get())
					if err != nil {
						if view.Config.Debug.Mode {
							return err
						} else {
							return errors.New("An internal error ocoured")
						}
					}
					if !ok {
						return errors.New("Wrong email and password combination")
					}
					return nil
				},
			}
			return form, nil
		},
	)
}

///////////////////////////////////////////////////////////////////////////////
// LoginFormModel

type LoginFormModel struct {
	Email    model.Email
	Password model.Password
}
