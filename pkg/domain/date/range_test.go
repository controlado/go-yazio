package date

import (
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

func TestRange_String(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name string
		dr   *Range
		want string
	}{
		{
			name: "valid start/end dates",
			dr: &Range{
				Start: time.Date(2023, 6, 1, 10, 0, 0, 0, time.UTC),
				End:   time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			},
			want: "1 June 2023 - 1 December 2023 (183 days)",
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.dr.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
