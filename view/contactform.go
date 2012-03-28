package view

import (
	"fmt"
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/model"
)

// ContactFormModel is a default form model for contact forms.
type ContactFormModel struct {
	Name    model.String `gostart:"label=Your name|maxlen=40"`
	Email   model.Email  `gostart:"label=Your email address|required|maxlen=40"`
	Subject model.String `gostart:"label=Subject|maxlen=40"`
	Message model.Text   `gostart:"label=Your message|required|cols=40|rows=10"`
}

// NewContactForm creates a new contact form that sends submitted data to recipientEmail.
func NewContactForm(recipientEmail, subjectPrefix, formClass, buttonClass, formID string) *Form {
	return &Form{
		Class:          formClass,
		ButtonClass:    buttonClass,
		ButtonText:     "Send",
		SuccessMessage: "Message sent",
		FormID:         formID,
		GetModel: func(form *Form, context *Context) (interface{}, error) {
			return &ContactFormModel{}, nil
		},
		OnSubmit: func(form *Form, formModel interface{}, context *Context) error {
			model := formModel.(*ContactFormModel)
			subject := fmt.Sprintf("%sFrom %s <%s>: %s", subjectPrefix, model.Name, model.Email, model.Subject)
			return email.NewBriefMessage(subject, model.Message.Get(), recipientEmail).Send()
		},
	}
}
