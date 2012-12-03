package mongo

import (
	"fmt"
	"strconv"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/config"
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
)

func NewCollection(name string, documentLabelSelectors ...string) *Collection {
	if _, ok := Collections[name]; ok {
		panic(fmt.Sprintf("Collection '%s' already exists", name))
	}
	collection := &Collection{
		Name:              name,
		DocLabelSelectors: documentLabelSelectors,
		DocLabelSeparator: " ",
	}
	collection.Init()
	return collection
}

/*
Collection represents a MongoDB collection and implements mongo.Query for all
documents in the collection.

Example for creating, modifying and saving a document:

	user := models.Users.NewDocument().(*models.User)

	user.Name.First.Set("Erik")
	user.Name.Last.Set("Unger")

	err := user.Save()

*/
type Collection struct {
	query_base
	Name              string
	DocLabelSelectors []string
	DocLabelSeparator string
	collection        *mgo.Collection
	// foreignRefs  []ForeignRef
}

func (self *Collection) Init() {
	self.thisQuery = self
	Collections[self.Name] = self
	if Database != nil {
		self.collection = Database.C(self.Name)
	}
	// self.foreignRefs = []ForeignRef{}
}

func (self *Collection) checkDBConnection() {
	if self == nil {
		panic("mongo.Collection is nil")
	}
	if self.collection == nil {
		panic("mongo.Collection.collection is nil")
	}
	if self.collection.Database == nil {
		panic("mongo.Collection.collection.Database is nil")
	}
	if self.collection.Database.Session == nil {
		panic("mongo.Collection.collection.Database.Session is nil")
	}
}

func (self *Collection) String() string {
	return fmt.Sprintf("mongo.Collection '%s'", self.Name)
}

func (self *Collection) Selector() string {
	return ""
}

func (self *Collection) bsonSelector() bson.M {
	return bson.M{}
}

func (self *Collection) mongoQuery() (q *mgo.Query, err error) {
	self.checkDBConnection()
	return self.collection.Find(nil), nil
}

// func (self *Collection) subDocumentType(docType reflect.Type, fieldName string, subDocSelectors []string) (reflect.Type, error) {
// 	if fieldName == "" {
// 		return docType, errs.Format("Collection '%s', selector '%s': Empty field name", self.Name, strings.Join(subDocSelectors, "."))
// 	}

// 	switch docType.Kind() {
// 	case reflect.Struct:
// 		bsonName := strings.ToLower(fieldName)
// 		field := reflection.FindFlattenedStructField(docType, MatchBsonField(bsonName))
// 		if field != nil {
// 			return field.Type, nil
// 		}
// 		return nil, errs.Format("Collection '%s', selector '%s': Struct has no field '%s'", self.Name, strings.Join(subDocSelectors, "."), fieldName)

// 	case reflect.Array, reflect.Slice:
// 		_, numberErr := strconv.Atoi(fieldName)
// 		if numberErr == nil || fieldName == "$" {
// 			return docType, nil
// 		}
// 		return docType.Elem(), nil

// 	case reflect.Ptr, reflect.Interface:
// 		return self.subDocumentType(docType.Elem(), fieldName, subDocSelectors)
// 	}

// 	return nil, errs.Format("Collection '%s', selector '%s': Can't select sub-document '%s' of type '%s'", self.Name, strings.Join(subDocSelectors, "."), fieldName, docType.String())
// }

// todo
func (self *Collection) ValidateSelectors(selectors ...string) (err error) {
	// docType := self.DocumentType
	// for _, selector := range selectors {
	// 	fields := strings.Split(selector, ".")
	// 	if selector == "" || len(fields) == 0 {
	// 		return errs.Format("In`valid empty selector in '%s'", selector)
	// 	}
	// 	for _, field := range fields {
	// 		docType, err = self.subDocumentType(docType, field, selectors)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	return nil
}

func (self *Collection) InitDocument(doc Document) {
	doc.Init(self, doc)
}

func (self *Collection) InitAndSaveDocument(doc Document) error {
	self.InitDocument(doc)
	return doc.Save()
}

