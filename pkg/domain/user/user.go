package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/controlado/go-yazio/pkg/domain/date"
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

func (u *User) SinceRegist() date.Range {
	return date.Range{
		Start: u.Registration,
		End:   now(),
	}
}

type Email struct {
	Value       string
	IsConfirmed bool
}

func (e *Email) String() string {
	return strings.ToLower(e.Value)
}
