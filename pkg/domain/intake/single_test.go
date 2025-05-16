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
					{Water, defaultDate, 2},
					{Water, defaultDate, 3},
					{Water, defaultDate, 5},
					{Water, defaultDate, 7},
				},
				want: SingleAverage{Water, 4, 4.25},
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

func TestSingleAverage_String(t *testing.T) {
	t.Parallel()

	var (
		testBlocks = []struct {
			name string
			sa   SingleAverage
			want string
		}{
			{
				name: "correct call (vitamin)",
				sa:   SingleAverage{VitaminA, 720, 5},
				want: "720 days: 5.0mcg",
			},
			{
				name: "correct call (water)",
				sa:   SingleAverage{Water, 320, 2223},
				want: "320 days: 2223.0ml",
			},
			{
				name: "zero-value",
				sa:   SingleAverage{},
				want: "Empty intakes data to calculate the average",
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.sa.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
