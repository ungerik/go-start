package mongomedia

import (
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
)

type Image struct {
	mongo.DocumentBase `bson:",inline"`
	media.Image        `bson:",inline"`
}
