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

type User struct {
	client       *client.Client
	expiresAt    time.Time
	accessToken  string
	refreshToken string
}

// AddFood implements application.User.
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

// Data implements application.User.
func (u *User) Data(ctx context.Context) (d user.User, err error) {
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

// Intake implements application.User.
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

// Macros implements application.User.
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

// IsExpired implements application.User.
func (u *User) IsExpired() bool {
	timeNow := time.Now()
	return u.expiresAt.Before(timeNow)
}
