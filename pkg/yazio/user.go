package yazio

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/infra/client"
	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/meal"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
	"github.com/google/uuid"
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
	client *client.Client
	token  application.Token
}

// Token returns the [application.Token] held by u.
func (u *User) Token() application.Token {
	return u.token
}

// EntryFood adds a food-intake to the authenticated user's diary.
//
// It always targets today, saving the [meal.Time] [food.ID] [food.Serving].
//
//   - YAZIO does not validate the product ID:
//
//     If foodID doesn't point to a valid food, the request succeeds
//     but the entry is silently discarded. Ensure the product exists
//     before invoking this method.
//
// On failure the error wraps either:
//   - [ErrExpiredToken]
//   - [ErrRequestingToYazio]
func (u *User) EntryFood(ctx context.Context, mealTime meal.Time, foodID food.ID, serving food.Serving) error {
	if u.token.IsExpired() {
		return ErrExpiredToken
	}

	var (
		req = client.Request{
			Method:   http.MethodPost,
			Endpoint: entryFoodEndpoint,
			Headers:  defaultHeaders(u.token),
			Body: client.Payload[any]{
				"products": []map[string]any{
					{
						"id":               uuid.New(),
						"date":             time.Now().Format(layoutDate),
						"daytime":          mealTime,
						"product_id":       foodID,
						"serving":          serving.Kind,
						"amount":           serving.Amount,
						"serving_quantity": 1,
					},
				},
				"simple_products": []any{},
				"recipe_portions": []any{},
			},
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return ErrExpiredToken
			case http.StatusConflict:
				// theoretically it's not possible because
				// we generate a uuid for each call.
				return food.ErrAlreadyExists
			}
		}
		return fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	return nil
}

// AddFood registers a new food (product) using the account.
//
// AddFood doesn't entry a new intake. Just regist a new food.
//
// On failure the error wraps either:
//   - [ErrExpiredToken]
//   - [ErrRequestingToYazio]
//   - [food.ErrAlreadyExists]
//   - [food.ErrMissingNutrients] f [food.Food] nutrients must have [intake.Energy] [intake.Fat] [intake.Protein] [intake.Carb]
func (u *User) AddFood(ctx context.Context, f food.Food, vis visibility.Food) error {
	if u.token.IsExpired() {
		return ErrExpiredToken
	}

	requiredNutrients := []intake.Kind{
		intake.Energy, intake.Fat,
		intake.Protein, intake.Carb,
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
			Headers:  defaultHeaders(u.token),
		}
	)

	if resp, err := u.client.Request(ctx, req); err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				return food.ErrMissingNutrients
			case http.StatusUnauthorized:
				return ErrExpiredToken
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
//   - [ErrExpiredToken]
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
func (u *User) Data(ctx context.Context) (d user.Data, err error) {
	if u.token.IsExpired() {
		return d, ErrExpiredToken
	}

	var (
		dto GetUserDataDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: userDataEndpoint,
			Headers:  defaultHeaders(u.token),
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return d, ErrExpiredToken
			}
		}
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
//   - [ErrExpiredToken]
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
//   - Other: generic (DTO related)
func (u *User) Intake(ctx context.Context, k intake.Kind, r date.Range) (intake.SingleRange, error) {
	if u.token.IsExpired() {
		return nil, ErrExpiredToken
	}

	var (
		dto GetSingleIntakeDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: singleIntakesEndpoint,
			Headers:  defaultHeaders(u.token),
			QueryParams: client.Payload[string]{
				"start":    r.Start.Format(layoutISO),
				"end":      r.End.Format(layoutISO),
				"nutrient": k.ID(),
			},
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return nil, ErrExpiredToken
			}
		}
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDecodingResponse, err)
	}

	return dto.toRangeSingle(k)
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
//   - [ErrExpiredToken]
//   - [ErrRequestingToYazio]
//   - [ErrDecodingResponse]
//   - Other: generic (DTO related)
func (u *User) Macros(ctx context.Context, r date.Range) (intake.MacrosRange, error) {
	if u.token.IsExpired() {
		return nil, ErrExpiredToken
	}

	var (
		dto GetMacroIntakeDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: macrosIntakesEndpoint,
			Headers:  defaultHeaders(u.token),
			QueryParams: client.Payload[string]{
				"start": r.Start.Format(layoutISO),
				"end":   r.End.Format(layoutISO),
			},
		}
	)

	resp, err := u.client.Request(ctx, req)
	if err != nil {
		if resp.Response != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return nil, ErrExpiredToken
			}
		}
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrRequestingToYazio, err)
	}

	return dto.toRangeMacro()
}
