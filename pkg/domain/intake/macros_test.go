package intake

import (
	"strings"
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

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

func TestMacrosAverage_String(t *testing.T) {
	t.Parallel()

	var (
		charBreak  = "\n"
		testBlocks = []struct {
			name      string
			ma        MacrosAverage
			wantLines []string
		}{
			{
				name: "correct call",
				ma: MacrosAverage{
					DaysLength: 1278,
					Energy:     1700,
					Carb:       150,
					Fat:        40,
					Protein:    180,
				},
				wantLines: []string{
					`Average (1278 days)`,
					`Energy: 1700.00 (kcal)`,
					`Carb: 150.00 (g)`,
					`Fat: 40.00 (g)`,
					`Protein: 180.00 (g)`,
				},
			},
			{
				name: "zero-value",
				ma:   MacrosAverage{},
				wantLines: []string{
					`Average (0 days)`,
					`Energy: 0.00 (kcal)`,
					`Carb: 0.00 (g)`,
					`Fat: 0.00 (g)`,
					`Protein: 0.00 (g)`,
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.ma.String()
			assert.Equal(
				t,
				got,
				strings.Join(tb.wantLines, charBreak),
			)
		})
	}
}
