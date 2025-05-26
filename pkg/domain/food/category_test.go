package food

import (
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
)

func TestCategory_String(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name string
		c    Category
		want string
	}{
		{"meat", Meat, "meat"},
		{"misc", Miscellaneous, "miscellaneous"},
		{"choco", Chocolate, "chocolate"},
		{"non-alcoholic", NonAlcoholicDrink, "drinksnonalcoholic"},
		{"cheese", Cheese, "cheese"},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.c.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
