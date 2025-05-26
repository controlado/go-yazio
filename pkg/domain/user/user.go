package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/google/uuid"
)

// Data represents the detailed profile information of a YAZIO user.
//
// It includes unique identifiers, personal details, authentication tokens,
// and important dates related to the user's account.
type Data struct {
	ID           uuid.UUID // ID is the unique identifier for the user.
	Token        string    // Token is an unknown (for now) token associated with user.
	FirstName    string    // FirstName is the user's first name.
	LastName     string    // LastName is the user's last name.
	IconURL      string    // IconURL is the URL pointing to the user's profile picture.
	Email        Email     // Email holds the user's email informations.
	Registration time.Time // Registration is the timestamp when the user registered their account.
	Birth        time.Time // Birth is the user's date of birth.
}

func (d *Data) String() string {
	return fmt.Sprintf("User(%v %v)", d.FirstName, d.LastName)
}

// SinceRegist calculates and returns a [date.Range]
// representing the period from the user's registration
// time up to the current moment.
func (d *Data) SinceRegist() date.Range {
	var (
		timeNow = time.Now()
	)
	return d.SinceRegistAt(timeNow)
}

func (d *Data) SinceRegistAt(until time.Time) date.Range {
	return date.Range{Start: d.Registration, End: until}
}

type Email struct {
	Value       string // Value is the actual email address string.
	IsConfirmed bool
}

func (e *Email) String() string {
	return strings.ToLower(e.Value)
}
