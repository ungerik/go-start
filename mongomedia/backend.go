package mongomedia

import (
	"io"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
)

type Backend struct {
	GridFS *mgo.GridFS
	Images *mongo.Collection
	Blobs  *mongo.Collection

	imageRefCollectionSelectors map[*mongo.Collection][]string
	blobRefCollectionSelectors  map[*mongo.Collection][]string
}

func (self *Backend) Init(gridFSName string) {
	self.GridFS = mongo.Database.GridFS(gridFSName)
	self.Images = mongo.NewCollection(gridFSName+".images", "Title")
	self.Blobs = mongo.NewCollection(gridFSName+".blobs", "Title")
}

func (self *Backend) FileWriter(filename, contentType string) (writer io.WriteCloser, id string, err error) {
	file, err := self.GridFS.Create("")
	if err != nil {
		return nil, "", err
	}
	id = file.Id().(bson.ObjectId).Hex()
	file.SetName(id + "/" + filename)
	file.SetContentType(contentType)
	return file, id, err
}

func (self *Backend) FileReader(id string) (reader io.ReadCloser, filename, contentType string, err error) {
	file, err := self.GridFS.OpenId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return nil, "", "", media.ErrNotFound(id)
	} else if err != nil {
		return nil, "", "", err
	}
	filename = id[strings.IndexRune(id, '/')+1:]
	return file, filename, file.ContentType(), nil
}

func (self *Backend) DeleteFile(id string) error {
	return self.GridFS.RemoveId(bson.ObjectIdHex(id))
}

func (self *Backend) LoadBlob(id string) (*media.Blob, error) {
	var doc BlobDoc
	err := self.Blobs.DocumentWithID(bson.ObjectIdHex(id), &doc)
	if err == mgo.ErrNotFound {
		return nil, media.ErrNotFound(id)
	} else if err != nil {
		return nil, err
	}
	return doc.Blob.Init(), nil
}

func (self *Backend) SaveBlob(blob *media.Blob) error {
	if blob.ID == "" {
		blob.ID.Set(bson.NewObjectId().Hex())
	}
	var blobDoc BlobDoc
	blobDoc.SetObjectId(bson.ObjectIdHex(blob.ID.Get()))
	blobDoc.Blob = *blob
	return self.Blobs.InitAndSaveDocument(&blobDoc)
}

func (self *Backend) DeleteBlob(blob *media.Blob) error {
	return self.Blobs.DeleteWithID(bson.ObjectIdHex(blob.ID.Get()))
}

func (self *Backend) BlobIterator() model.Iterator {
	return model.ConversionIterator(
		self.Blobs.Iterator(),
		new(BlobDoc),
		func(doc interface{}) interface{} {
			return &doc.(*BlobDoc).Blob
		},
	)
}

func (self *Backend) getBlobRefCollectionSelectors() map[*mongo.Collection][]string {
	return nil

	if self.imageRefCollectionSelectors == nil {
		colSel := make(map[*mongo.Collection][]string)
		for _, collection := range mongo.Collections {
			var doc BlobDoc // todo needs to be document type of collection
			collection.InitDocument(&doc)
			model.SetAllSliceLengths(&doc, 1)
			model.Visit(&doc, model.FieldOnlyVisitor(
				func(data *model.MetaData) error {
					if _, ok := data.Value.Addr().Interface().(*media.BlobRef); ok {
						if _, ok := colSel[collection]; !ok {
							colSel[collection] = nil
						}
						colSel[collection] = append(colSel[collection], data.WildcardSelector())
					}
					return nil
				},
			))
		}
		self.blobRefCollectionSelectors = colSel
	}
	return self.blobRefCollectionSelectors
}

func (self *Backend) CountBlobRefs(blobID string) (count int, err error) {
	colSel := self.getBlobRefCollectionSelectors()
	for collection, selectors := range colSel {
		for _, selector := range selectors {
			c, err := collection.Filter(selector, blobID).Count()
			if err != nil {
				return 0, err
			}
			count += c
		}
	}
	return count, nil
}

func (self *Backend) RemoveAllBlobRefs(blobID string) (count int, err error) {
	colSel := self.getImageRefCollectionSelectors()
	for collection, selectors := range colSel {
		for _, selector := range selectors {
			c, err := collection.Filter(selector, blobID).UpdateAll(selector, "")
			if err != nil {
				return count, err
			}
			count += c
		}
	}
	return count, nil
}

func (self *Backend) LoadImage(id string) (*media.Image, error) {
	var doc ImageDoc
	err := self.Images.DocumentWithID(bson.ObjectIdHex(id), &doc)
	if err == mgo.ErrNotFound {
		return nil, media.ErrNotFound(id)
	} else if err != nil {
		return nil, err
	}
	return doc.Image.Init(), nil
}

func (self *Backend) SaveImage(image *media.Image) error {
	if image.ID == "" {
		image.ID.Set(bson.NewObjectId().Hex())
	}
	var imageDoc ImageDoc
	imageDoc.SetObjectId(bson.ObjectIdHex(image.ID.Get()))
	imageDoc.Image = *image
	return self.Images.InitAndSaveDocument(&imageDoc)
}

func (self *Backend) DeleteImage(image *media.Image) error {
	for i := range image.Versions {
		err := self.DeleteFile(image.Versions[i].ID.Get())
		if err != nil {
			return err
		}
	}
	return self.Images.DeleteWithID(bson.ObjectIdHex(image.ID.Get()))
}

// func (self *Backend) ImageVersionWriter(version *media.ImageVersion) (writer io.WriteCloser, err error) {
// 	if version.ID != "" {
// 		err = self.GridFS.RemoveId(bson.ObjectIdHex(version.ID.Get()))
// 		if err != nil {
// 			return nil, err
// 		}
// 		version.ID = ""
// 	}
// 	file, err := self.GridFS.Create("")
// 	if err != nil {
// 		return nil, err
// 	}
// 	id := file.Id().(bson.ObjectId).Hex()
// 	file.SetName(id + "/" + version.Filename.Get())
// 	file.SetMeta(version)
// 	version.ID.Set(id)
// 	return file, err
// }

func (self *Backend) ImageIterator() model.Iterator {
	return model.ConversionIterator(
		self.Images.Iterator(),
		new(ImageDoc),
		func(doc interface{}) interface{} {
			return doc.(*ImageDoc).Image.Init()
		},
	)
}

func (self *Backend) getImageRefCollectionSelectors() map[*mongo.Collection][]string {
	return nil

	if self.imageRefCollectionSelectors == nil {
		colSel := make(map[*mongo.Collection][]string)
		for _, collection := range mongo.Collections {
			var doc ImageDoc // todo needs to be document type of collection
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
