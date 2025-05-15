package yazio

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
	"github.com/controlado/go-yazio/internal/testutil/times"
	"github.com/controlado/go-yazio/pkg/client"
	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/unit"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
	"github.com/google/uuid"
)

func TestUser_Data(t *testing.T) {
	t.Parallel()

	var (
		fakeID = uuid.New()
		want   = user.Data{
			ID:        fakeID,
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
	)

	handler := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		assert.Equal(t, r.URL.Path, userDataEndpoint)

		respBody := GetUserDataDTO{
			ID:           fakeID.String(),
			Token:        want.Token,
			FirstName:    want.FirstName,
			LastName:     want.LastName,
			IconURL:      want.IconURL,
			Email:        want.Email.Value,
			EmailStatus:  confirmedEmailStatus,
			Registration: "2023-02-06 21:22:46",
			BirthDate:    "2005-08-26",
		}

		err := json.NewEncoder(w).Encode(respBody)
		assert.NoError(t, err)
	}

	var (
		ctx = context.Background()
		srv = server.New(t, handler)
		c   = client.New(
			client.WithBaseURL(srv.URL),
		)
		u = User{
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
	assert.Equal(t, want, userData)
}

func TestUser_Macros(t *testing.T) {
	t.Parallel()

	startDate, err := time.Parse(layoutISO, "2025-04-12")
	assert.NoError(t, err)

	endDate, err := time.Parse(layoutISO, "2025-04-13")
	assert.NoError(t, err)

	want := intake.MacrosRange{
		{
			Date:    startDate,
			Energy:  1288.68,
			Carb:    85.37,
			Fat:     62.17,
			Protein: 94.38,
		},
		{
			Date:    endDate,
			Energy:  1768.78,
			Carb:    156.53,
			Fat:     38.76,
			Protein: 182.95,
		},
	}

	handler := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		assert.Equal(t, r.URL.Path, macrosIntakesEndpoint)

		var (
			q              = r.URL.Query()
			queryStartDate = q.Get("start")
			queryEndDate   = q.Get("end")
		)

		respBody := GetMacroIntakeDTO{
			{
				Date:    queryStartDate,
				Energy:  1288.68,
				Carb:    85.37,
				Fat:     62.17,
				Protein: 94.38,
			},
			{
				Date:    queryEndDate,
				Energy:  1768.78,
				Carb:    156.53,
				Fat:     38.76,
				Protein: 182.95,
			},
		}

		assert.Write(t, w, respBody)
	}

	var (
		ctx = context.Background()
		srv = server.New(t, handler)
		c   = client.New(
			client.WithBaseURL(srv.URL),
		)
		u = User{
			client: c,
			token: &Token{
				expiresAt: times.Future(),
				access:    uuid.NewString(),
				refresh:   uuid.NewString(),
			},
		}
	)

	macroRange := date.Range{
		Start: startDate,
		End:   endDate,
	}
	rm, err := u.Macros(ctx, macroRange)
	assert.NoError(t, err)
	assert.DeepEqual(t, rm, want)
}

func TestUser_Intake(t *testing.T) {
	t.Parallel()

	startDate, err := time.Parse(layoutISO, "2025-04-12")
	assert.NoError(t, err)

	endDate, err := time.Parse(layoutISO, "2025-04-13")
	assert.NoError(t, err)

	type args struct {
		ctx       context.Context
		kind      intake.Kind
		dateRange date.Range
	}

	var (
		ctx        = context.Background()
		testBlocks = []struct {
			name          string
			args          args
			wantErr       bool
			want          intake.SingleRange
			serverHandler server.Handler
		}{
			{
				name: "valid args",
				args: args{
					ctx:  ctx,
					kind: intake.Sugar,
					dateRange: date.Range{
						Start: startDate, End: endDate,
					},
				},
				wantErr: false,
				want: intake.SingleRange{
					{Kind: intake.Sugar, Date: startDate, Value: 89.83},
					{Kind: intake.Sugar, Date: endDate, Value: 50.38},
				},
				serverHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.URL.Path, singleIntakesEndpoint)
					assert.Equal(t, r.Method, http.MethodGet)

					var (
						q               = r.URL.Query()
						queryStartDate  = q.Get("start")
						queryEndDate    = q.Get("end")
						queryIntakeKind = q.Get("nutrient")
					)

					if queryIntakeKind != intake.Sugar.ID() {
						t.Fatalf("want sugar intake id, got %q", queryIntakeKind)
					}

					respBody := GetSingleIntakeDTO{
						queryStartDate: 89.83,
						queryEndDate:   50.38,
					}
					assert.Write(t, w, respBody)
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			var (
				srv = server.New(t, tb.serverHandler)
				c   = client.New(client.WithBaseURL(srv.URL))
				u   = User{
					client: c,
					token: &Token{
						expiresAt: times.Future(),
						access:    uuid.NewString(),
						refresh:   uuid.NewString(),
					},
				}
			)

			rm, err := u.Intake(
				tb.args.ctx,
				tb.args.kind,
				tb.args.dateRange,
			)
			assert.NoError(t, err)
			assert.DeepEqual(t, rm, tb.want)
		})
	}
}

func TestUser_AddFood(t *testing.T) {
	t.Parallel()

	defaultHandler := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, addFoodEndpoint)
		assert.Equal(t, r.Method, http.MethodPost)

		headers := client.Payload{
			`Cache-Control`: ` no-cache, private`,
			`Connection`:    ` keep-alive`,
			`Date`:          ` Wed, 14 May 2025 19:03:19 GMT`,
			`Server`:        ` nginx`,
		}
		responseHeaders := w.Header()
		headers.Set(responseHeaders)

		w.WriteHeader(http.StatusNoContent)
	}

	var (
		ctx       = context.Background()
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
		testBlocks = []struct {
			name         string
			wantErr      bool
			food         food.Food
			serverHandle server.Handler
		}{
			{name: "valid food", food: validFood},
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
				name: "server respond (StatusBadRequest)", wantErr: true, food: validFood,
				serverHandle: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				},
			},
			{
				name: "server respond (StatusConflict)", wantErr: true, food: validFood,
				serverHandle: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusConflict)
				},
			},
		}
	)

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			t.Parallel()

			handler := defaultHandler
			if tb.serverHandle != nil {
				handler = tb.serverHandle
			}

			var (
				srv = server.New(t, handler)
				c   = client.New(client.WithBaseURL(srv.URL))
				u   = &User{
					client: c,
					token: &Token{
						expiresAt: times.Future(),
						access:    uuid.NewString(),
						refresh:   uuid.NewString(),
					},
				}
			)

			err := u.AddFood(
				ctx,
				tb.food,
				visibility.PrivateFood,
			)
			if (err != nil) != tb.wantErr {
				t.Fatalf("want (%v), got err: %v", tb.wantErr, err)
			}
		})
	}
}
