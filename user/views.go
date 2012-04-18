package user

import (
	"errors"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
)

// The confirmation code will be passed in the GET parameter "code"
func EmailConfirmationView(profileURL view.URL) view.View {
	return view.DynamicView(
		func(context *view.Context) (view.View, error) {
			confirmationCode, ok := context.Params["code"]
			if !ok {
				return view.DIV("error", view.HTML("Invalid email confirmation code!")), nil
			}

			doc, email, confirmed, err := ConfirmEmail(confirmationCode)
			if !confirmed {
				return view.DIV("error", view.HTML("Invalid email confirmation code!")), err
			}

			Login(context, doc)

			return view.Views{
				view.DIV("success", view.Printf("Email address %s confirmed!", email)),
				&view.If{
					Condition: profileURL != nil,
					Content: view.P(
						view.HTML("Continue to your "),
						view.A(profileURL, "profile..."),
					),
				},
			}, nil
		},
	)
}

func NewLoginForm(buttonText, class, errorMessageClass, successMessageClass string, redirectURL view.URL) view.View {
	return view.DynamicView(
		func(context *view.Context) (v view.View, err error) {
			if from, ok := context.Params["from"]; ok {
				redirectURL = view.StringURL(from)
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
				ButtonText:          buttonText,
				FormID:              "gostart_user_login",
				GetModel:            view.FormModel(model),
				Redirect:            redirectURL,
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

type LoginFormModel struct {
	Email    model.Email
	Password model.Password
}

// If redirect is nil, the redirect will go to "/"
func LogoutView(redirect view.URL) view.View {
	return view.RenderView(
		func(context *view.Context, writer *utils.XMLWriter) (err error) {
			Logout(context)
			if redirect != nil {
				return view.Redirect(redirect.URL(context))
			}
			return view.Redirect("/")
		},
	)
}

// confirmationPage must have the confirmation code as first URL parameter
func NewSignupForm(buttonText, class, errorMessageClass, successMessageClass string, confirmationURL, redirectURL view.URL) *view.Form {
	return &view.Form{
		Class:               class,
		ErrorMessageClass:   errorMessageClass,
		SuccessMessageClass: successMessageClass,
		SuccessMessage:      Config.ConfirmationSent,
		ButtonText:          buttonText,
		FormID:              "gostart_user_signup",
		GetModel: func(form *view.Form, context *view.Context) (interface{}, error) {
			return &EmailPasswordFormModel{}, nil
		},
		Redirect: redirectURL,
		OnSubmit: func(form *view.Form, formModel interface{}, context *view.Context) error {
			m := formModel.(*EmailPasswordFormModel)
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
