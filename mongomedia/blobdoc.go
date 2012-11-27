package mongomedia

import (
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
)

type BlobDoc struct {
	mongo.DocumentBase `bson:",inline"`
	media.Blob         `bson:",inline"`
}

func (self *BlobDoc) Init(collection *mongo.Collection, embeddingStructPtr interface{}) {
	self.DocumentBase.Init(collection, embeddingStructPtr)
}

func (self *BlobDoc) Save() error {
	return self.DocumentBase.Save()
}

func (self *BlobDoc) Delete() error {
	return self.DocumentBase.Delete()
}
