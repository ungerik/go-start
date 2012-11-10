package user

import (
	"errors"

	"github.com/ungerik/go-start/view"
)

// The confirmation code will be passed in the GET parameter "code"
func EmailConfirmationView(profileURL view.URL) view.View {
	return view.DynamicView(
		func(ctx *view.Context) (view.View, error) {

			confirmationCode, ok := ctx.Request.Params["code"]
			if !ok {
				return view.DIV("error", view.HTML("Invalid email confirmation code!")), nil
			}

			userID, email, confirmed, err := ConfirmEmail(confirmationCode)
			if !confirmed {
				return view.DIV("error", view.HTML("Invalid email confirmation code!")), err
			}

			LoginID(ctx.Session, userID)

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
		func(ctx *view.Context) (v view.View, err error) {
			if from, ok := ctx.Request.Params["from"]; ok {
				redirectURL = view.StringURL(from)
			}
			model := &LoginFormModel{}
			if email, ok := ctx.Request.Params["email"]; ok {
				model.Email.Set(email)
			}
			form := &view.Form{
				Class:               class,
				ErrorMessageClass:   errorMessageClass,
				SuccessMessageClass: successMessageClass,
				SuccessMessage:      "Login successful",
				SubmitButtonText:    buttonText,
				FormID:              "gostart_user_login",
				GetModel:            view.FormModel(model),
				OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (string, view.URL, error) {
					m := formModel.(*LoginFormModel)
					ok, err := LoginEmailPassword(ctx.Session, m.Email.Get(), m.Password.Get())
					if err != nil {
						if view.Config.Debug.Mode {
							return "", nil, err
						} else {
							return "", nil, errors.New("An internal error ocoured")
						}
					}
					if !ok {
						return "", nil, errors.New("Wrong email and password combination")
					}
					return "", redirectURL, nil
				},
			}
			return form, nil
		},
	)
}

// If redirect is nil, the redirect will go to "/"
func LogoutView(redirect view.URL) view.View {
	return view.RenderView(
		func(ctx *view.Context) (err error) {
			Logout(ctx.Session)
			if redirect != nil {
				return view.Redirect(redirect.URL(ctx))
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
		SuccessMessage:      Config.ConfirmationMessage.Sent,
		SubmitButtonText:    buttonText,
		FormID:              "gostart_user_signup",
		GetModel: func(form *view.Form, ctx *view.Context) (interface{}, error) {
			return &EmailPasswordFormModel{}, nil
		},
		OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (string, view.URL, error) {
			m := formModel.(*EmailPasswordFormModel)
			email := m.Email.Get()
			password := m.Password1.Get()
			var user User
			found, err := WithEmail(email, &user)
			if err != nil {
				return "", nil, err
			}
			if found {
				if user.EmailPasswordConfirmed() {
					return "", nil, errors.New("A user with that email and a password already exists")
				}
				user.Password.SetHashed(password)
			} else {
				// Config.Collection.InitDocument(&user)
				err = user.SetEmailPassword(email, password)
				if err != nil {
					return "", nil, err
				}
			}
			err = <-user.Email[0].SendConfirmationEmail(ctx, confirmationURL)
			if err != nil {
				return "", nil, err
			}

			if found {
				err = Config.Collection.UpdateSubDocumentWithID(user.ID, "", &user)
			} else {
				err = Config.Collection.InitAndSaveDocument(&user)
			}
			return "", redirectURL, err
		},
	}
}
