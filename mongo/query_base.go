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
// query_base

type query_base struct {
	thisQuery   Query
	parentQuery Query
}

func (self *query_base) init(thisQuery Query, parentQuery Query) {
	self.thisQuery = thisQuery
	self.parentQuery = parentQuery
}

func (self *query_base) subDocumentSelector() string {
	return ""
}

func (self *query_base) bsonSelector() bson.M {
	return self.thisQuery.bsonSelector()
}

func (self *query_base) ParentQuery() Query {
	return self.parentQuery
}

func (self *query_base) Collection() *Collection {
	for query := self.thisQuery; query != nil; query = query.ParentQuery() {
		if collection, ok := query.(*Collection); ok {
			return collection
		}
	}
	return nil
}

func (self *query_base) Count() (n int, err error) {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return 0, err
	}
	return q.Count()
}

func (self *query_base) SubDocument(selector string) Query {
	q := &query_subDocument{selector: selector}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) Skip(skip int) Query {
	if skip < 0 {
		return &QueryError{self.ParentQuery(), errs.Format("Invalid negative skip count: %d", skip)}
	}
	q := &query_skip{skip: skip}
	q.init(q, self.thisQuery)
	return q
}

func (self *query_base) Limit(limit int) Query {
	if limit < 0 {
		return &QueryError{self.ParentQuery(), errs.Format("Invalid negative limit: %d", limit)}
	}
	q := &query_limit{limit: limit}
	q.init(q, self.thisQuery)
	return q
}

func (self *query_base) Sort(selectors ...string) Query {
	for i := range selectors {
		selectors[i] = strings.ToLower(selectors[i])
	}
	q := &query_sort{selectors: selectors}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) SortFunc(compareFunc interface{}) model.Iterator {
	return model.NewSortIterator(self.Iterator(), compareFunc)
}

func (self *query_base) Explain() string {
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

func (self *query_base) IsFilter() bool {
	return false
}

func (self *query_base) Filter(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterEqual{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterWhere(javascript string) Query {
	return self.Filter("$where", javascript)
}

func (self *query_base) FilterFunc(passFilter model.FilterFunc) model.Iterator {
	return &model.FilterIterator{self.Iterator(), passFilter}
}

func (self *query_base) FilterRef(selector string, refs ...Ref) Query {
	selector = strings.ToLower(selector)
	q := &query_filterRef{selector: selector, refs: refs}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterID(ids ...bson.ObjectId) Query {
	if len(ids) == 1 {
		return self.Filter("_id", ids[0])
	}
	values := make([]interface{}, len(ids))
	for i := range ids {
		values[i] = ids[i]
	}
	return self.FilterIn("_id", values)
}

func (self *query_base) FilterEqualCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterEqualCaseInsensitive{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterNotEqual(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterNotEqual{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterLess(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterLess{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterGreater(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterGreater{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterLessEqual(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterLessEqual{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterGreaterEqual(selector string, value interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterGreaterEqual{selector: selector, value: value}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterModulo(selector string, divisor, result interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterModulo{selector: selector, divisor: divisor, result: result}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterIn(selector string, values ...interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterIn{selector: selector, values: values}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterNotIn(selector string, values ...interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterNotIn{selector: selector, values: values}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterAllIn(selector string, values ...interface{}) Query {
	selector = strings.ToLower(selector)
	q := &query_filterAllIn{selector: selector, values: values}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterArraySize(selector string, size int) Query {
	selector = strings.ToLower(selector)
	q := &query_filterArraySize{selector: selector, size: size}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterStartsWith(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterStartsWith{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterStartsWithCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterStartsWith{selector: selector, str: str, caseInsensitive: true}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterEndsWith(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterEndsWith{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterEndsWithCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterEndsWith{selector: selector, str: str, caseInsensitive: true}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterContains(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterContains{selector: selector, str: str}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterContainsCaseInsensitive(selector string, str string) Query {
	selector = strings.ToLower(selector)
	q := &query_filterContains{selector: selector, str: str, caseInsensitive: true}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) FilterExists(selector string, exists bool) Query {
	selector = strings.ToLower(selector)
	q := &query_filterExists{selector: selector, exists: exists}
	q.init(q, self.thisQuery)
	return checkQuery(q)
}

func (self *query_base) Or() Query {
	if !self.thisQuery.IsFilter() {
		return &QueryError{self.thisQuery, errs.Format("Or() can only be called after Filter()")}
	}
	q := &query_or{}
	q.init(q, self.thisQuery)
	return q
}

func (self *query_base) OneDocument(resultRef interface{}) error {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return err
	}
	err = q.One(resultRef)
	if err != nil {
		return err
	}
	// resultRef has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.Collection().InitDocument(documentFromResultPtr(resultRef))
	return nil
}

func (self *query_base) TryOneDocument(resultRef interface{}) (found bool, err error) {
	err = self.OneDocument(resultRef)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (self *query_base) OneSubDocument(selector string, resultRef interface{}) error {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return err
	}
	q = q.Select(bson.M{selector: 1})
	err = q.One(resultRef)
	if err != nil {
		return err
	}
	// resultRef has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.Collection().InitSubDocument(resultRef)
	return nil
}

func (self *query_base) TryOneSubDocument(selector string, resultRef interface{}) (found bool, err error) {
	err = self.OneSubDocument(selector, resultRef)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (self *query_base) Iterator() model.Iterator {
	return NewMongoIterator(self.thisQuery)
}

func (self *query_base) OneDocumentID() (id bson.ObjectId, err error) {
	q, err := self.thisQuery.mongoQuery()
	if err != nil {
		return bson.ObjectId(""), err
	}
	var doc DocumentBase
	err = q.One(&doc)
	return doc.ID, err
}

func (self *query_base) TryOneDocumentID() (id bson.ObjectId, found bool, err error) {
	id, err = self.OneDocumentID()
	if err == mgo.ErrNotFound {
		return id, false, nil
	}
	return id, err == nil, err
}

func (self *query_base) DocumentIDs() (ids []bson.ObjectId, err error) {
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

func (self *query_base) Refs() (refs []Ref, err error) {
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

func (self *query_base) UpdateOne(selector string, value interface{}) error {
	bsonQuery, err := bsonQuery(self.thisQuery)
	if err != nil {
		return err
	}
	return self.Collection().collection.Update(bsonQuery, bson.M{"$set": bson.M{selector: value}})
}

func (self *query_base) UpdateAll(selector string, value interface{}) (numUpdated int, err error) {
	bsonQuery, err := bsonQuery(self.thisQuery)
	if err != nil {
		return 0, err
	}
	info, err := self.Collection().collection.UpdateAll(bsonQuery, bson.M{"$set": bson.M{selector: value}})
	return info.Updated, err
}

func (self *query_base) UpdateSubDocument(selector string, subDocument interface{}) error {
	bsonQuery, err := bsonQuery(self.thisQuery)
	if err != nil {
		return err
	}
	if selector == "" {
		valMap := make(bson.M)
		model.VisitMaxDepth(subDocument, 1, model.FieldOnlyVisitor(
			func(field *model.MetaData) error {
				valMap[field.Name] = field.Value.Addr().Interface()
				return nil
			},
		))
		if len(valMap) == 0 {
			return nil
		}
		return self.Collection().collection.Update(bsonQuery, bson.M{"$set": valMap})
	}
	return self.Collection().collection.Update(bsonQuery, bson.M{"$set": bson.M{selector: subDocument}})
}

func (self *query_base) RemoveAll() (numRemoved int, err error) {
	info, err := self.Collection().collection.RemoveAll(self.bsonSelector())
	return info.Removed, err
}
