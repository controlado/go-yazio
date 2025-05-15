package unit

const (
	Kilocalorie Base = "kcal"
	Milliliter  Base = "ml"
	Gram        Base = "g"
	Milligram   Base = "mg"
	Microgram   Base = "mcg"
)

// Base represents the food unit of measurement.
type Base string

func (b Base) String() string {
	return string(b)
}
