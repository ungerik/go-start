package mongomedia

import (
	"io"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mgo"
	"github.com/ungerik/go-start/mgo/bson"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
)

type backend struct {
	gridFS                      *mgo.GridFS
	images                      *mongo.Collection
	imageRefCollectionSelectors map[*mongo.Collection][]string
}

func (self *backend) LoadImage(id string) (*media.Image, error) {
	doc, err := self.images.DocumentWithID(bson.ObjectIdHex(id))
	if err != nil {
		return nil, err
	}
	return doc.(*ImageDoc).GetAndInitImage(), nil
}

func (self *backend) SaveImage(image *media.Image) error {
	if image.ID == "" {
		doc := self.images.NewDocument().(*ImageDoc)
		doc.Image = *image
		id, err := self.images.Insert(doc)
		if err != nil {
			return err
		}
		image.ID.Set(id.Hex())
		doc.Image.ID = image.ID
		return self.images.Update(id, doc)
	}

	id := bson.ObjectIdHex(image.ID.Get())
	doc := self.images.NewDocument().(*ImageDoc)
	doc.SetObjectId(id)
	doc.Image = *image
	return self.images.Update(id, doc)
}

func (self *backend) DeleteImage(image *media.Image) error {
	for i := range image.Versions {
		err := self.DeleteImageVersion(image.Versions[i].ID.Get())
		if err != nil {
			return err
		}
	}
	return self.images.Remove(bson.ObjectIdHex(image.ID.Get()))
}

func (self *backend) DeleteImageVersion(id string) error {
	return self.gridFS.RemoveId(id)
}

func (self *backend) ImageVersionReader(id string) (reader io.ReadCloser, ctype string, err error) {
	file, err := self.gridFS.OpenId(bson.ObjectIdHex(id))
	if err == mgo.NotFound {
		return nil, "", media.ErrInvalidImageID(id)
	} else if err != nil {
		return nil, "", err
	}
	return file, file.ContentType(), nil
}

func (self *backend) ImageVersionWriter(version *media.ImageVersion) (writer io.WriteCloser, err error) {
	if version.ID != "" {
		err = self.gridFS.RemoveId(version.ID)
		if err != nil {
			return nil, err
		}
		version.ID = ""
	}
	file, err := self.gridFS.Create("")
	if err != nil {
		return nil, err
	}
	id := file.Id().(bson.ObjectId).Hex()
	file.SetName(id + "/" + version.Filename.Get())
	file.SetMeta(version)
	version.ID.Set(id)
	return file, err
}

func (self *backend) ImageIterator() model.Iterator {
	return model.ConvertIterator(self.images.Iterator(),
		func(doc interface{}) interface{} {
			return doc.(*ImageDoc).GetAndInitImage()
		},
	)
}

func (self *backend) getImageRefCollectionSelectors() map[*mongo.Collection][]string {
	if self.imageRefCollectionSelectors == nil {
		colSel := make(map[*mongo.Collection][]string)
		for _, collection := range mongo.Collections {
			doc := collection.NewDocument()
			model.SetAllSliceLengths(doc, 1)
			model.Visit(doc, model.FieldOnlyVisitor(
				func(data *model.MetaData) error {
					if _, ok := data.Value.Addr().Interface().(*media.ImageRef); ok {
						if _, ok := colSel[collection]; !ok {
							colSel[collection] = nil
						}
						colSel[collection] = append(colSel[collection], data.WildcardSelector())
					}
					return nil
				},
			))
		}
		self.imageRefCollectionSelectors = colSel
	}
	return self.imageRefCollectionSelectors
}

func (self *backend) CountImageRefs(imageID string) (int, error) {
	colSel := self.getImageRefCollectionSelectors()
	// debug.Dump(colSel)
	count := 0
	for collection, selectors := range colSel {
		for _, sel := range selectors {
			c, err := collection.Filter(sel, imageID).Count()
			if err != nil {
				return 0, err
			}
			count += c
		}
	}
	return count, nil
}
