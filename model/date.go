package model

import "time"

const DateFormat = "2006-01-02"

func NewDate(value string) *Date {
	return (*Date)(&value)
}

type Date string

func (self *Date) Get() string {
	return string(*self)
}

func (self *Date) Set(value string) error {
	*self = Date(value) // set value in any case, so that user can see wrong value in form

	if value != "" {
		if _, err := time.Parse(DateFormat, value); err != nil {
			return err
		}
	}
	*self = Date(value)
	return nil
}

func (self *Date) SetTodayUTC() {
	self.SetTime(time.Now().UTC())
}

func (self *Date) Time() time.Time {
	time, err := time.Parse(DateFormat, self.Get())
	if err != nil {
		panic(err)
	}
	return time
}

func (self *Date) SetTime(t time.Time) {
	*self = Date(t.Format(DateFormat))
}

func (self *Date) UnixNanoseconds() int64 {
	return self.Time().UnixNano()
}

func (self *Date) SetUnixNanoseconds(nanos int64) {
	self.SetTime(time.Unix(0, nanos).UTC())
}

func (self *Date) Format(format string) string {
	if *self == "" {
		return ""
	}
	return self.Time().Format(format)
}

func (self *Date) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Date) SetEmpty() {
	*self = ""
}

func (self *Date) String() string {
	return self.Get()
}

func (self *Date) SetString(str string) error {
	return self.Set(str)
}

func (self *Date) FixValue(metaData *MetaData) {
}

// todo min max
func (self *Date) Validate(metaData *MetaData) error {
	value := self.Get()
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	if self.Required(metaData) || value != "" {
		if _, err := time.Parse(DateFormat, value); err != nil {
			return err
		}
	}
	return nil
}

func (self *Date) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}
