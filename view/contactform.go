package view

import (
	"fmt"
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/model"
)

// ContactFormModel is a default form model for contact forms.
type ContactFormModel struct {
	Name    model.String `view:"label=Your name" model:"maxlen=40"`
	Email   model.Email  `view:"label=Your email address" model:"required|maxlen=40"`
	Subject model.String `view:"label=Subject" model:"maxlen=40"`
	Message model.Text   `view:"label=Your message|cols=40|rows=10" model:"required"`
}

// NewContactForm creates a new contact form that sends submitted data to recipientEmail.
func NewContactForm(recipientEmail, subjectPrefix, formClass, buttonClass, formID string) *Form {
	return &Form{
		Class:             formClass,
		SubmitButtonClass: buttonClass,
		SubmitButtonText:  "Send",
		SuccessMessage:    "Message sent",
		FormID:            formID,
		GetModel: func(form *Form, ctx *Context) (interface{}, error) {
			return &ContactFormModel{}, nil
		},
		OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
			model := formModel.(*ContactFormModel)
			subject := fmt.Sprintf("%sFrom %s <%s>: %s", subjectPrefix, model.Name, model.Email, model.Subject)
			err := email.NewBriefMessage(subject, model.Message.Get(), recipientEmail).Send()
			return "", nil, err
		},
	}
}
