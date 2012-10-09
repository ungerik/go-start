package mongo

import (
	"labix.org/v2/mgo"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// sortQuery

type sortQuery struct {
	queryBase
	selectors []string
}

func (self *sortQuery) mongoQuery() (q *mgo.Query, err error) {
	selectors := self.selectors
	for query := self.parentQuery; query != nil; query = query.ParentQuery() {
		s, ok := query.(*sortQuery)
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

func (self *sortQuery) Selector() string {
	return strings.Join(self.selectors, ",")
}
