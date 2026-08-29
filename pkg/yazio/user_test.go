package yazio

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/controlado/go-yazio/v3/internal/infra/client"
	"github.com/controlado/go-yazio/v3/internal/testutil/assert"
	"github.com/controlado/go-yazio/v3/internal/testutil/server"
	"github.com/controlado/go-yazio/v3/internal/testutil/times"
	"github.com/controlado/go-yazio/v3/pkg/domain/date"
	"github.com/controlado/go-yazio/v3/pkg/domain/food"
	"github.com/controlado/go-yazio/v3/pkg/domain/intake"
	"github.com/controlado/go-yazio/v3/pkg/domain/meal"
	"github.com/controlado/go-yazio/v3/pkg/domain/unit"
	"github.com/controlado/go-yazio/v3/pkg/domain/user"
	"github.com/controlado/go-yazio/v3/pkg/visibility"
	"github.com/google/uuid"
)

func TestUser_Data(t *testing.T) {
	t.Parallel()

	var (
		staticID = uuid.New()
	)

	srv, err := server.New(t,
		server.AssertMethod(http.MethodGet),
		server.AssertEndpoint(userDataEndpoint),
		server.RespondBodyAny(map[string]string{
			"uuid":                      staticID.String(),
			"user_token":                "c000a7769600a98abae7cefe56174e48240ee297e06be3052cc3e743f12bcfd5",
			"first_name":                "João Brito",
			"last_name":                 "da Silva",
			"profile_image":             "https://images.yazio-cdn.com/process/plain/app/profile/user/2025/d297247d-51d4-4e04-9e87-c99fdf693585.jpg",
			"email":                     "joaodasilva@gmail.com",
			"email_confirmation_status": "confirmed",
			"registration_date":         "2023-02-06 21:22:46",
			"date_of_birth":             "2005-08-26",
		}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	c := client.New(
		client.WithBaseURL(srv.URL),
	)
	assert.NotNil(t, c)

	var (
		ctx = context.Background()
		u   = User{
			client: c,
			token: &Token{
				expiresAt: times.Future(),
				access:    uuid.NewString(),
				refresh:   uuid.NewString(),
			},
		}
	)

	userData, err := u.Data(ctx)
	assert.NoError(t, err)

	want := user.Data{
		ID:        staticID,
		Token:     "c000a7769600a98abae7cefe56174e48240ee297e06be3052cc3e743f12bcfd5",
		FirstName: "João Brito",
		LastName:  "da Silva",
		IconURL:   "https://images.yazio-cdn.com/process/plain/app/profile/user/2025/d297247d-51d4-4e04-9e87-c99fdf693585.jpg",
		Email: user.Email{
			Value:       "joaodasilva@gmail.com",
			IsConfirmed: true,
		},
		Registration: time.Date(2023, 02, 06, 21, 22, 46, 0, time.UTC),
		Birth:        time.Date(2005, 8, 26, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, want, userData)
}

func TestUser_Macros(t *testing.T) {
	t.Parallel()

	startDate, err := time.Parse(layoutISO, "2025-04-12")
	assert.NoError(t, err)

	endDate, err := time.Parse(layoutISO, "2025-04-13")
	assert.NoError(t, err)

	srv, err := server.New(t,
		server.AssertEndpoint(macrosIntakesEndpoint),
		server.AssertQueryParams(map[string]string{
			"start": "2025-04-12",
			"end":   "2025-04-13",
		}),
		server.RespondBodyAny([]map[string]any{
			{
				"date":        "2025-04-12",
				"energy":      1288.68,
				"carb":        85.37,
				"protein":     94.38,
				"fat":         62.17,
				"energy_goal": 1935,
			},
			{
				"date":        "2025-04-13",
				"energy":      1768.78,
				"carb":        156.53,
				"protein":     182.95,
				"fat":         38.76,
				"energy_goal": 1800,
			},
		}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	c := client.New(
		client.WithBaseURL(srv.URL),
	)
	assert.NotNil(t, c)

	var (
		ctx        = context.Background()
		macroRange = date.Range{Start: startDate, End: endDate}
		u          = User{
			client: c,
			token: &Token{
				expiresAt: times.Future(),
				access:    uuid.NewString(),
				refresh:   uuid.NewString(),
			},
		}
	)
	rm, err := u.Macros(ctx, macroRange)
	assert.NoError(t, err)

	want := intake.MacrosRange{
		{Date: startDate, Energy: 1288.68, Carb: 85.37, Fat: 62.17, Protein: 94.38},
		{Date: endDate, Energy: 1768.78, Carb: 156.53, Fat: 38.76, Protein: 182.95},
	}
	assert.EqualSlicesItems(t, rm, want)
}

func TestUser_Intake(t *testing.T) {
	t.Parallel()

	var (
		startDate = times.PastDate(time.UTC)
		endDate   = times.FutureDate(time.UTC)
	)

	type args struct {
		ctx       context.Context
		kind      intake.Kind
		dateRange date.Range
	}

	var (
		ctx        = context.Background()
		testBlocks = []struct {
			name    string
			wantErr bool
			args    args
			want    intake.SingleRange
		}{
			{
				name: "valid args",
				args: args{
					ctx:       ctx,
					kind:      intake.Sugar,
					dateRange: date.Range{Start: startDate, End: endDate},
				},
				want: intake.SingleRange{
					{Kind: intake.Sugar, Date: startDate, Value: 89.83},
					{Kind: intake.Sugar, Date: endDate, Value: 50.38},
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			serverBody := make(map[string]any, len(tb.want))
			for _, si := range tb.want {
				d := si.Date.Format(layoutISO)
				serverBody[d] = si.Value
			}

			srv, err := server.New(t,
				server.AssertMethod(http.MethodGet),
				server.AssertEndpoint(singleIntakesEndpoint),
				server.AssertQueryParams(
					map[string]string{
						"start":    tb.args.dateRange.Start.Format(layoutISO),
						"end":      tb.args.dateRange.End.Format(layoutISO),
						"nutrient": tb.args.kind.ID(),
					},
				),
				server.RespondBodyAny(serverBody),
			)
			assert.NoError(t, err)
			assert.NotNil(t, srv)

			c := client.New(
				client.WithBaseURL(srv.URL),
			)
			assert.NotNil(t, c)

			u := User{
				client: c,
				token: &Token{
					expiresAt: times.Future(),
					access:    uuid.NewString(),
					refresh:   uuid.NewString(),
				},
			}
			rm, err := u.Intake(
				tb.args.ctx,
				tb.args.kind,
				tb.args.dateRange,
			)
			assert.NoError(t, err)
			assert.EqualSlicesItems(t, rm, tb.want)
		})
	}
}

func TestUser_AddFood(t *testing.T) {
	t.Parallel()

	var ( // static data
		validFood = food.Food{
			ID:                      uuid.New(),
			Name:                    "banana",
			BaseUnit:                unit.Gram,
			Category:                food.Miscellaneous,
			NutrientReferenceAmount: 100,
			Nutrients: food.Nutrients{
				intake.Energy:  10,
				intake.Fat:     10,
				intake.Protein: 10,
				intake.Carb:    10,
			},
			Servings: []food.Serving{
				{
					Kind:             food.Piece,
					AmountPerServing: 1,
				},
			},
		}
		nutrientsPerUnit = map[string]any{
			"energy.energy":    0.1,
			"nutrient.fat":     0.1,
			"nutrient.protein": 0.1,
			"nutrient.carb":    0.1,
		}
	)

	var (
		ctx        = context.Background()
		testBlocks = []struct {
			name           string
			wantErr        bool
			serverStatus   int // default (success): StatusNoContent
			food           food.Food
			foodVisibility visibility.Food
			wantNutrients  map[string]any
		}{
			{
				name:           "valid public food",
				food:           validFood,
				foodVisibility: visibility.PublicFood,
				wantNutrients:  nutrientsPerUnit,
			},
			{
				name: "valid food using custom nutrient reference amount",
				food: func() food.Food {
					customReferenceFood := validFood
					customReferenceFood.NutrientReferenceAmount = 18
					customReferenceFood.Nutrients = food.Nutrients{
						intake.Energy:  18,
						intake.Fat:     18,
						intake.Protein: 18,
						intake.Carb:    18,
					}
					return customReferenceFood
				}(),
				foodVisibility: visibility.PrivateFood,
				wantNutrients: map[string]any{
					"energy.energy":    1.0,
					"nutrient.fat":     1.0,
					"nutrient.protein": 1.0,
					"nutrient.carb":    1.0,
				},
			},
			{
				name:    "food missing nutrients",
				wantErr: true,
				food: func() food.Food {
					invalidFood := validFood
					invalidFood.Nutrients = food.Nutrients{}
					return invalidFood
				}(),
				foodVisibility: visibility.PrivateFood,
			},
			{
				name:           "server respond - StatusBadRequest",
				wantErr:        true,
				serverStatus:   http.StatusBadRequest,
				food:           validFood,
				foodVisibility: visibility.PrivateFood,
				wantNutrients:  nutrientsPerUnit,
			},
			{
				name:           "server respond - StatusConflict",
				wantErr:        true,
				serverStatus:   http.StatusConflict,
				food:           validFood,
				foodVisibility: visibility.PrivateFood,
				wantNutrients:  nutrientsPerUnit,
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			if tb.serverStatus == 0 {
				tb.serverStatus = http.StatusNoContent
			}

			srv, err := server.New(t,
				server.AssertEndpoint(addFoodEndpoint),
				server.AssertMethod(http.MethodPost),
				server.AssertBody(map[string]any{
					"id":         tb.food.ID.String(),
					"name":       tb.food.Name,
					"category":   tb.food.Category.String(),
					"base_unit":  tb.food.BaseUnit.String(),
					"is_private": tb.foodVisibility.IsPrivate(),
					"nutrients":  tb.wantNutrients,
					"servings": []any{
						map[string]any{
							"serving": "piece",
							"amount":  1.0,
						},
					},
				}),
				server.RespondHeaders(
					map[string]string{
						"Cache-Control": "no-cache, private",
						"Connection":    "keep-alive",
						"Date":          "Wed, 14 May 2025 19:03:19 GMT",
						"Server":        "nginx",
					},
				),
				server.RespondStatus(tb.serverStatus),
			)
			assert.NoError(t, err)
			assert.NotNil(t, srv)

			c := client.New(
				client.WithBaseURL(srv.URL),
			)
			assert.NotNil(t, c)

			u := &User{
				client: c,
				token: &Token{
					expiresAt: times.Future(),
					access:    uuid.NewString(),
					refresh:   uuid.NewString(),
				},
			}
			err = u.AddFood(
				ctx,
				tb.food,
				tb.foodVisibility,
			)
			assert.WantErr(t, tb.wantErr, err)
		})
	}
}

func TestUser_AddFood_InvalidNutrientReferenceAmount(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name   string
		amount float64
	}{
		{name: "zero", amount: 0},
		{name: "negative", amount: -1},
		{name: "NaN", amount: math.NaN()},
		{name: "infinity", amount: math.Inf(1)},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			srv, err := server.New(t, server.AssertRequest(false))
			assert.NoError(t, err)

			u := &User{
				client: client.New(client.WithBaseURL(srv.URL)),
				token:  &Token{expiresAt: times.Future()},
			}
			f := food.Food{
				NutrientReferenceAmount: tb.amount,
				Nutrients: food.Nutrients{
					intake.Energy:  1,
					intake.Fat:     1,
					intake.Protein: 1,
					intake.Carb:    1,
				},
			}

			err = u.AddFood(context.Background(), f, visibility.PrivateFood)
			assert.Equal(t, err, food.ErrInvalidNutrientReferenceAmount)
		})
	}
}

func TestUser_EntryFood(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx             context.Context
		meal            meal.Time
		id              food.ID
		serving         food.Serving
		servingQuantity float64
	}

	var ( // static
		validFoodID  = uuid.New()
		validServing = food.Serving{Kind: food.Portion, AmountPerServing: 100}
		validToken   = &Token{expiresAt: times.Future(), access: "valid-access", refresh: "valid-refresh"}
		expiredToken = &Token{expiresAt: times.Past(), access: "invalid-access", refresh: "invalid-refresh"}
	)

	testBlocks := []struct {
		name          string
		wantErr       bool
		wantNoRequest bool
		respondStatus int
		token         *Token
		a             args
	}{
		{
			name:  "valid path",
			token: validToken,
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         validServing,
				servingQuantity: 1,
			},
		},
		{
			name:          "using a blank serving.Kind should explode",
			wantErr:       true,
			wantNoRequest: true,
			token:         validToken,
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         food.Serving{Kind: "", AmountPerServing: 50},
				servingQuantity: 1,
			},
		},
		{
			name:          "using a invalid serving.AmountPerServing (0) should explode",
			wantErr:       true,
			wantNoRequest: true,
			token:         validToken,
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         food.Serving{Kind: food.Gram, AmountPerServing: 0},
				servingQuantity: 1,
			},
		},
		{
			name:          "using a invalid servingQuantity (0) should explode",
			wantErr:       true,
			wantNoRequest: true,
			token:         validToken,
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         validServing,
				servingQuantity: 0,
			},
		},
		{
			name:          "expired token",
			wantErr:       true,
			wantNoRequest: true,
			token:         expiredToken,
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         validServing,
				servingQuantity: 1,
			},
		},
		{
			name:          "server -> http.StatusUnauthorized: invalid token",
			wantErr:       true,
			respondStatus: http.StatusUnauthorized,
			token:         validToken, // to pass first check
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         validServing,
				servingQuantity: 1,
			},
		},
		{
			name:          "server -> http.StatusConflict: action uuid already exists",
			wantErr:       true,
			respondStatus: http.StatusConflict,
			token:         validToken,
			a: args{
				ctx:             context.Background(),
				meal:            meal.Dinner,
				id:              validFoodID,
				serving:         validServing,
				servingQuantity: 1,
			},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			if tb.respondStatus == 0 {
				tb.respondStatus = http.StatusNoContent
			}

			srv, err := server.New(t,
				server.AssertMethod(http.MethodPost),
				server.AssertEndpoint("/v20/user/consumed-items"),
				server.AssertRequest(!tb.wantNoRequest),
				server.RespondStatus(tb.respondStatus),
			)
			assert.NoError(t, err)
			assert.NotNil(t, srv)

			u := &User{
				token: tb.token,
				client: client.New(
					client.WithBaseURL(srv.URL),
				),
			}
			err = u.EntryFood(
				tb.a.ctx,
				tb.a.meal,
				tb.a.id,
				tb.a.serving,
				tb.a.servingQuantity,
			)
			assert.WantErr(t, tb.wantErr, err)
		})
	}
}

