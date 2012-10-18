package mongo

import "labix.org/v2/mgo"

///////////////////////////////////////////////////////////////////////////////
// query_limit

type query_limit struct {
	query_base
	limit int
}

func (self *query_limit) mongoQuery() (q *mgo.Query, err error) {
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Limit(self.limit), nil
}

func (self *query_limit) Selector() string {
	return ""
}
