package yazio

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/controlado/go-yazio/v4/internal/infra/client"
	"github.com/controlado/go-yazio/v4/pkg/domain/date"
	"github.com/controlado/go-yazio/v4/pkg/domain/food"
	"github.com/controlado/go-yazio/v4/pkg/domain/intake"
	"github.com/controlado/go-yazio/v4/pkg/domain/meal"
	"github.com/controlado/go-yazio/v4/pkg/domain/user"
	"github.com/controlado/go-yazio/v4/pkg/visibility"
	"github.com/google/uuid"
)

// User represents an authenticated YAZIO account.
//
// A value of this type keeps the HTTP client used to talk to the
// YAZIO mobile API, the current access token and its expiration
// instant, plus the refresh token needed to renew credentials.
//
// The zero value is not functional; obtain a User through the
// login flow provided in [API].
type User struct {
	client *client.Client
	token  *Token
}

// Token returns the [Token] held by u.
func (u *User) Token() *Token {
	return u.token
}

// EntryFood logs a food in today's diary.
//
// serving describes the selected measure and how much it represents in the
// food's base unit. servingQuantity is the number of those measures consumed.
// For example, two slices of 10 g are logged as 20 g.
//
//   - YAZIO does not validate the product ID:
//
//     If foodID doesn't point to a valid food, the request succeeds
//     but the entry is silently discarded. Ensure the product exists
//     before invoking this method.
//
// On failure the error wraps either:
//   - [ErrExpiredToken]
//   - [food.ErrInvalidServing]
//   - [food.ErrInvalidServingQuantity]
//   - [ErrRequestingToYazio]
func (u *User) EntryFood(ctx context.Context, mealTime meal.Time, foodID food.ID, serving food.Serving, servingQuantity float64) error {
	if u.token.IsExpired() {
		return ErrExpiredToken
	}
	if serving.Kind == "" || !isPositiveFinite(serving.AmountPerServing) {
		return food.ErrInvalidServing
	}
	if !isPositiveFinite(servingQuantity) {
		return food.ErrInvalidServingQuantity
	}

	var (
		req = client.Request{
			Method:   http.MethodPost,
			Endpoint: entryFoodEndpoint,
			Headers:  defaultHeaders(u.token),
			Body: newEntryFoodBody(
				uuid.New(),
				time.Now(),
				mealTime,
				foodID,
				serving,
				servingQuantity,
			),
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
		return fmt.Errorf("%w: %w", ErrRequestingToYazio, err)
	}

	return nil
}

func isPositiveFinite(v float64) bool {
	return v > 0 && !math.IsInf(v, 0)
}

func newEntryFoodBody(entryID uuid.UUID, at time.Time, mealTime meal.Time, foodID food.ID, serving food.Serving, servingQuantity float64) client.Payload[any] {
	return client.Payload[any]{
		"products": []map[string]any{
			{
				"id":               entryID,
				"date":             at.Format(layoutDate),
				"daytime":          mealTime,
				"product_id":       foodID,
				"serving":          serving.Kind,
				"serving_quantity": servingQuantity,
				"amount":           serving.AmountPerServing * servingQuantity,
			},
		},
		"simple_products": []any{},
		"recipe_portions": []any{},
	}
}

// AddFood registers a new food (product) using the account.
//
// A successful call confirms that the food was registered on YAZIO's servers.
// Some versions of the YAZIO app keep a local product catalog and may not show
// foods created by external clients in search until that catalog is refreshed.
// Creating a product in the app was observed to trigger this refresh. There is
// no known API endpoint that forces it.
//
// AddFood does not create an intake entry. Use [User.EntryFood] to log the food.
//
// On failure the error wraps either:
//   - [ErrExpiredToken]
//   - [ErrRequestingToYazio]
//   - [food.ErrAlreadyExists]
//   - [food.ErrMissingNutrients] f [food.Food] nutrients must have [intake.Energy] [intake.Fat] [intake.Protein] [intake.Carb]
//   - [food.ErrInvalidNutrientReferenceAmount]
func (u *User) AddFood(ctx context.Context, f food.Food, vis visibility.Food) error {
	if u.token.IsExpired() {
		return ErrExpiredToken
	}
	if !isPositiveFinite(f.NutrientReferenceAmount) {
		return food.ErrInvalidNutrientReferenceAmount
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
		return fmt.Errorf("%w: %w", ErrRequestingToYazio, err)
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
		dto getUserDataDTO
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
		return d, fmt.Errorf("%w: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return d, fmt.Errorf("%w: %w", ErrDecodingResponse, err)
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
		dto getSingleIntakeDTO
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
		return nil, fmt.Errorf("%w: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodingResponse, err)
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
		dto getMacroIntakeDTO
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
		return nil, fmt.Errorf("%w: %w", ErrRequestingToYazio, err)
	}

	if err := resp.BodyStruct(&dto); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodingResponse, err)
	}

	return dto.toRangeMacro()
}
