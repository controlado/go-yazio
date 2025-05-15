package date

import (
	"fmt"
	"time"
)

const (
	humanLayout = "2 January 2006"
)

type Range struct {
	Start time.Time
	End   time.Time
}

func (r *Range) String() string {
	var (
		duration = r.End.Sub(r.Start) / (24 * time.Hour)
	)
	return fmt.Sprintf("%s - %s (%d days)",
		r.Start.Format(humanLayout),
		r.End.Format(humanLayout),
		int(duration),
	)
}
