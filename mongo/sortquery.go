package mongo

import (
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// sortQuery

type sortQuery struct {
	queryBase
	selector  string
	direction int
}

func (self *sortQuery) mongoQuery() (q *mgo.Query, err error) {
	elems := bson.D{{self.selector, self.direction}}
	for query := self.parentQuery; query != nil; query = query.ParentQuery() {
		s, ok := query.(*sortQuery)
		if !ok {
			break
		}
		elems = append(elems, bson.DocElem{s.selector, s.direction})
	}
	ReverseBsonD(elems)
	q, err = self.parentQuery.mongoQuery()
	if err != nil {
		return nil, err
	}
	return q.Sort(elems), nil
}

func (self *sortQuery) Selector() string {
	return self.selector
}

func ReverseBsonD(d bson.D) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}
