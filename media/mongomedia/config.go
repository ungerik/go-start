package mongomedia

import (
	"errors"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
)

var Config = Configuration{
	GridFSName: "media",
}

type Configuration struct {
	GridFSName string
	Backend    Backend
}

func (self *Configuration) Name() string {
	return "mongomedia"
}

func (self *Configuration) Init() error {
	if mongo.Database == nil {
		panic("Package mongo must be initialized before mongomedia")
	}
	self.Backend.Init(self.GridFSName)
	media.Config.Backend = &self.Backend
	return nil
}

func (self *Configuration) Close() error {
	return nil
}

// Init must be called after mongo.Init()
func Init(name string) error {
	if name == "" {
		return errors.New("media.Init() called with empty name")
	}
	Config.GridFSName = name
	return Config.Init()
}
