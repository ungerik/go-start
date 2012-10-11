package mongo

import (
	"reflect"
	"sort"
	"strings"

	"github.com/ungerik/go-start/config"
	// "github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/reflection"
	"github.com/ungerik/go-start/utils"
	"labix.org/v2/mgo/bson"

//	"github.com/ungerik/go-start/debug"
)

// func AddCollection(collection *Collection) (err error) {
// 	if _, ok := collections[collection.Name]; ok {
// 		return errs.Format("Collection %s already added", collection.Name)
// 	}
// 	collection.Init()
// 	collections[collection.Name] = collection
// 	return nil
// }

func documentFromResultPtr(resultPtr interface{}) Document {
	if resultPtr == nil {
		return nil
	}
	if doc, ok := resultPtr.(Document); ok {
		return doc
	}
	return documentFromResultPtr(reflect.ValueOf(resultPtr).Elem().Interface())
}

func CollectionByName(name string) (collection *Collection, ok bool) {
	collection, ok = Collections[name]
	return collection, ok
}

const javascriptRegexSpecialChars string = "\\^$.?*!+:=()[]{}|"

func escapeStringForRegex(str string) string {
	for i := range javascriptRegexSpecialChars {
		c := javascriptRegexSpecialChars[i]
		str = strings.Replace(str, string(c), string([]byte{'\\', c}), -1)
	}
	str = strings.Replace(str, "\t", "\\t", -1)
	str = strings.Replace(str, "\v", "\\v", -1)
	str = strings.Replace(str, "\f", "\\f", -1)
	str = strings.Replace(str, "\n", "\\n", -1)
	str = strings.Replace(str, "\r", "\\r", -1)
	return str
}

func collectionAndSubDocumentSelectors(query Query) (collection *Collection, selectors []string) {
	selectors = []string{}
	for collection == nil {
		collection, _ = query.(*Collection)
		selector := query.subDocumentSelector()
		if selector != "" {
			selectors = append(selectors, selector)
		}
		query = query.ParentQuery()
	}
	utils.ReverseStringSlice(selectors)
	return collection, selectors
}

func checkQuery(query Query) Query {
	if Config.CheckQuerySelectors {
		collection, selectors := collectionAndSubDocumentSelectors(query)
		selector := query.Selector()
		if selector != "" {
			selectors = append(selectors, selector)
		}
		if err := collection.ValidateSelectors(selectors...); err != nil {
			return &QueryError{query.ParentQuery(), err}
		}
	}
	return query
}

func MatchBsonField(bsonName string) reflection.MatchStructFieldFunc {
	return func(field *reflect.StructField) bool {
		var name string
		bsonTag := field.Tag.Get("bson")
		if bsonTag != "" {
			name = strings.SplitN(bsonTag, ",", 1)[0]
		} else {
			name = strings.ToLower(field.Name)
		}
		return name == bsonName
	}
}

func RemoveRefWithIDFromSlice(refs []Ref, id bson.ObjectId) ([]Ref, bool) {
	if i, found := FindRefWithID(refs, id); found {
		return append(refs[0:i], refs[i+1:len(refs)]...), true
	}
	return refs, false
}

// nil refs will be ignored
func CheckRefs(refs []Ref) (validRefs []Ref, invalidRefs []Ref, err error) {
	for _, ref := range refs {
		if !ref.IsEmpty() {
			if ok, err := ref.CheckID(); !ok {
				if err != nil {
					return nil, nil, err
				}
				invalidRefs = append(invalidRefs, ref)
			} else {
				validRefs = append(validRefs, ref)
			}
		}
	}
	return validRefs, invalidRefs, nil
}

func FindRefWithID(refs []Ref, id bson.ObjectId) (index int, found bool) {
	for i := range refs {
		if refs[i].ID == id {
			return i, true
		}
	}
	return -1, false
}

type SortableRefs struct {
	Refs     []Ref
	LessFunc func(a, b *Ref) bool
}

