package user

import (
	"fmt"
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
	"github.com/ungerik/go-start/mgo/bson"
	"net/url"
)

///////////////////////////////////////////////////////////////////////////////
// EmailIdentity

type EmailIdentity struct {
	//	mongo.SubDocumentBase
	Address          model.Email
	Description      model.String
	Confirmed        model.DateTime
	ConfirmationCode model.String
}

// EmailIdentity has to be saved after a successful call because the confirmation code could have changed
// confirmationPage needs to be a page with one URL parameter
func (self *EmailIdentity) SendConfirmationEmail(context *view.Context, confirmationURL view.URL) <-chan error {
	errChan := make(chan error, 1)

	confirmationCode := self.ConfirmationCode.Get()
	if confirmationCode == "" {
		confirmationCode = bson.NewObjectId().Hex()
		self.ConfirmationCode.SetString(confirmationCode)
	}

	subject := fmt.Sprintf(Config.ConfirmationEmailSubject, view.Config.SiteName)
	confirm := confirmationURL.URL(context.PathArgs...) + "?code=" + url.QueryEscape(confirmationCode)
	message := fmt.Sprintf(Config.ConfirmationEmailMessage, view.Config.SiteName, confirm)

	go func() {
		errChan <- email.NewBriefMessage(subject, message, self.Address.Get()).Send()
		close(errChan)
	}()

	return errChan
}

func (self *EmailIdentity) MailtoURL() string {
	return "mailto:" + self.Address.Get()
}

func (self *EmailIdentity) URL(args ...string) string {
	return self.MailtoURL()
}

func (self *EmailIdentity) LinkContent(context *view.Context) view.View {
	return view.Escape(self.Address.Get())
}

func (self *EmailIdentity) LinkTitle(context *view.Context) string {
	return self.Address.Get()
}

func (self *EmailIdentity) LinkRel(context *view.Context) string {
	return ""
}
