package models

import (
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/user"
)

var Users = mongo.NewCollection("users")

func NewUser() *User {
	var doc User
	Users.InitDocument(&doc)
	return &doc
}

type User struct {
	user.User `bson:",inline"`

	Image  media.ImageRef
	Gender model.Choice `model:"options=Male,Female"`
}
