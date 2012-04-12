package utils

import "time"

func ConvertTimeString(value, formatIn, formatOut string) (result string, err error) {
	if value == "" {
		return "", nil
	}
	t, err := time.Parse(formatIn, value)
	if err != nil {
		return "", err
	}
	return t.Format(formatOut), nil
}

func DayTimeRange(someTimeOfTheDay time.Time) (from, until time.Time) {
	from = DayBeginningTime(someTimeOfTheDay)
	until = from.Add(time.Hour * 24)
	return from, until
}

func TimeInRange(t, from, until time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Before(until) || t.Equal(until))
}

func DayBeginningTime(someTimeOfTheDay time.Time) time.Time {
	year, month, day := someTimeOfTheDay.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, someTimeOfTheDay.Location())
}

type SortableTimeSlice []time.Time

// Len is the number of elements in the collection.
func (self SortableTimeSlice) Len() int {
	return len(self)
}

// Less returns whether the element with index i should sort
// before the element with index j.
func (self SortableTimeSlice) Less(i, j int) bool {
	return self[i].Before(self[j])
}

// Swap swaps the elements with indexes i and j.
func (self SortableTimeSlice) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
