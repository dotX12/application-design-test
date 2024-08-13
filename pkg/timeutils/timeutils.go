package timeutils

import (
	"errors"
	"time"
)

var ErrFromAfterTo = errors.New("from date is after to date")

func TimeToDate(timestamp time.Time) time.Time {
	return time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		0,
		0,
		0,
		0,
		timestamp.Location(),
	)
}

func DateToTime(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func DaysBetween(from time.Time, to time.Time) ([]time.Time, error) {
	from = TimeToDate(from)
	to = TimeToDate(to)

	if from.After(to) {
		return nil, ErrFromAfterTo
	}

	days := make([]time.Time, 0, int(to.Sub(from).Hours()/24)+1)
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days, nil
}
