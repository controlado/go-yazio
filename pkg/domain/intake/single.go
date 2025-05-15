package intake

import (
	"fmt"
	"time"
)

type Single struct {
	Kind  Kind
	Date  time.Time
	Value float64
}

type SingleRange []Single

func (sr SingleRange) Average() SingleAverage {
	var (
		kindSample  Kind
		totalValues float64
		rangeLength = len(sr)
	)

	if rangeLength == 0 {
		return SingleAverage{}
	}

	for i, intake := range sr {
		if i == 0 {
			kindSample = intake.Kind
		}
		totalValues += intake.Value
	}

	return SingleAverage{
		Kind:       kindSample,
		DaysLength: rangeLength,
		Average:    totalValues / float64(rangeLength),
	}
}

type SingleAverage struct {
	Kind       Kind
	DaysLength int
	Average    float64
}

func (sa SingleAverage) String() string {
	if sa.DaysLength < 1 {
		return "Empty intakes data to calculate the average"
	}

	return fmt.Sprintf("%d days: %.1f%s",
		sa.DaysLength,
		sa.Average,
		sa.Kind.baseUnit,
	)
}
