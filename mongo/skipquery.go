package mongo

import "labix.org/v2/mgo"

///////////////////////////////////////////////////////////////////////////////
// skipQuery

type skipQuery struct {
	queryBase
	skip int
}

func (self *skipQuery) mongoQuery() (q *mgo.Query, err error) {
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Skip(self.skip), nil
}

func (self *skipQuery) Selector() string {
	return ""
}
