package serving

import (
	"testing"

	"github.com/controlado/go-yazio/v4/internal/testutil/assert"
	"github.com/controlado/go-yazio/v4/pkg/domain/food"
)

func TestBaseUnitServingConstructors(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name                 string
		newServing           func() food.Serving
		wantKind             string
		wantAmountPerServing float64
	}{
		{name: "one gram serving", newServing: Gram, wantKind: "gram", wantAmountPerServing: 1},
		{name: "one milliliter serving", newServing: Milliliter, wantKind: "milliliter", wantAmountPerServing: 1},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			s := tb.newServing()
			assert.Equal(t, string(s.Kind), tb.wantKind)
			assert.Equal(t, s.AmountPerServing, tb.wantAmountPerServing)
		})
	}
}
