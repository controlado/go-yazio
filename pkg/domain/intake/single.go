package intake

import (
	"fmt"
	"time"
)

type Single struct {
	Date  time.Time
	Value float64
}

type SingleRange []Single

func (sr SingleRange) Average() SingleAverage {
	var (
		totalValues float64
		rangeLength = len(sr)
	)

	if rangeLength == 0 {
		return SingleAverage{}
	}

	for _, intake := range sr {
		totalValues += intake.Value
	}

	return SingleAverage{
		DaysLength: rangeLength,
		Average:    totalValues / float64(rangeLength),
	}
}

type SingleAverage struct {
	DaysLength int
	Average    float64
}

func (sa SingleAverage) String() string {
	return fmt.Sprintf("%d days: %.2f", sa.DaysLength, sa.Average)
}
