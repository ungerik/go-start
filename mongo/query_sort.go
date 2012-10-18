package mongo

import (
	"labix.org/v2/mgo"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// query_sort

type query_sort struct {
	query_base
	selectors []string
}

func (self *query_sort) mongoQuery() (q *mgo.Query, err error) {
	selectors := self.selectors
	for query := self.parentQuery; query != nil; query = query.ParentQuery() {
		s, ok := query.(*query_sort)
		if !ok {
			break
		}
		selectors = append(s.selectors, selectors...)
	}
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Sort(selectors...), nil
}

func (self *query_sort) Selector() string {
	return strings.Join(self.selectors, ",")
}
