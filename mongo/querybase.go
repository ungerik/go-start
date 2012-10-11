package mongo

import (
	"fmt"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
)

//	http://www.mongodb.org/display/DOCS/Advanced+Queries
//	http://www.10gen.com/static/downloads/mongodb_qrc_queries.pdf

func init() {
	debug.Nop()
}

///////////////////////////////////////////////////////////////////////////////
// queryBase

type queryBase struct {
	thisQuery   Query
	parentQuery Query
}

func (self *queryBase) init(thisQuery Query, parentQuery Query) {
	self.thisQuery = thisQuery
	self.parentQuery = parentQuery
}

func (self *queryBase) subDocumentSelector() string {
	return ""
}

func (self *queryBase) bsonSelector() bson.M {
	return self.thisQuery.bsonSelector()
}

func (self *queryBase) ParentQuery() Query {
	return self.parentQuery
}

func (self *queryBase) Collection() *Collection {
	for query := self.thisQuery; query != nil; query = query.ParentQuery() {
		if collection, ok := query.(*Collection); ok {
			return collection
		}
	}
	return nil
}

func (self *queryBase) Count() (n int, err error) {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return 0, err
	}
	return q.Count()
}

func (self *queryBase) SubDocument(selector string) Query {
	q := &subDocumentQuery{selector: selector}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) Skip(skip int) Query {
	if skip < 0 {
		return &QueryError{self.ParentQuery(), errs.Format("Invalid negative skip count: %d", skip)}
	}
	q := &skipQuery{skip: skip}
	q.init(q, self.thisQuery)
	return q
}

func (self *queryBase) Limit(limit int) Query {
	if limit < 0 {
		return &QueryError{self.ParentQuery(), errs.Format("Invalid negative limit: %d", limit)}
	}
	q := &limitQuery{limit: limit}
	q.init(q, self.thisQuery)
	return q
}

func (self *queryBase) Sort(selectors ...string) Query {
	for i := range selectors {
		selectors[i] = strings.ToLower(selectors[i])
	}
	q := &sortQuery{selectors: selectors}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) SortFunc(lessFunc func(a, b interface{}) bool) model.Iterator {
	return model.SortIterator(self.Iterator(), lessFunc)
}

func (self *queryBase) Explain() string {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return err.Error()
	}
	m := bson.M{}
	err = q.Explain(m)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Explain: %#v\n", m)
}

func (self *queryBase) IsFilter() bool {
	return false
}

func (self *queryBase) Filter(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterQuery{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterWhere(javascript string) Query {
	return self.Filter("$where", javascript)
}

func (self *queryBase) FilterFunc(passFilter model.FilterFunc) model.Iterator {
	return &model.FilterIterator{self.Iterator(), passFilter}
}

func (self *queryBase) FilterRef(selector string, ref ...Ref) Query {
	selector = strings.ToLower(selector)
	q := &filterRefQuery{selector: selector, refs: ref}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterEqualCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterEqualCaseInsensitiveQuery{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterNotEqual(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterNotEqualQuery{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterLess(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterLessQuery{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterGreater(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterGreaterQuery{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterLessEqual(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterLessEqualQuery{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterGreaterEqual(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterGreaterEqualQuery{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterModulo(selector string, divisor, result interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterModuloQuery{selector: selector, divisor: divisor, result: result}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterIn(selector string, values ...interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterInQuery{selector: selector, values: values}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterNotIn(selector string, values ...interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterNotInQuery{selector: selector, values: values}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterAllIn(selector string, values ...interface{}) Query {
	selector = strings.ToLower(selector)
	q := &filterAllInQuery{selector: selector, values: values}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterArraySize(selector string, size int) Query {
	selector = strings.ToLower(selector)
	q := &filterArraySizeQuery{selector: selector, size: size}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterStartsWith(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterStartsWithQuery{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterStartsWithCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterStartsWithQuery{selector: selector, str: str, caseInsensitive: true}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterEndsWith(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterEndsWithQuery{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterEndsWithCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterEndsWithQuery{selector: selector, str: str, caseInsensitive: true}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterContains(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterContainsQuery{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterContainsCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &filterContainsQuery{selector: selector, str: str, caseInsensitive: true}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) FilterExists(selector string, exists bool) Query {
	selector = strings.ToLower(selector)
	q := &filterExistsQuery{selector: selector, exists: exists}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *queryBase) Or() Query {
	if !self.thisQuery.IsFilter() {
		return &QueryError{self.thisQuery, errs.Format("Or() can only be called after Filter()")}
	}
	q := &orQuery{}
	q.init(q, self.thisQuery)
	return q
}

func (self *queryBase) OneDocument(resultPtr interface{}) error {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return err
	}
	err = q.One(resultPtr)
	if err != nil {
		return err
	}
	// resultPtr has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.Collection().InitDocument(documentFromResultPtr(resultPtr))
	return nil
}

func (self *queryBase) TryOneDocument(resultPtr interface{}) (found bool, err error) {
	err = self.OneDocument(resultPtr)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (self *queryBase) Iterator() model.Iterator {
	return newIterator(self.thisQuery)
}

func (self *queryBase) OneDocumentID() (id bson.ObjectId, err error) {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return bson.ObjectId(""), err
	}
	var doc DocumentBase
	err = q.One(&doc)
	return doc.ID, err
}

func (self *queryBase) TryOneDocumentID() (id bson.ObjectId, found bool, err error) {
	id, err = self.OneDocumentID()
	if err == mgo.ErrNotFound {
		return id, false, nil
	}
	return id, err == nil, err
}

func (self *queryBase) DocumentIDs() (ids []bson.ObjectId, err error) {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	i := q.Select(bson.M{"_id": 1}).Iter()
	var doc DocumentBase
	for i.Next(&doc) {
		ids = append(ids, doc.ID)
	}
	if i.Err() != nil {
		return nil, i.Err()
	}
	return ids, nil
}

func (self *queryBase) Refs() (refs []Ref, err error) {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	i := q.Select(bson.M{"_id": 1}).Iter()
	collection := self.Collection()
	var doc DocumentBase
	for i.Next(&doc) {
		refs = append(refs, collection.RefForID(doc.ID))
	}
	if i.Err() != nil {
		return nil, i.Err()
	}
	return refs, nil
}

func (self *queryBase) UpdateOne(selector string, value interface{}) error {
	bsonQuery, err := bsonQuery(self.thisQuery)
	if err != nil {
		return err
	}
	return self.Collection().collection.Update(bsonQuery, bson.M{"$set": bson.M{selector: value}})
}

func (self *queryBase) UpdateAll(selector string, value interface{}) (numUpdated int, err error) {
	bsonQuery, err := bsonQuery(self.thisQuery)
	if err != nil {
		return 0, err
	}
	info, err := self.Collection().collection.UpdateAll(bsonQuery, bson.M{"$set": bson.M{selector: value}})
	return info.Updated, err
}

func (self *queryBase) RemoveAll() (numRemoved int, err error) {
	info, err := self.Collection().collection.RemoveAll(self.bsonSelector())
	return info.Removed, err
}
