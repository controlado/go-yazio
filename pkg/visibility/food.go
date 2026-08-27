package visibility

type Food bool

const (
	PublicFood  Food = true
	PrivateFood Food = false
)

func (f Food) IsPrivate() bool {
	return f == PrivateFood
}
