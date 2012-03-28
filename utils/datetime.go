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
	return time.Date(
		someTimeOfTheDay.Year(),
		someTimeOfTheDay.Month(),
		someTimeOfTheDay.Day(),
		0,
		0,
		0,
		0,
		someTimeOfTheDay.Location(),
	)
}
