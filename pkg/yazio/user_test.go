package yazio

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/controlado/go-yazio/internal/domain"
	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/internal/testutil/server"
	"github.com/controlado/go-yazio/pkg/client"
	"github.com/google/uuid"
)

func TestUser_Data(t *testing.T) {
	t.Parallel()

	parsedFakeID, err := uuid.Parse("21a7e919-b3f2-4abc-a6b8-83dddfe311a6")
	assert.NoError(t, err)

	want := domain.User{
		ID:        parsedFakeID,
		Token:     "c000a7769600a98abae7cefe56174e48240ee297e06be3052cc3e743f12bcfd5",
		FirstName: "João Brito",
		LastName:  "da Silva",
		IconURL:   "https://images.yazio-cdn.com/process/plain/app/profile/user/2025/d297247d-51d4-4e04-9e87-c99fdf693585.jpg",
		Email: domain.Email{
			Value:       "joaodasilva@gmail.com",
			IsConfirmed: true,
		},
		Registration: time.Date(2023, 02, 06, 21, 22, 46, 0, time.UTC),
		Birth:        time.Date(2005, 8, 26, 0, 0, 0, 0, time.UTC),
	}

	handler := func(t *testing.T, w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		assert.Equal(t, r.URL.Path, userDataEndpoint)

		respBody := GetUserDataDTO{
			ID:           parsedFakeID.String(),
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
			client:       c,
			expiresAt:    time.Now().Add(time.Hour),
			accessToken:  "302af606a79142cb2ab862bf9488cfd4",
			refreshToken: "302af606a79142cb2ab862bf9488cfd4",
		}
	)

	userData, err := u.Data(ctx)
	assert.NoError(t, err)
	assert.Equal(t, want, userData)
}

func TestUser_IsExpired(t *testing.T) {
	t.Parallel()

	var (
		c = client.New(
			client.WithBaseURL(DefaultBaseURL),
		)
		u = User{
			client:       c,
			expiresAt:    time.Now().Add(-time.Hour),
			accessToken:  "302af606a79142cb2ab862bf9488cfd4",
			refreshToken: "302af606a79142cb2ab862bf9488cfd4",
		}
	)

	want := true
	got := u.IsExpired()
	assert.Equal(t, got, want)
}

func TestUser_Macros(t *testing.T) {
	t.Parallel()

	startDate, err := time.Parse(layoutISO, "2025-04-12")
	assert.NoError(t, err)

	endDate, err := time.Parse(layoutISO, "2025-04-13")
	assert.NoError(t, err)

	want := domain.MacrosRange{
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
		assert.Equal(t, r.URL.Path, macrosEndpoint)

		q := r.URL.Query()
		queryStartDate := q.Get("start")
		queryEndDate := q.Get("end")

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
			client:       c,
			expiresAt:    time.Now().Add(time.Hour),
			accessToken:  "302af606a79142cb2ab862bf9488cfd4",
			refreshToken: "302af606a79142cb2ab862bf9488cfd4",
		}
	)

	macroRange := domain.DateRange{
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
		kind      domain.IntakeKind
		dateRange domain.DateRange
	}

	var (
		ctx        = context.Background()
		testBlocks = []struct {
			name          string
			args          args
			wantErr       bool
			want          domain.SingleRange
			serverHandler server.Handler
		}{
			{
				name: "valid args",
				args: args{
					ctx:       ctx,
					kind:      domain.Sugar,
					dateRange: domain.DateRange{Start: startDate, End: endDate},
				},
				wantErr: false,
				want: domain.SingleRange{
					{
						Date:  startDate,
						Value: 89.83,
					},
					{
						Date:  endDate,
						Value: 50.38,
					},
				},
				serverHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, r.URL.Path, intakeEndpoint)
					assert.Equal(t, r.Method, http.MethodGet)

					var (
						q               = r.URL.Query()
						queryStartDate  = q.Get("start")
						queryEndDate    = q.Get("end")
						queryIntakeKind = q.Get("nutrient")
					)

					switch domain.IntakeKind(queryIntakeKind) {
					case domain.Salt:
					case domain.Sugar:
					case domain.Fiber:
					case domain.Water:
					default:
						t.Fatalf("unexpected kind %q", queryIntakeKind)
					}

					respBody := GetSingleIntakeDTO{
						queryStartDate: 89.83,
						queryEndDate:   50.38,
					}
					err := json.NewEncoder(w).Encode(respBody)
					assert.NoError(t, err)
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
					client:       c,
					expiresAt:    time.Now().Add(time.Hour),
					accessToken:  "302af606a79142cb2ab862bf9488cfd4",
					refreshToken: "302af606a79142cb2ab862bf9488cfd4",
				}
			)

			rm, err := u.Intake(tb.args.ctx, tb.args.kind, tb.args.dateRange)
			assert.NoError(t, err)
			assert.DeepEqual(t, rm, tb.want)
		})
	}
}
