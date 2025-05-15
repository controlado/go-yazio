package food

const (
	Meat          Category = "meat"
	Miscellaneous Category = "miscellaneous"
)

type Category string

func (c Category) String() string {
	return string(c)
}
