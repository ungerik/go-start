package mongo

import "labix.org/v2/mgo"

///////////////////////////////////////////////////////////////////////////////
// limitQuery

type limitQuery struct {
	queryBase
	limit int
}

func (self *limitQuery) mongoQuery() (q *mgo.Query, err error) {
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Limit(self.limit), nil
}

func (self *limitQuery) Selector() string {
	return ""
}
