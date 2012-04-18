package user

import (
	"github.com/ungerik/go-start/view"
	// "net/url"
)

///////////////////////////////////////////////////////////////////////////////
// EmailConfirmationView

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

/*
// The confirmation code will be passed in the GET parameter "confirm"
type EmailConfirmationView struct {
	view.ViewBase

}

func (self *EmailConfirmationView) Render(context *view.Context, writer *utils.XMLWriter) (err error) {
	const invalidConfirmationCode = "<div class='error'>Invalid email confirmation code!</div>"

	confirmationCode, ok := context.Params["code"]
	if !ok {
		writer.Content(invalidConfirmationCode)
		return nil
	}

	email, confirmed, err := ConfirmEmail(confirmationCode)
	if err != nil {
		return err
	}
	if !confirmed {
		writer.Content(invalidConfirmationCode)
		return nil
	}

	writer.Printf("<div class='success'>Email address %s confirmed!</div>", email)
	if self.LoginURL != nil {
		loginURL := self.LoginURL.URL(context) + "?email=" + url.QueryEscape(email)
		writer.Printf("<p>You may <a href='%s'>login</a> now.</p>", loginURL)
	}

	return nil
}
*/
