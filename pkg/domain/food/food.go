package food

import (
	"fmt"

	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/unit"
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
	// Nutrients represents a map of nutrient kinds to their respective values.
	//
	// The key is an [intake.Kind] (e.g., protein, carbohydrates, fat, energy)
	// and the value is a float64 representing the amount of that nutrient,
	// typically per 100 units of the food's [BaseUnit] (e.g., 100g or 100ml).
	Nutrients map[intake.Kind]float64

	// Food represents a food item, detailing its identification, nutritional
	// information, category, and available serving sizes.
	//
	// Each food item has a unique ID, a descriptive name, a base unit for
	// its nutrient values (e.g., grams or milliliters), and belongs to a specific
	// food category. The Nutrients field provides a breakdown of its nutritional
	// content, and Servings lists predefined ways the food can be measured or logged.
	Food struct {
		ID        uuid.UUID // ID is the unique identifier for the food item.
		Name      string    // Name is the descriptive name of the food item.
		BaseUnit  unit.Base // BaseUnit specifies the food fundamental unit of measurement.
		Category  Category  // Category classifies the food item; [Meat] [Miscellaneous]...
		Nutrients Nutrients // Nutrients holds the nutritional composition of the food.
		Servings  []Serving // Servings lists the food predefined serving sizes; [Bottle]...
	}
)

// New creates and returns a new [Food] item with a generated ID.
//
// It initializes the food with the provided name, cat [Category], and
// nut [Nutrients]. Optional [Option] functions can be passed to customize
// the food item further, such as setting a specific [BaseUnit] or adding
// custom [Serving] sizes.
//
// If no serving options are provided, a default serving (e.g., 100g portion)
// is automatically added to the food item. The default [BaseUnit] is [Grams]
// unless overridden by an [Option].
//
// On failure the error wraps either:
//   - [ErrInvalidName] if name length is less than 3 characters
func New(name string, cat Category, nut Nutrients, opts ...Option) (f Food, err error) {
	if len(name) < nameMinLength {
		return f, fmt.Errorf("%w: %q should have at least 3 chars", ErrInvalidName, name)
	}

	f = Food{
		ID:        uuid.New(),
		Name:      name,
		BaseUnit:  unit.Gram,
		Category:  cat,
		Nutrients: nut,
		Servings:  []Serving{},
	}
	f.apply(opts...)

	if len(f.Servings) == 0 { // not defined with options
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
