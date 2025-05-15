package food

type Option func(f *Food)

func WithBaseUnit(bu BaseUnit) Option {
	return func(f *Food) {
		f.BaseUnit = bu
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
