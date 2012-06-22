package mongo

import (
	"launchpad.net/mgo"
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// sortQuery

type sortQuery struct {
	queryBase
	selector string
}

func (self *sortQuery) mongoQuery() (q *mgo.Query, err error) {
	elems := []string{self.selector}
	for query := self.parentQuery; query != nil; query = query.ParentQuery() {
		s, ok := query.(*sortQuery)
		if !ok {
			break
		}
		elems = append(elems, s.selector)
	}
	utils.ReverseStringSlice(elems)
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Sort(elems...), nil
}

func (self *sortQuery) Selector() string {
	return self.selector
}

// func ReverseBsonD(d bson.D) {
// 	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
// 		d[i], d[j] = d[j], d[i]
// 	}
// }
