package mongomedia

import (
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/media"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"errors"
	"io"
)

// Init must be called after mongo.Init()
func Init(name string) error {
	if name == "" {
		return errors.New("media.Init() called with empty name")
	}
	media.Config.Backend = &Backend{
		gridFS: mongo.Database.GridFS(name),
		images: mongo.NewCollection(name+".images", (*Image)(nil)),
	}
	return nil
}

type Backend struct {
	gridFS *mgo.GridFS
	images *mongo.Collection
}

func (self *Backend) Image(id string) (*media.Image, error) {
	doc, err := self.images.DocumentWithID(bson.ObjectIdHex(id))
	if err != nil {
		return nil, err
	}
	return &doc.(*Image).Image, nil
}

func (self *Backend) LoadImage(id string) (image *media.BackendImage, found bool, err error) {
	file, err := self.gridFS.OpenId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	image = &media.BackendImage{Reader: file}
	err = file.GetMeta(image)
	if err != nil {
		return nil, false, err
	}
	return image, true, nil
}

func (self *Backend) SaveImage(id string, image *media.BackendImage) (err error) {
	file, err := self.gridFS.OpenId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}
	defer image.Reader.Close()
	_, err = io.Copy(file, image.Reader)
	if err != nil {
		return err
	}
	file.SetMeta(image)
	return nil
}

func (self *Backend) SaveNewImage(image *media.BackendImage) (id string, err error) {
	file, err := self.gridFS.Create("")
	if err != nil {
		return "", err
	}
	defer image.Reader.Close()
	_, err = io.Copy(file, image.Reader)
	if err != nil {
		return "", err
	}
	file.SetMeta(image)
	return file.Id().(bson.ObjectId).Hex(), nil
}
