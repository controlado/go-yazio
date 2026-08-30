package intake

import (
	"testing"

	"github.com/controlado/go-yazio/v5/internal/testutil/assert"
	"github.com/controlado/go-yazio/v5/pkg/domain/unit"
)

func TestKind_BaseUnit(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name string
		kind Kind
		want unit.Base
	}{
		{name: "cholesterol", kind: Cholesterol, want: unit.Milligram},
		{name: "sodium", kind: Sodium, want: unit.Milligram},
		{name: "water", kind: Water, want: unit.Gram},
		{name: "alcohol", kind: Alcohol, want: unit.Gram},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tb.want, tb.kind.BaseUnit())
			assert.Equal(t, tb.want.String(), tb.kind.Unit())
		})
	}
}
