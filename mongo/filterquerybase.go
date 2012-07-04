package mongo

import (
	"labix.org/v2/mgo/bson"
	"github.com/ungerik/go-start/errs"
	"labix.org/v2/mgo"

//	"github.com/ungerik/go-start/debug"
)

///////////////////////////////////////////////////////////////////////////////
// filterQueryBase

type filterQueryBase struct {
	queryBase
}

func (self *filterQueryBase) mongoQuery() (q *mgo.Query, err error) {
	chainedFilters := []Query{self.thisQuery}

	orChained := false
	for parent := self.parentQuery; parent != nil; parent = parent.ParentQuery() {
		if parent.IsFilter() {
			chainedFilters = append(chainedFilters, parent)
		}
		if _, or := parent.(*orQuery); or {
			// todo check if filter and or interleave
			orChained = true
		}
	}

	var bsonQuery bson.M

	if len(chainedFilters) == 1 {
		bsonQuery = self.thisQuery.bsonSelector()
	} else if orChained {
		bsonQueries := make([]bson.M, len(chainedFilters))
		for i, filter := range chainedFilters {
			bsonQueries[i] = filter.bsonSelector()
		}
		bsonQuery = bson.M{"$or": bsonQueries}
	} else {
		bsonQuery = bson.M{}
		for _, filter := range chainedFilters {
			for key, value := range filter.bsonSelector() {
				if existingValue, hasKey := bsonQuery[key]; hasKey {
					return nil, errs.Format("Can't filter %s for %v and %v", key, existingValue, value)
				}
				bsonQuery[key] = value
			}
		}
	}

	collection := self.Collection()
	collection.checkDBConnection()
	return collection.collection.Find(bsonQuery), nil
}

func (self *filterQueryBase) IsFilter() bool {
	return true
}