func (self *Collection) InitSubDocument(subDoc interface{}) {
	InitRefs(subDoc)
}

func (self *Collection) RefForID(id bson.ObjectId) Ref {
	return Ref{id, self.Name}
}

// betterError returns a better error description
func (self *Collection) betterError(err error, context string) error {
	if err == mgo.ErrNotFound {
		return fmt.Errorf("MongoDB document with id %s not found in collection %s", context, self.Collection().Name)
	}
	return err
}

func (self *Collection) logIdNotFoundError(id bson.ObjectId) {
	config.Logger.Printf("MongoDB document with id %s not found in collection %s", id.Hex(), self.Collection().Name)
	debug.LogCallStack()
}

func (self *Collection) documentWithID(id bson.ObjectId, resultRef interface{}) error {
	if id == "" {
		return errs.Format("mongo.Collection %s: Can't get document with empty id", self.Name)
	}
	self.checkDBConnection()
	q := self.collection.FindId(id)
	err := q.One(resultRef)
	if err != nil {
		return err
	}
	// resultRef has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.InitDocument(documentFromResultPtr(resultRef))
	return nil
}

func (self *Collection) DocumentWithID(id bson.ObjectId, resultRef interface{}) error {
	err := self.documentWithID(id, resultRef)
	if err == mgo.ErrNotFound {
		self.logIdNotFoundError(id)
	}
	return err
}

