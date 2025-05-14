package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	now = time.Now
)

type User struct {
	ID           uuid.UUID
	Token        string
	FirstName    string
	LastName     string
	IconURL      string
	Email        Email
	Registration time.Time
	Birth        time.Time
}

func (u *User) String() string {
	return fmt.Sprintf("User(%v %v)", u.FirstName, u.LastName)
}

func (u *User) SinceRegist() DateRange {
	return DateRange{
		Start: u.Registration,
		End:   now(),
	}
}
