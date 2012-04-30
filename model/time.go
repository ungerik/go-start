package model

import "time"

const unixDateSec int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * (24 * 60 * 60)

// Time in milliseconds since January 1, year 1 00:00:00 UTC.
// Time values are always in UTC.
type Time int64

func (self *Time) Get() time.Time {
	unixsec := int64(*self)/1000 - unixDateSec
	unixmsec := int64(*self) % 1000
	return time.Unix(unixsec, unixmsec*1e6)
}

func (self *Time) Set(t time.Time) {
	unixsec := t.Unix()
	unixmsec := int64(t.Nanosecond()) / 1e6
	*self = Time((unixsec+unixDateSec)*1000 + unixmsec)
}

func (self *Time) IsEmpty() bool {
	return *self == 0
}
