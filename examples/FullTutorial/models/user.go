package models

import (
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/user"
)

// We are using user.NewCollection here instead of mongo.NewCollection
// because user.NewCollection sets the correct mongo.Collection.DocLabelSelectors
// so that mongo.Collection.DocumentLabel(id) returns a label for
// the document with id composed of the name modelext.Name components
// Prefix + First + Middle + Last + Postfix + Organization.
var Users = user.NewCollection("users")

func NewUser() *User {
	var doc User
	Users.InitDocument(&doc)
	return &doc
}

type User struct {
	user.User `bson:",inline"`

	Image  media.ImageRef
	Gender model.Choice `model:"options=,Male,Female"`
}
