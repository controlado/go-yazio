package food

const (
	Gram       ServingKind = "gram"
	Milliliter ServingKind = "milliliter"
	Bar        ServingKind = "bar"
	Teaspoon   ServingKind = "teaspoon"
	Tablespoon ServingKind = "tablespoon"
	Glass      ServingKind = "glass"
	Slice      ServingKind = "slice"
	Bottle     ServingKind = "bottle"
	Can        ServingKind = "can"
	Piece      ServingKind = "piece"
	Tablet     ServingKind = "tablet"
	Portion    ServingKind = "portion"
	Cup        ServingKind = "cup"
	Pack       ServingKind = "package"
	Each       ServingKind = "each"
)

// Serving describes a specific quantity or
// measure of a food item.
//
// AmountPerServing represents the quantity
// in [Food.BaseUnit] represented by one
// unit of [Serving.Kind].
type Serving struct {
	Kind             ServingKind
	AmountPerServing float64
}

type ServingKind string

func (sk ServingKind) String() string {
	return string(sk)
}
