package models

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/user"
)

var People = mongo.NewCollection("people", (*Person)(nil))

///////////////////////////////////////////////////////////////////////////////
// Person

type Person struct {
	user.User  `bson:",inline"`
	SuperAdmin model.Bool
	Gender     model.Choice `gostart:"options=,Male,Female|required"`
}

func (self *Person) String() string {
	return self.User.Name.String()
}
