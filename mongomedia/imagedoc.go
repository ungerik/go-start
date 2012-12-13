package mongomedia

import (
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
)

type ImageDoc struct {
	mongo.DocumentBase `bson:",inline"`
	media.Image        `bson:",inline"`
}

func (self *ImageDoc) Init(collection *mongo.Collection, embeddingStructPtr interface{}) {
	self.DocumentBase.Init(collection, embeddingStructPtr)
}

func (self *ImageDoc) Save() error {
	return self.DocumentBase.Save()
}

func (self *ImageDoc) Delete() error {
	return self.DocumentBase.Delete()
}