func TestIsPositiveFinite(t *testing.T) {
	t.Parallel()

	testBlocks := []struct {
		name  string
		value float64
		want  bool
	}{
		{name: "positive", value: 1, want: true},
		{name: "zero", value: 0},
		{name: "NaN", value: math.NaN()},
		{name: "positive infinity", value: math.Inf(1)},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, isPositiveFinite(tb.value), tb.want)
		})
	}
}

func TestNewEntryFoodBody(t *testing.T) {
	t.Parallel()

	var (
		entryID = uuid.New()
		foodID  = uuid.New()
		at      = time.Date(2026, 8, 26, 22, 9, 12, 0, time.Local)

		servingCases = []struct {
			servingKind      food.ServingKind
			amountPerServing float64
			servingQuantity  float64
			wantAmount       float64
		}{
			{food.Gram, 1, 58, 58},
			{food.Slice, 10, 2, 20},
			{food.Cup, 237, 1.5, 355.5},
			{food.Portion, 100, 2, 200},
			{food.Portion, 100, 0.5, 50},
			{food.Milliliter, 1, 250, 250},
		}
	)

	for _, sc := range servingCases {
		t.Run(fmt.Sprintf("%s %f x %f = %f", sc.servingKind, sc.amountPerServing, sc.servingQuantity, sc.wantAmount), func(t *testing.T) {
			t.Parallel()

			var (
				serving = food.Serving{Kind: sc.servingKind, AmountPerServing: sc.amountPerServing}
				got     = newEntryFoodBody(entryID, at, meal.Lunch, foodID, serving, sc.servingQuantity)
				want    = client.Payload[any]{
					"products": []map[string]any{
						{
							"id":               entryID,
							"date":             at.Format(layoutDate),
							"daytime":          meal.Lunch,
							"product_id":       foodID,
							"serving":          sc.servingKind,
							"serving_quantity": sc.servingQuantity,
							"amount":           sc.wantAmount,
						},
					},
					"simple_products": []any{},
					"recipe_portions": []any{},
				}
			)
			assert.DeepEqual(t, got, want)
		})
	}
}
