package food

const (
	Meat              Category = "meat"
	Miscellaneous     Category = "miscellaneous"
	Chocolate         Category = "chocolate"
	NonAlcoholicDrink Category = "drinksnonalcoholic"
	Cheese            Category = "cheese"
)

// Category represents the classification of a food item.
type Category string

func (c Category) String() string {
	return string(c)
}
