package user

import (
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/pkg/domain/date"
)

func TestUser_SinceRegist(t *testing.T) {
	t.Parallel()

	var (
		originalNowFn     = now
		defaultRegistTime = time.Date(2023, 8, 26, 12, 0, 0, 0, time.UTC)
		nowResponse       = time.Date(2025, 5, 13, 15, 0, 0, 0, time.UTC)
	)

	defer func() { now = originalNowFn }()
	now = func() time.Time { return nowResponse }

	testBlocks := []struct {
		name string
		u    *User
		want date.Range
	}{
		{
			name: "valid call",
			u:    &User{Registration: defaultRegistTime},
			want: date.Range{Start: defaultRegistTime, End: nowResponse},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			// t.Parallel()
			got := tb.u.SinceRegist()
			assert.Equal(t, got, tb.want)
		})
	}
}
