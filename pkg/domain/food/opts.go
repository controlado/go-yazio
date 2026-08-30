package food

import (
	"github.com/controlado/go-yazio/v5/pkg/domain/unit"
	"github.com/google/uuid"
)

type Option func(f *Food)

func WithID(i uuid.UUID) Option {
	return func(f *Food) {
		f.ID = i
	}
}

// WithBaseUnit sets b as the food's base unit.
//
// The base unit defines how nutrient reference amounts
// and serving sizes are measured.
func WithBaseUnit(b unit.Base) Option {
	return func(f *Food) {
		f.BaseUnit = b
	}
}

// WithNutrientsPer sets amount as the quantity of the
// food's base unit represented by its nutrients.
func WithNutrientsPer(amount float64) Option {
	return func(f *Food) {
		f.NutrientReferenceAmount = amount
	}
}

// WithNewServing adds a saved serving option of kind k
// to a food.
//
// amountPerServing is the quantity of the food's base
// unit represented by one serving.
//
// This option does not create a diary entry.
func WithNewServing(k ServingKind, amountPerServing float64) Option {
	return func(f *Food) {
		s := Serving{
			Kind:             k,
			AmountPerServing: amountPerServing,
		}
		f.Servings = append(f.Servings, s)
	}
}

// WithServing adds s as a saved serving option to a
// food.
//
// s.AmountPerServing is the quantity of the food's base
// unit represented by one serving.
//
// This option does not create a diary entry.
func WithServing(s Serving) Option {
	return func(f *Food) {
		f.Servings = append(f.Servings, s)
	}
}
