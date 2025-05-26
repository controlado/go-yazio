package user

import (
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/pkg/domain/date"
)

func TestUser_SinceRegistAt(t *testing.T) {
	t.Parallel()

	var ( // static
		startTime = time.Date(2023, 8, 26, 12, 0, 0, 0, time.UTC)
		endTime   = time.Date(2025, 5, 13, 15, 0, 0, 0, time.UTC)
	)

	testBlocks := []struct {
		name string
		ud   *Data
		want date.Range
	}{
		{
			name: "valid call",
			ud:   &Data{Registration: startTime},
			want: date.Range{Start: startTime, End: endTime},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.ud.SinceRegistAt(endTime)
			assert.Equal(t, got, tb.want)
		})
	}
}
