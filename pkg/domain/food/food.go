package food

import (
	"fmt"

	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/google/uuid"
)

const (
	nameMinLength = 3
)

var (
	defaultServing = Serving{
		Kind:   Portion,
		Amount: 100,
	}
)

type (
	Nutrients map[intake.Kind]float64
	Food      struct {
		ID        uuid.UUID
		Name      string
		BaseUnit  BaseUnit
		Category  Category
		Nutrients Nutrients
		Servings  []Serving
	}
)

func New(name string, cat Category, nut Nutrients, opts ...Option) (f Food, err error) {
	if len(name) < nameMinLength {
		return f, fmt.Errorf("%w: (%q) should have at least 3 chars", ErrInvalidName, name)
	}

	f = Food{
		ID:        uuid.New(),
		Name:      name,
		BaseUnit:  Grams,
		Category:  cat,
		Nutrients: nut,
		Servings:  []Serving{},
	}
	f.apply(opts...)

	if len(f.Servings) == 0 {
		f.Servings = append(f.Servings, defaultServing)
	}

	return f, nil
}

func (f Food) String() string {
	return fmt.Sprintf("Food(%q, %s, %s)",
		f.Name,
		f.ID.String(),
		f.Category.String(),
	)
}

func (f *Food) apply(opts ...Option) {
	for _, o := range opts {
		o(f)
	}
}
