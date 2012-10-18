package mongo

import (
	"github.com/ungerik/go-start/errs"
	"labix.org/v2/mgo"
)

///////////////////////////////////////////////////////////////////////////////
// query_or

type query_or struct {
	query_base
}

func (self *query_or) mongoQuery() (q *mgo.Query, err error) {
	return nil, errs.Format("Can't create a mongo query. query_or needs a query_filterEqual chained after it")
}

func (self *query_or) Selector() string {
	return ""
}
