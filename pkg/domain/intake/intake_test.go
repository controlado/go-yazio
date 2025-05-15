package intake

import (
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

func TestSingleRange_Average(t *testing.T) {
	t.Parallel()

	var (
		defaultDate = time.Now()
		testBlocks  = []struct {
			name string
			sr   SingleRange
			want SingleAverage
		}{
			{
				name: "avarage should be 4.25",
				sr: SingleRange{
					{defaultDate, 2},
					{defaultDate, 3},
					{defaultDate, 5},
					{defaultDate, 7},
				},
				want: SingleAverage{
					DaysLength: 4,
					Average:    4.25,
				},
			},
			{
				name: "empty range should return zero average",
				sr:   SingleRange{},
				want: SingleAverage{},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.sr.Average()
			assert.Equal(t, tb.want, got)
		})
	}
}

func TestMacrosRange_Average(t *testing.T) {
	t.Parallel()

	var (
		defaultDate = time.Now()
		testBlocks  = []struct {
			name string
			mr   MacrosRange
			want MacrosAverage
		}{
			{
				name: "average of 3 days macros",
				mr: MacrosRange{
					{defaultDate, 2000, 200, 80, 120},
					{defaultDate, 2200, 220, 90, 130},
					{defaultDate, 2100, 210, 85, 125},
				},
				want: MacrosAverage{
					DaysLength: 3,
					Energy:     2100,
					Carb:       210,
					Fat:        85,
					Protein:    125,
				},
			},
			{
				name: "empty range should return zero average",
				mr:   MacrosRange{},
				want: MacrosAverage{},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.mr.Average()
			assert.Equal(t, got, tb.want)
		})
	}
}
