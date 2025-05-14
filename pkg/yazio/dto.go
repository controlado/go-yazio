package yazio

import (
	"fmt"
	"time"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/domain"
	"github.com/controlado/go-yazio/pkg/client"
	"github.com/google/uuid"
)

type LoginDTO struct {
	ExpiresInSec int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (d *LoginDTO) toUser(c *client.Client) (application.User, error) {
	switch {
	case d.ExpiresInSec == 0:
		return nil, fmt.Errorf(`zero "expires_in"`)
	case d.AccessToken == "":
		return nil, fmt.Errorf(`blank "access_token"`)
	case d.RefreshToken == "":
		return nil, fmt.Errorf(`blank "refresh_token"`)
	}

	var (
		timeNow   = time.Now()
		expiresAt = time.Duration(d.ExpiresInSec) * time.Second
		user      = &User{
			client:       c,
			expiresAt:    timeNow.Add(expiresAt),
			accessToken:  d.AccessToken,
			refreshToken: d.RefreshToken,
		}
	)

	return user, nil
}

type GetUserDataDTO struct {
	ID string `json:"uuid"`

	Token     string `json:"user_token"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IconURL   string `json:"profile_image"`

	Email       string `json:"email"`
	EmailStatus string `json:"email_confirmation_status"`

	Registration string `json:"registration_date"`
	BirthDate    string `json:"date_of_birth"`
}

func (d *GetUserDataDTO) toUserData() (u domain.User, err error) {
	parsedID, err := uuid.Parse(d.ID)
	if err != nil {
		return u, fmt.Errorf("parsing user uuid (%q): %w", d.ID, err)
	}

	birthTime, err := time.Parse(layoutISO, d.BirthDate)
	if err != nil {
		return u, fmt.Errorf("parsing user bith date (%q): %w", d.BirthDate, err)
	}

	registTime, err := time.Parse(layoutDate, d.Registration)
	if err != nil {
		return u, fmt.Errorf("parsing regist date (%q): %w", d.Registration, err)
	}

	u = domain.User{
		ID:        parsedID,
		Token:     d.Token,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		IconURL:   d.IconURL,
		Email: domain.Email{
			Value:       d.Email,
			IsConfirmed: d.EmailStatus == confirmedEmailStatus,
		},
		Registration: registTime,
		Birth:        birthTime,
	}

	return u, nil
}

type MacroIntakeDTO struct {
	Date    string  `json:"date"`
	Energy  float64 `json:"energy"`
	Carb    float64 `json:"carb"`
	Fat     float64 `json:"fat"`
	Protein float64 `json:"protein"`
}

type GetMacroIntakeDTO []MacroIntakeDTO

func (d GetMacroIntakeDTO) toRangeMacro() (rm domain.MacrosRange, err error) {
	for i, intake := range d {
		parsedDate, err := time.Parse(layoutISO, intake.Date)
		if err != nil {
			return rm, fmt.Errorf("parsing %d intake date: %w", i, err)
		}

		mi := domain.MacrosIntake{
			Date:    parsedDate,
			Energy:  intake.Energy,
			Carb:    intake.Carb,
			Fat:     intake.Fat,
			Protein: intake.Protein,
		}
		rm = append(rm, mi)
	}

	return rm, nil
}

type GetSingleIntakeDTO map[string]float64

func (d GetSingleIntakeDTO) toRangeSingle() (rs domain.SingleRange, err error) {
	for date, value := range d {
		parsedDate, err := time.Parse(layoutISO, date)
		if err != nil {
			return nil, fmt.Errorf("parsing single intake date %q: %w", date, err)
		}

		si := domain.SingleIntake{
			Date:  parsedDate,
			Value: value,
		}
		rs = append(rs, si)
	}

	return rs, nil
}
