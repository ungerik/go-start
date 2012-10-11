package mongomedia

import (
	"io"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
)

type Backend struct {
	GridFS                      *mgo.GridFS
	Images                      *mongo.Collection
	imageRefCollectionSelectors map[*mongo.Collection][]string
}

func (self *Backend) LoadImage(id string) (*media.Image, error) {
	var doc ImageDoc
	err := self.Images.DocumentWithID(bson.ObjectIdHex(id), &doc)
	if err != nil {
		return nil, err
	}
	return doc.GetAndInitImage(), nil
}

func (self *Backend) TryLoadImage(id string) (*media.Image, bool, error) {
	var doc ImageDoc
	found, err := self.Images.TryDocumentWithID(bson.ObjectIdHex(id), &doc)
	if !found {
		return nil, found, err
	}
	return doc.GetAndInitImage(), true, nil
}

func (self *Backend) SaveImage(image *media.Image) error {
	if image.ID == "" {
		var imageDoc ImageDoc
		imageDoc.Image = *image
		err := imageDoc.Save()
		if err != nil {
			return err
		}
		imageDoc.Image.ID.Set(imageDoc.ObjectId().Hex())
		return imageDoc.Save()
	}

	var imageDoc ImageDoc
	imageDoc.SetObjectId(bson.ObjectIdHex(image.ID.Get()))
	imageDoc.Image = *image
	return self.Images.InitAndSaveDocument(&imageDoc)
}

func (self *Backend) DeleteImage(image *media.Image) error {
	for i := range image.Versions {
		err := self.DeleteImageVersion(image.Versions[i].ID.Get())
		if err != nil {
			return err
		}
	}
	return self.Images.Delete(bson.ObjectIdHex(image.ID.Get()))
}

func (self *Backend) DeleteImageVersion(id string) error {
	return self.GridFS.RemoveId(id)
}

func (self *Backend) ImageVersionReader(id string) (reader io.ReadCloser, ctype string, err error) {
	file, err := self.GridFS.OpenId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return nil, "", media.ErrInvalidImageID(id)
	} else if err != nil {
		return nil, "", err
	}
	return file, file.ContentType(), nil
}

func (self *Backend) ImageVersionWriter(version *media.ImageVersion) (writer io.WriteCloser, err error) {
	if version.ID != "" {
		err = self.GridFS.RemoveId(version.ID)
		if err != nil {
			return nil, err
		}
		version.ID = ""
	}
	file, err := self.GridFS.Create("")
	if err != nil {
		return nil, err
	}
	id := file.Id().(bson.ObjectId).Hex()
	file.SetName(id + "/" + version.Filename.Get())
	file.SetMeta(version)
	version.ID.Set(id)
	return file, err
}

func (self *Backend) ImageIterator() model.Iterator {
	return model.ConversionIterator(self.Images.Iterator(),
		func(doc interface{}) interface{} {
			return doc.(*ImageDoc).GetAndInitImage()
		},
	)
}

func (self *Backend) getImageRefCollectionSelectors() map[*mongo.Collection][]string {
	if self.imageRefCollectionSelectors == nil {
		colSel := make(map[*mongo.Collection][]string)
		for _, collection := range mongo.Collections {
			var doc ImageDoc
			collection.InitDocument(&doc)
			model.SetAllSliceLengths(&doc, 1)
			model.Visit(&doc, model.FieldOnlyVisitor(
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

func (self *Backend) CountImageRefs(imageID string) (count int, err error) {
	colSel := self.getImageRefCollectionSelectors()
	// debug.Dump(colSel)
	for collection, selectors := range colSel {
		for _, selector := range selectors {
			c, err := collection.Filter(selector, imageID).Count()
			if err != nil {
				return 0, err
			}
			count += c
		}
	}
	return count, nil
}

func (self *Backend) RemoveAllImageRefs(imageID string) (count int, err error) {
	colSel := self.getImageRefCollectionSelectors()
	for collection, selectors := range colSel {
		for _, selector := range selectors {
			c, err := collection.Filter(selector, imageID).UpdateAll(selector, "")
			if err != nil {
				return count, err
			}
			count += c
		}
	}
	return count, nil
}
