package times

import "time"

func Future() time.Time {
	return time.Now().Add(time.Hour)
}

func Past() time.Time {
	return time.Now().Add(-time.Hour)
}

func extractDate(t time.Time, loc *time.Location) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		loc,
	)
}

func FutureDate(loc *time.Location) time.Time {
	nextYear := time.Now().AddDate(1, 0, 0)
	return extractDate(nextYear, loc)
}

func PastDate(loc *time.Location) time.Time {
	pastYear := time.Now().AddDate(-1, 0, 0)
	return extractDate(pastYear, loc)
}
