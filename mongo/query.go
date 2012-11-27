package mongo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/model"
)

///////////////////////////////////////////////////////////////////////////////
// Query

/*
Query is the interface with all query methods for mongo.
*/
type Query interface {
	subDocumentSelector() string
	bsonSelector() bson.M
	mongoQuery() (q *mgo.Query, err error)

	Selector() string

	ParentQuery() Query
	Collection() *Collection

	SubDocument(selector string) Query
	Skip(int) Query
	Limit(int) Query
	Sort(selectors ...string) Query // sort reverse by prefixing the selector with a minus -

	// SortFunc will sort the query according to compareFunc.
	// compareFunc must have two arguments that are assignable from
	// the iterated document type or pointers to such a type.
	// The result of compareFunc must be a bool indicating
	// if the first argument is less than the second.
	SortFunc(compareFunc interface{}) model.Iterator

	// FilterX must be the first query on a Collection
	IsFilter() bool
	Filter(selector string, value interface{}) Query
	FilterWhere(javascript string) Query

	// Filter via a Go function. Note that all documents have to be loaded
	// in memory in order for Go code to be able to filter it.
	FilterFunc(passFilter model.FilterFunc) model.Iterator
	FilterRef(selector string, refs ...Ref) Query
	FilterID(ids ...bson.ObjectId) Query
	FilterEqualCaseInsensitive(selector string, str string) Query
	FilterNotEqual(selector string, value interface{}) Query
	FilterLess(selector string, value interface{}) Query
	FilterGreater(selector string, value interface{}) Query
	FilterLessEqual(selector string, value interface{}) Query
	FilterGreaterEqual(selector string, value interface{}) Query
	FilterModulo(selector string, divisor, result interface{}) Query
	FilterIn(selector string, values ...interface{}) Query
	FilterNotIn(selector string, values ...interface{}) Query
	FilterAllIn(selector string, values ...interface{}) Query
	FilterArraySize(selector string, size int) Query
	FilterStartsWith(selector string, str string) Query
	FilterStartsWithCaseInsensitive(selector string, str string) Query
	FilterEndsWith(selector string, str string) Query
	FilterEndsWithCaseInsensitive(selector string, str string) Query
	FilterContains(selector string, str string) Query
	FilterContainsCaseInsensitive(selector string, str string) Query
	FilterExists(selector string, exists bool) Query

	Or() Query

	// Statistics
	Count() (n int, err error)
	// Distinct() int
	Explain() string

	// Read
	OneDocument(resultRef interface{}) error
	TryOneDocument(resultRef interface{}) (found bool, err error)

	OneSubDocument(selector string, resultRef interface{}) error
	TryOneSubDocument(selector string, resultRef interface{}) (found bool, err error)

	Iterator() model.Iterator

	OneDocumentID() (id bson.ObjectId, err error)
	TryOneDocumentID() (id bson.ObjectId, found bool, err error)
	DocumentIDs() (ids []bson.ObjectId, err error)
	Refs() (refs []Ref, err error)

	// Write
	UpdateOne(selector string, value interface{}) error
	UpdateAll(selector string, value interface{}) (numUpdated int, err error)
	// UpdateSubDocument updates a sub-part of queried documents.
	// selector can be an empty string if subDocument is on
	// the root level of the queried documents.
	UpdateSubDocument(selector string, subDocument interface{}) error

	// RemoveAll ignores Skip() and Limit()
	RemoveAll() (numRemoved int, err error)
}
