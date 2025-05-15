package food

const (
	Grams      BaseUnit = "g"
	Milliliter BaseUnit = "ml"
)

// BaseUnit represents the food unit of measurement.
type BaseUnit string

func (bu BaseUnit) String() string {
	return string(bu)
}
