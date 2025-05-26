package meal

const (
	Breakfast Time = "breakfast"
	Lunch     Time = "lunch"
	Dinner    Time = "dinner"
	Snack     Time = "snack"
)

// Time is the meal time.
type Time string

func (t Time) String() string {
	return string(t)
}
