package food

const (
	Grams      BaseUnit = "g"
	Milliliter BaseUnit = "ml"
)

type BaseUnit string

func (bu BaseUnit) String() string {
	return string(bu)
}
