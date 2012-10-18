package mongo

import "labix.org/v2/mgo"

///////////////////////////////////////////////////////////////////////////////
// query_skip

type query_skip struct {
	query_base
	skip int
}

func (self *query_skip) mongoQuery() (q *mgo.Query, err error) {
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Skip(self.skip), nil
}

func (self *query_skip) Selector() string {
	return ""
}
