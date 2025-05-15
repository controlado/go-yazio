package times

import "time"

func Future() time.Time {
	return time.Now().Add(time.Hour)
}

func Past() time.Time {
	return time.Now().Add(-time.Hour)
}
