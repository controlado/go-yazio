package unit

import (
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

func TestBase_String(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name string
		b    Base
		want string
	}{
		{
			name: "kilocalorie",
			b:    Kilocalorie,
			want: "kcal",
		},
		{
			name: "milliliter",
			b:    Milliliter,
			want: "ml",
		},
		{
			name: "gram",
			b:    Gram,
			want: "g",
		},
		{
			name: "miligram",
			b:    Miligram,
			want: "mg",
		},
		{
			name: "microgram",
			b:    Microgram,
			want: "mcg",
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.b.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
