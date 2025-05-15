package food

import "github.com/controlado/go-yazio/pkg/domain/unit"

type Option func(f *Food)

func WithBaseUnit(b unit.Base) Option {
	return func(f *Food) {
		f.BaseUnit = b
	}
}

func WithNewServing(k ServingKind, amount float64) Option {
	return func(f *Food) {
		s := Serving{
			Kind:   k,
			Amount: amount,
		}
		f.Servings = append(f.Servings, s)
	}
}

func WithServing(s Serving) Option {
	return func(f *Food) {
		f.Servings = append(f.Servings, s)
	}
}
