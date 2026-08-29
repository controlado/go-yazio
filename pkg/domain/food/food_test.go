package food

import (
	"fmt"
	"testing"

	"github.com/controlado/go-yazio/v4/internal/testutil/assert"
	"github.com/controlado/go-yazio/v4/pkg/domain/intake"
	"github.com/controlado/go-yazio/v4/pkg/domain/unit"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	t.Parallel()

	var ( // static
		staticUUID = uuid.New()
	)

	type args struct {
		name string
		cat  Category
		nut  Nutrients
		opts []Option
	}

	testBlocks := []struct {
		name    string
		args    args
		want    Food
		wantErr bool
	}{
		{name: "invalid name", wantErr: true},
		{
			name: "using default unit/serving",
			args: args{
				name: "Banana",
				cat:  Miscellaneous,
				nut:  Nutrients{intake.AddedSugar: 0.1},
				opts: []Option{
					WithID(staticUUID),
				},
			},
			want: Food{
				ID:                      staticUUID,
				Name:                    "Banana",
				BaseUnit:                unit.Gram,
				Category:                Miscellaneous,
				Nutrients:               Nutrients{intake.AddedSugar: 0.1},
				NutrientReferenceAmount: defaultNutrientReferenceAmount,
				Servings:                []Serving{defaultServing},
			},
		},
		{
			name: "using options: defining unit/serving/nutrient reference",
			args: args{
				name: "Liquid",
				cat:  Miscellaneous,
				nut:  Nutrients{intake.Water: 0.1},
				opts: []Option{
					WithID(staticUUID),
					WithBaseUnit(unit.Milliliter),
					WithNutrientsPer(18),
					WithNewServing(Pack, 50),
				},
			},
			want: Food{
				ID:                      staticUUID,
				Name:                    "Liquid",
				BaseUnit:                unit.Milliliter,
				Category:                Miscellaneous,
				Nutrients:               Nutrients{intake.Water: 0.1},
				NutrientReferenceAmount: 18,
				Servings: []Serving{
					{
						Kind:             Pack,
						AmountPerServing: 50,
					},
				},
			},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(
				tb.args.name,
				tb.args.cat,
				tb.args.nut,
				tb.args.opts...,
			)
			assert.WantErr(t, tb.wantErr, err)
			assert.DeepEqual(t, got, tb.want)
		})
	}
}

func TestFood_String(t *testing.T) {
	t.Parallel()

	var ( // static
		staticUUID = uuid.New()
	)

	testBlocks := []struct {
		name string
		f    Food
		want string
	}{
		{
			name: "valid food",
			f: Food{
				ID:       staticUUID,
				Name:     "Sadness",
				BaseUnit: unit.Microgram,
				Category: Meat,
				Nutrients: Nutrients{
					intake.Salt:    50,
					intake.Protein: 10,
					intake.Water:   20,
				},
			},
			want: fmt.Sprintf("Food(%q, %s, %s)", "Sadness", staticUUID, Meat),
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()
			got := tb.f.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
