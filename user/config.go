package user

import (
	"github.com/ungerik/go-start/mongo"
)

type Configuration struct {
	ConfirmationEmailSubject string
	ConfirmationEmailMessage string
	ConfirmationSent         string
	Collection               *mongo.Collection
}

var Config = Configuration{
	ConfirmationEmailSubject: "Please confirm your email address for %s",
	ConfirmationEmailMessage: "Please confirm your email address for %s by opening the following link:\n\n%s",
	ConfirmationSent:         "We sent you an email with a verification link.<br/>It might some time to show up, but when it does you will be ready to use this site.",
}

func Init(collection *mongo.Collection) {
	Config.Collection = collection
}
