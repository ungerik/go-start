package user

import (
	"fmt"

	"github.com/ungerik/go-start/mongo"
)

var Config = Configuration{
	ConfirmationMessage: ConfirmationMessage{
		EmailSubject: "Please confirm your email address for %s",
		EmailMessage: "Please confirm your email address for %s by opening the following link:\n\n%s",
		Sent:         "We sent you an email with a verification link. It might some time to show up, but when it does you will be ready to use this site.",
	},
}

func Init(collection *mongo.Collection) {
	Config.Collection = collection
	Config.CollectionName = collection.Name
}

type ConfirmationMessage struct {
	EmailSubject string
	EmailMessage string
	Sent         string
}

type Configuration struct {
	ConfirmationMessage ConfirmationMessage
	CollectionName      string
	Collection          *mongo.Collection
}

func (self *Configuration) Name() string {
	return "user"
}

func (self *Configuration) Init() error {
	collection, found := mongo.CollectionByName(self.CollectionName)
	if !found {
		return fmt.Errorf("Can't find mongo collection with name '%s'", self.CollectionName)
	}
	self.Collection = collection
	return nil
}

func (self *Configuration) Close() error {
	return nil
}
