package serving

import "github.com/controlado/go-yazio/v5/pkg/domain/food"

// Gram returns a [food.Serving] representing one gram.
func Gram() food.Serving {
	return food.Serving{Kind: food.Gram, AmountPerServing: 1}
}

// Milliliter returns a [food.Serving] representing one milliliter.
func Milliliter() food.Serving {
	return food.Serving{Kind: food.Milliliter, AmountPerServing: 1}
}

func Bar(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Bar, AmountPerServing: amountPerServing}
}

func Teaspoon(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Teaspoon, AmountPerServing: amountPerServing}
}

func Tablespoon(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Tablespoon, AmountPerServing: amountPerServing}
}

func Glass(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Glass, AmountPerServing: amountPerServing}
}

func Slice(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Slice, AmountPerServing: amountPerServing}
}

func Bottle(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Bottle, AmountPerServing: amountPerServing}
}

func Can(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Can, AmountPerServing: amountPerServing}
}

func Piece(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Piece, AmountPerServing: amountPerServing}
}

func Tablet(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Tablet, AmountPerServing: amountPerServing}
}

func Portion(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Portion, AmountPerServing: amountPerServing}
}

func Cup(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Cup, AmountPerServing: amountPerServing}
}

func Pack(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Pack, AmountPerServing: amountPerServing}
}

func Each(amountPerServing float64) food.Serving {
	return food.Serving{Kind: food.Each, AmountPerServing: amountPerServing}
}
