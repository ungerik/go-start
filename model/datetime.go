package model

import "time"

const DateTimeFormat = "2006-01-02 15:04:05"
const ShortDateTimeFormat = "2006-01-02 15:04"

func NewDateTime(value string) *DateTime {
	return (*DateTime)(&value)
}

// DateTime 
type DateTime string

func (self *DateTime) Get() string {
	return string(*self)
}

func (self *DateTime) Set(value string) (err error) {
	*self = DateTime(value) // set value in any case, so that user can see wrong value in form

	if value != "" {
		if _, err = time.Parse(DateTimeFormat, value); err != nil {
			if _, err = time.Parse(ShortDateTimeFormat, value); err != nil {
				return err
			} else {
				value += ":00"
			}
		}
	}
	*self = DateTime(value)
	return nil
}

func (self *DateTime) SetNowUTC() {
	self.SetTime(time.Now().UTC())
}

func (self *DateTime) Time() time.Time {
	time, err := time.Parse(DateTimeFormat, self.Get())
	if err != nil {
		panic(err)
	}
	return time
}

func (self *DateTime) SetTime(t time.Time) {
	*self = DateTime(t.Format(DateTimeFormat))
}

func (self *DateTime) UnixNanoseconds() int64 {
	return self.Time().UnixNano()
}

func (self *DateTime) SetUnixNanoseconds(nanos int64) {
	self.SetTime(time.Unix(0, nanos).UTC())
}

func (self *DateTime) Format(format string) string {
	if *self == "" {
		return ""
	}
	return self.Time().Format(format)
}

func (self *DateTime) IsEmpty() bool {
	return len(*self) == 0
}

func (self *DateTime) SetEmpty() {
	*self = ""
}

func (self *DateTime) String() string {
	return self.Get()
}

func (self *DateTime) SetString(str string) error {
	return self.Set(str)
}

func (self *DateTime) FixValue(metaData *MetaData) {
}

// todo min max
func (self *DateTime) Validate(metaData *MetaData) error {
	value := self.Get()
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	if self.Required(metaData) || value != "" {
		if _, err := time.Parse(DateTimeFormat, value); err != nil {
			return err
		}
	}
	return nil
}

func (self *DateTime) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}