func (self *SortableRefs) Len() int {
	return len(self.Refs)
}

func (self *SortableRefs) Less(i, j int) bool {
	return self.LessFunc(&self.Refs[i], &self.Refs[j])
}

func (self *SortableRefs) Swap(i, j int) {
	self.Refs[i], self.Refs[j] = self.Refs[j], self.Refs[i]
}

func SortRefs(refs []Ref, lessFunc func(a, b *Ref) bool) {
	sort.Sort(&SortableRefs{refs, lessFunc})
}

//func SortQueryByRef(query Query, selector string) (model.Iterator, error) {
//	collection, selectors := collectionAndSubDocumentSelectors(query)
//	selectors = append(selectors, selector)
//	if err := collection.ValidateSelector(selectors...); err != nil {
//		return nil, err
//	}
//
//	// todo
//
//	docs := []interface{}{}
//	//refMap := make(map[bson.ObjectId]interface{})
//
//	i := query.Iterator()
//	for doc := i.Next(); doc != nil; doc = i.Next() {
//
//	}
//
//	return model.NewSliceIterator(docs...), nil
//}

func InitRefs(document interface{}) {
	model.Visit(document, model.FieldOnlyVisitor(
		func(data *model.MetaData) error {
			if data.Value.CanAddr() {
				if ref, ok := data.Value.Addr().Interface().(*Ref); ok && ref.CollectionName == "" {
					ref.CollectionName, ok = data.Attrib(model.StructTagKey, "to")
					if !ok {
						panic(data.Selector() + " is missing the 'to' meta-data tag")
					}
				}
			}
			return nil
		},
	))
}

type InvalidRefData struct {
	Ref      *Ref
	MetaData *model.MetaData
}

func RemoveInvalidRefs(document interface{}) (invalidRefs []InvalidRefData, err error) {
	err = model.Visit(document, model.FieldTypeVisitor(
		func(ref *Ref, metaData *model.MetaData) error {
			ok, err := ref.CheckID()
			if err != nil {
				return err
			}
			if !ok {
				invalidRefs = append(invalidRefs, InvalidRefData{ref, metaData})
				ref.Set(nil)
			}
			return nil
		},
	))
	if err != nil {
		return nil, err
	}
	return invalidRefs, nil
}

func RemoveInvalidRefsInAllCollections() (invalidCollectionRefs map[string][]InvalidRefData, err error) {
	config.Logger.Println("mongo.RemoveInvalidRefsInAllCollections()")
	for name, collection := range Collections {
		refs, err := collection.RemoveInvalidRefs()
		if err != nil {
			return nil, err
		}
		invalidCollectionRefs[name] = refs
	}
	return invalidCollectionRefs, nil
}

// // Returns an iterator of dereferenced refs, or an error iterator if there was an error
// func DereferenceIterator(refs ...Ref) model.Iterator {
// 	var docs []interface{}
// 	for i := range refs {
// 		doc, err := refs[i].Get()
// 		if err != nil {
// 			err = errs.Format("%s: %s", refs[i].ID.Hex(), err.Error())
// 			return model.NewErrorOnlyIterator(err)
// 		} else if doc != nil {
// 			docs = append(docs, doc)
// 		}
// 	}
// 	return model.NewSliceIterator(docs...)
// }

// // Returns an iterator for all refs without error and a slice of the errors
// func FailsafeDereferenceIterator(refs ...Ref) (i model.Iterator, errors []error) {
// 	var docs []interface{}
// 	for i := range refs {
// 		doc, err := refs[i].Get()
// 		if err != nil {
// 			err = errs.Format("%s: %s", refs[i].ID.Hex(), err.Error())
// 			errors = append(errors, err)
// 		} else if doc != nil {
// 			docs = append(docs, doc)
// 		}
// 	}
// 	return model.NewSliceIterator(docs...), errors
// }

func ReverseBsonD(d bson.D) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}
