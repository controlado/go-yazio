package meal

import (
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

func TestMeal_String(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name string
		time Time
		want string
	}{
		{name: "breakfast", time: Breakfast, want: "breakfast"},
		{name: "lunch", time: Lunch, want: "lunch"},
		{name: "dinner", time: Dinner, want: "dinner"},
		{name: "snack", time: Snack, want: "snack"},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.time.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
