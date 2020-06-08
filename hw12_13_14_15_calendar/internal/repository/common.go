package repository

import "time"

type Period int

const (
	Day Period = iota
	Week
	Month
)

func GetTimeRange(dt time.Time, p Period) (from time.Time, to time.Time) {
	from = dt.Round(24 * time.Hour)

	switch p {
	case Day:
		to = dt.AddDate(0, 0, 1)
	case Week:
		to = dt.AddDate(0, 0, 7)
	case Month:
		to = dt.AddDate(0, 1, 0)
	}

	return
}
