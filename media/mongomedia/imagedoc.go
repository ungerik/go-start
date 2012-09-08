package mongomedia

import (
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
)

type ImageDoc struct {
	mongo.DocumentBase `bson:",inline"`
	media.Image        `bson:",inline"`
}

func (self *ImageDoc) GetAndInitImage() *media.Image {
	self.Image.Init()
	return &self.Image
}
