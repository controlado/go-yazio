package yazio

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/infra/client"
	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
	"github.com/controlado/go-yazio/internal/testutil/times"
	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/meal"
	"github.com/controlado/go-yazio/pkg/domain/unit"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
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
			ID:       uuid.New(),
			Name:     "banana",
			BaseUnit: unit.Gram,
			Category: food.Miscellaneous,
			Nutrients: food.Nutrients{
				intake.Energy:  10,
				intake.Fat:     10,
				intake.Protein: 10,
				intake.Carb:    10,
			},
			Servings: []food.Serving{
				{
					Kind:   food.Piece,
					Amount: 1,
				},
			},
		}
	)

	var (
		ctx        = context.Background()
		testBlocks = []struct {
			name         string
			wantErr      bool
			serverStatus int // default (success): StatusNoContent
			food         food.Food
		}{
			{
				name: "valid food",
				food: validFood,
			},
			{
				name:    "food missing nutrients",
				wantErr: true,
				food: func() food.Food {
					invalidFood := validFood
					invalidFood.Nutrients = food.Nutrients{}
					return invalidFood
				}(),
			},
			{
				name:         "server respond - StatusBadRequest",
				wantErr:      true,
				serverStatus: http.StatusBadRequest,
				food:         validFood,
			},
			{
				name:         "server respond - StatusConflict",
				wantErr:      true,
				serverStatus: http.StatusConflict,
				food:         validFood,
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
				visibility.PrivateFood,
			)
			assert.WantErr(t, tb.wantErr, err)
		})
	}
}

func TestUser_EntryFood(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		meal    meal.Time
		id      food.ID
		serving food.Serving
	}

	var ( // static
		validFoodID  = uuid.New()
		validServing = food.Serving{Kind: food.Portion, Amount: 100}
		validToken   = &Token{expiresAt: times.Future(), access: "valid-access", refresh: "valid-refresh"}
		expiredToken = &Token{expiresAt: times.Past(), access: "invalid-access", refresh: "invalid-refresh"}
	)

	testBlocks := []struct {
		name          string
		wantErr       bool
		respondStatus int
		token         *Token
		a             args
	}{
		{
			name:  "valid path",
			token: validToken,
			a: args{
				ctx:     context.Background(),
				meal:    meal.Dinner,
				id:      validFoodID,
				serving: validServing,
			},
		},
		{
			name:    "expired token",
			wantErr: true,
			token:   expiredToken,
			a: args{
				ctx:     context.Background(),
				meal:    meal.Dinner,
				id:      validFoodID,
				serving: validServing,
			},
		},
		{
			name:          "server -> http.StatusUnauthorized: invalid token",
			wantErr:       true,
			respondStatus: http.StatusUnauthorized,
			token:         validToken, // to pass first check
			a: args{
				ctx:     context.Background(),
				meal:    meal.Dinner,
				id:      validFoodID,
				serving: validServing,
			},
		},
		{
			name:          "server -> http.StatusConflict: action uuid already exists",
			wantErr:       true,
			respondStatus: http.StatusConflict,
			token:         validToken,
			a: args{
				ctx:     context.Background(),
				meal:    meal.Dinner,
				id:      validFoodID,
				serving: validServing,
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
				server.AssertEndpoint("/v18/user/consumed-items"),
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
			)
			assert.WantErr(t, tb.wantErr, err)
		})
	}
}