func (self *Collection) TryDocumentWithID(id bson.ObjectId, resultRef interface{}) (found bool, err error) {
	if id == "" {
		return false, nil
	}
	err = self.documentWithID(id, resultRef)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (self *Collection) SubDocumentWithID(id bson.ObjectId, selector string, resultRef interface{}) error {
	if id == "" {
		return errs.Format("mongo.Collection %s: Can't get document with empty id", self.Name)
	}
	self.checkDBConnection()
	q := self.collection.FindId(id).Select(bson.M{selector: 1})
	err := q.One(resultRef)
	if err != nil {
		return err
	}
	// resultRef has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.InitSubDocument(resultRef)
	return nil
}

func (self *Collection) TrySubDocumentWithID(id bson.ObjectId, selector string, resultRef interface{}) (found bool, err error) {
	err = self.SubDocumentWithID(id, selector, resultRef)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (self *Collection) UpdateSubDocumentWithID(id bson.ObjectId, selector string, subDocument interface{}) error {
	return self.Filter("_id", id).UpdateSubDocument(selector, subDocument)
}

type documentWithIDIterator struct {
	collection *Collection
	id         bson.ObjectId
	err        error
	nextCalled bool
}

func (self *documentWithIDIterator) Next(resultRef interface{}) (ok bool) {
	if self.nextCalled {
		return false
	}
	self.err = self.collection.DocumentWithID(self.id, resultRef)
	self.nextCalled = true
	return self.err == nil
}

func (self *documentWithIDIterator) Err() error {
	return self.err
}

type tryDocumentWithIDIterator struct {
	collection *Collection
	id         bson.ObjectId
	err        error
	nextCalled bool
}

func (self *tryDocumentWithIDIterator) Next(resultRef interface{}) (ok bool) {
	if self.nextCalled {
		return false
	}
	ok, self.err = self.collection.TryDocumentWithID(self.id, resultRef)
	self.nextCalled = true
	return ok
}

func (self *tryDocumentWithIDIterator) Err() error {
	return self.err
}

func (self *Collection) DocumentWithIDIterator(id bson.ObjectId) model.Iterator {
	return &documentWithIDIterator{collection: self, id: id}
}

func (self *Collection) TryDocumentWithIDIterator(id bson.ObjectId) model.Iterator {
	return &tryDocumentWithIDIterator{collection: self, id: id}
}

// FilterReferenced filters the collection for documents,
// that are referenced by refs.
func (self *Collection) FilterReferenced(refs ...Ref) Query {
	return self.FilterRef("_id", refs...)
}

func (self *Collection) Count() (n int, err error) {
	return self.collection.Count()
}

func (self *Collection) DeleteWithID(id bson.ObjectId) error {
	self.checkDBConnection()
	return self.collection.RemoveId(id)
}

func (self *Collection) DeleteAllWithIDs(ids ...bson.ObjectId) (numDeleted int, err error) {
	self.checkDBConnection()
	info, err := self.collection.RemoveAll(bson.M{"_id": bson.M{"$in": ids}})
	return info.Removed, err
}

func (self *Collection) DeleteAllNotWithIDs(ids ...bson.ObjectId) (numDeleted int, err error) {
	self.checkDBConnection()
	info, err := self.collection.RemoveAll(bson.M{"_id": bson.M{"$nin": ids}})
	return info.Removed, err
}

func (self *Collection) DocumentLabel(docId bson.ObjectId, labelSelectors ...string) (string, error) {
	if !docId.Valid() {
		return "", fmt.Errorf("Invalid bson.ObjectId")
	}
	if len(labelSelectors) == 0 {
		labelSelectors = self.DocLabelSelectors
	}

	if len(labelSelectors) == 0 {
		// No label selectors, use hex ID as label
		count, err := self.collection.FindId(docId).Count()
		if err != nil {
			return "", err
		}
		if count == 0 {
			return "", fmt.Errorf("No document with ID '%s' in collection '%s'", docId.Hex(), self.Name)
		}
		return docId.Hex(), nil
	}

	labels := make([]string, len(labelSelectors))
	for i, selector := range labelSelectors {
		if len(selector) == 0 {
			return "", fmt.Errorf("Label selector must not be an empty string")
		}

		selector = strings.ToLower(selector)

		var doc bson.M
		err := self.collection.FindId(docId).Select(bson.M{selector: 1}).One(&doc)
		if err != nil {
			return "", err
		}

		var label interface{}
		if strings.IndexRune(selector, '.') != -1 {
			for _, subSel := range strings.Split(selector, ".") {
				label, _ = doc[subSel]
				if label == nil {
					break
				}
				doc, _ = label.(bson.M)
				if doc == nil {
					break
				}
			}
		} else {
			label, _ = doc[selector]
		}

		if label != nil {
			// Support the basic BSON types from unmarshalling
			switch s := label.(type) {
			case string:
				labels[i] = s

			case int:
				labels[i] = strconv.Itoa(s)

			case float64:
				labels[i] = strconv.FormatFloat(s, 'f', -1, 64)

			case bool:
				labels[i] = strconv.FormatBool(s)

			case bson.ObjectId:
				labels[i] = s.Hex()

			default:
				return "", fmt.Errorf("Can't get label of %s for '%s' because it is of type %T", docId.Hex(), selector, label)
			}
		}
	}

	label := strings.Join(labels, self.DocLabelSeparator)
	label = strings.TrimSpace(utils.RemoveMultipleWhiteSpace(label))

	return label, nil
}

// RemoveInvalidRefs removes invalid refs from all documents and saves
// the changes.
func (self *Collection) RemoveInvalidRefs() (invalidRefs []InvalidRefData, err error) {
	config.Logger.Println("RemoveInvalidRefs() @ " + self.String())
	i := self.Iterator()
	var doc Document
	for i.Next(&doc) {
		refsData, err := doc.RemoveInvalidRefs()
		if err != nil {
			return nil, err
		}
		if len(refsData) > 0 {
			config.Logger.Println("Found invalid refs in document " + doc.ObjectId().Hex())
			for _, refData := range refsData {
				config.Logger.Printf("Invalid ref from %s to ", refData.MetaData.Selector(), refData.Ref)
			}
			err = doc.Save()
			if err != nil {
				config.Logger.Printf("mongo.Collection.RemoveInvalidRefs(): Error while saving document %#v", doc)
				return nil, err
			}
			invalidRefs = append(invalidRefs, refsData...)
		}
	}
	if i.Err() != nil {
		config.Logger.Printf("mongo.Collection.RemoveInvalidRefs(): Error while iterating collection documents: '%s'", i.Err())
		return nil, i.Err()
	}
	return invalidRefs, nil
}

///////////////////////////////////////////////////////////////////////////////
// ForeignRef

type ForeignRef struct {
	Collection *Collection
	Selector   string
}
