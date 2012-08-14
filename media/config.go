package media

import (
// "github.com/ungerik/go-start/mongo"
// "github.com/ungerik/go-start/mgo"
// "errors"
)

var Config Configuration

type Configuration struct {
	Backend Backend
}

func (self *Configuration) Name() string {
	return "media"
}

func (self *Configuration) Init() error {
	return nil
}

func (self *Configuration) Close() error {
	return nil
}
