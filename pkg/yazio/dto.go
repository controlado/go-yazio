package yazio

import (
	"fmt"
	"time"

	"github.com/controlado/go-yazio/internal/infra/client"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
	"github.com/google/uuid"
)

type LoginDTO struct {
	ExpiresInSec int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (d *LoginDTO) toUser(c *client.Client) (*User, error) {
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
			client: c,
			token: &Token{
				expiresAt: timeNow.Add(expiresAt),
				access:    d.AccessToken,
				refresh:   d.RefreshToken,
			},
		}
	)

	return user, nil
}

type GetUserDataDTO struct {
	ID           string `json:"uuid"`
	Token        string `json:"user_token"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	IconURL      string `json:"profile_image"`
	Email        string `json:"email"`
	EmailStatus  string `json:"email_confirmation_status"`
	Registration string `json:"registration_date"`
	BirthDate    string `json:"date_of_birth"`
}

func (d *GetUserDataDTO) toUserData() (u user.Data, err error) {
	parsedID, err := uuid.Parse(d.ID)
	if err != nil {
		return u, fmt.Errorf("parsing user uuid (%q): %w", d.ID, err)
	}

	registTime, err := time.Parse(layoutDate, d.Registration)
	if err != nil {
		return u, fmt.Errorf("parsing regist date (%q): %w", d.Registration, err)
	}

	birthTime, err := time.Parse(layoutISO, d.BirthDate)
	if err != nil {
		return u, fmt.Errorf("parsing user bith date (%q): %w", d.BirthDate, err)
	}

	u = user.Data{
		ID:        parsedID,
		Token:     d.Token,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		IconURL:   d.IconURL,
		Email: user.Email{
			Value:       d.Email,
			IsConfirmed: d.EmailStatus == confirmedEmailStatus,
		},
		Registration: registTime,
		Birth:        birthTime,
	}

	return u, nil
}

type (
	GetMacroIntakeDTO []MacroIntakeDTO
	MacroIntakeDTO    struct {
		Date    string  `json:"date"`
		Energy  float64 `json:"energy"`
		Carb    float64 `json:"carb"`
		Fat     float64 `json:"fat"`
		Protein float64 `json:"protein"`
	}
)

func (d GetMacroIntakeDTO) toRangeMacro() (mr intake.MacrosRange, err error) {
	for i, macroIntake := range d {
		parsedDate, err := time.Parse(layoutISO, macroIntake.Date)
		if err != nil {
			return mr, fmt.Errorf("parsing %d intake date: %w", i, err)
		}

		mi := intake.Macros{
			Date:    parsedDate,
			Energy:  macroIntake.Energy,
			Carb:    macroIntake.Carb,
			Fat:     macroIntake.Fat,
			Protein: macroIntake.Protein,
		}
		mr = append(mr, mi)
	}

	return mr, nil
}

type GetSingleIntakeDTO map[string]float64

func (d GetSingleIntakeDTO) toRangeSingle(k intake.Kind) (sr intake.SingleRange, err error) {
	for date, value := range d {
		parsedDate, err := time.Parse(layoutISO, date)
		if err != nil {
			return nil, fmt.Errorf("parsing single intake date %q: %w", date, err)
		}

		s := intake.Single{
			Kind:  k,
			Date:  parsedDate,
			Value: value,
		}
		sr = append(sr, s)
	}

	return sr, nil
}

type (
	ServingsDTO []ServingDTO
	ServingDTO  struct {
		Type   string  `json:"serving"`
		Amount float64 `json:"amount"`
	}
)

func mapNutrients(nuts map[intake.Kind]float64) map[string]float64 {
	var (
		nutsLength = len(nuts)
		out        = make(map[string]float64, nutsLength)
	)

	if nutsLength < 1 {
		return out
	}

	for kind, value := range nuts {
		nutrientID := kind.ID()
		out[nutrientID] = value
	}

	return out
}

func mapServings(servs []food.Serving) ServingsDTO {
	var (
		servingsLength = len(servs)
		out            = make(ServingsDTO, servingsLength)
	)

	if servingsLength < 1 {
		return out
	}

	for i, s := range servs {
		out[i] = ServingDTO{
			Type:   s.Kind.String(),
			Amount: s.Amount,
		}
	}

	return out
}

func newAddFoodBody(f food.Food, vis visibility.Food) client.Payload[any] {
	return client.Payload[any]{
		"id":         f.ID.String(),
		"name":       f.Name,
		"category":   f.Category.String(),
		"base_unit":  f.BaseUnit.String(),
		"is_private": vis,
		"nutrients":  mapNutrients(f.Nutrients),
		"servings":   mapServings(f.Servings),
	}
}
