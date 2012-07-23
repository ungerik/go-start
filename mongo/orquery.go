package mongo

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/mgo"
)

///////////////////////////////////////////////////////////////////////////////
// orQuery

type orQuery struct {
	queryBase
}

func (self *orQuery) mongoQuery() (q *mgo.Query, err error) {
	return nil, errs.Format("Can't create a mongo query. orQuery needs a filterQuery chained after it")
}

func (self *orQuery) Selector() string {
	return ""
}
