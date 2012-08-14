package mongomedia

import (
	"errors"
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
)

var Config = Configuration{
	GridFSName: "media",
}

type Configuration struct {
	GridFSName string
}

func (self *Configuration) Name() string {
	return "mongomedia"
}

func (self *Configuration) Init() error {
	media.Config.Backend = &Backend{
		gridFS: mongo.Database.GridFS(self.GridFSName),
		images: mongo.NewCollection(self.GridFSName+".images", (*ImageDoc)(nil)),
	}
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
