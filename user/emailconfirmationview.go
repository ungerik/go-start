package user

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
	"net/url"
)

///////////////////////////////////////////////////////////////////////////////
// EmailConfirmationView

// The confirmation code will be passed in the GET parameter "confirm"
type EmailConfirmationView struct {
	view.ViewBase
	LoginURL view.URL
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
