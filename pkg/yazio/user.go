package yazio

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/controlado/go-yazio/pkg/client"
	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
)

// User represents an authenticated YAZIO account.
//
// A value of this type keeps the HTTP client used to talk to the
// YAZIO mobile API, the current access token and its expiration
// instant, plus the refresh token needed to renew credentials.
//
// The zero value is not functional; obtain a User through the
// login flow provided in application.API.
type User struct {
	client       *client.Client
	expiresAt    time.Time
	accessToken  string
	refreshToken string
}

// AddFood registers a new food (product) using the account.
//
// AddFood doesn't entry a new intake. Just regist a new food.
//
// On failure the error wraps either:
//   - [ErrRequestingToYazio]
//   - [food.ErrAlreadyExists]
//   - [food.ErrMissingNutrients] f [food.Food] nutrients must have [intake.Energy] [intake.Fat] [intake.Protein] [intake.Carb]
func (u *User) AddFood(ctx context.Context, f food.Food, vis visibility.Food) error {
	requiredNutrients := []intake.Kind{
		intake.Energy,
		intake.Fat,
		intake.Protein,
		intake.Carb,
	}

	for _, k := range requiredNutrients {
		if _, ok := f.Nutrients[k]; !ok {
			return food.ErrMissingNutrients
		}
	}

	var (
		req = client.Request{
			Method:   http.MethodPost,
			Endpoint: addFoodEndpoint,
			Body:     newAddFoodBody(f, vis),
			Headers: client.Payload{
				`accept`:          `*/*`,
				`accept-charset`:  `UTF-8`,
				`accept-encoding`: `application/json`,
				`connection`:      `Keep-Alive`,
				`host`:            `yzapi.yazio.com`,
				`authorization`:   fmt.Sprintf("Bearer %s", u.accessToken),
				`user-agent`:      `YAZIO/12.31.0 (com.yazio.android; build:411052340; Android 34) Ktor`,
			},
		}
	)

	if resp, err := u.client.Request(ctx, req); err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				return food.ErrMissingNutrients
			case http.StatusConflict:
				return food.ErrAlreadyExists
			}
		}
		return fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	return nil
}

// Data retrieves the profile metadata for
// the authenticated user u.
//
// The returned d [user.Data] mirrors the public
// information exposed by the YAZIO API.
//
// On failure the error wraps either:
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
func (u *User) Data(ctx context.Context) (d user.Data, err error) {
	var (
		dto GetUserDataDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: userDataEndpoint,
			Headers: client.Payload{
				`accept`:          `*/*`,
				`accept-charset`:  `UTF-8`,
				`accept-encoding`: `application/json`,
				`connection`:      `Keep-Alive`,
				`host`:            `yzapi.yazio.com`,
				`authorization`:   fmt.Sprintf("Bearer %s", u.accessToken),
				`user-agent`:      `YAZIO/12.31.0 (com.yazio.android; build:411052340; Android 34) Ktor`,
			},
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		return d, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return d, fmt.Errorf("%s: %w", ErrDecodingResponse, err)
	}

	return dto.toUserData()
}

// Intake returns a series of single-nutrient
// intake values for the given date range.
//
// On failure the error wraps either:
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
//   - Other: generic (DTO related)
func (u *User) Intake(ctx context.Context, k intake.Kind, r date.Range) (intake.SingleRange, error) {
	var (
		dto GetSingleIntakeDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: singleIntakesEndpoint,
			Headers: client.Payload{
				`accept`:          `*/*`,
				`accept-charset`:  `UTF-8`,
				`accept-encoding`: `application/json`,
				`connection`:      `Keep-Alive`,
				`host`:            `yzapi.yazio.com`,
				`authorization`:   fmt.Sprintf("Bearer %s", u.accessToken),
				`user-agent`:      `YAZIO/12.31.0 (com.yazio.android; build:411052340; Android 34) Ktor`,
			},
			QueryParams: client.Payload{
				"start":    r.Start.Format(layoutISO),
				"end":      r.End.Format(layoutISO),
				"nutrient": k.String(),
			},
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDecodingResponse, err)
	}

	return dto.toRangeSingle()
}

// Macros returns aggregated values for each
// day within the provided date range:
//
//	[intake.Macros]
//	  - Energy
//	  - Carbohydrate
//	  - Fat
//	  - Protein
//
// On failure the error wraps either:
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
//   - Other: generic (DTO related)
func (u *User) Macros(ctx context.Context, r date.Range) (intake.MacrosRange, error) {
	var (
		dto GetMacroIntakeDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: macrosIntakesEndpoint,
			Headers: client.Payload{
				`accept`:          `*/*`,
				`accept-charset`:  `UTF-8`,
				`accept-encoding`: `application/json`,
				`connection`:      `Keep-Alive`,
				`host`:            `yzapi.yazio.com`,
				`authorization`:   fmt.Sprintf("Bearer %s", u.accessToken),
				`user-agent`:      `YAZIO/12.31.0 (com.yazio.android; build:411052340; Android 34) Ktor`,
			},
			QueryParams: client.Payload{
				"start": r.Start.Format(layoutISO),
				"end":   r.End.Format(layoutISO),
			},
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	return dto.toRangeMacro()
}

// IsExpired reports whether the access token held
// by u has already expired relative to the current
// time.
func (u *User) IsExpired() bool {
	timeNow := time.Now()
	return u.expiresAt.Before(timeNow)
}
