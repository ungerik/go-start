package media

import (
	"github.com/ungerik/go-start/mongo"
	"launchpad.net/mgo"
	"errors"
)

var gridFS *mgo.GridFS

// Init must be called after mongo.Init()
func Init(name string) error {
	if name == "" {
		return errors.New("media.Init() called with empty name")
	}
	gridFS = mongo.Database.GridFS(name)
	return nil
}
