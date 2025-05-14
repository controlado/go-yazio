package yazio

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/controlado/go-yazio/internal/domain"
	"github.com/controlado/go-yazio/internal/infra/client"
)

type User struct {
	client       *client.Client
	expiresAt    time.Time
	accessToken  string
	refreshToken string
}

// Data implements application.User.
func (u *User) Data(ctx context.Context) (d domain.User, err error) {
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
		return d, fmt.Errorf("%s: %w", ErrDecodingToInternalDTO, err)
	}

	return dto.toUserData()
}

// Intake implements application.User.
func (u *User) Intake(ctx context.Context, k domain.IntakeKind, r domain.DateRange) (domain.SingleRange, error) {
	var (
		dto GetSingleIntakeDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: intakeEndpoint,
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
		return nil, fmt.Errorf("%s: %w", ErrDecodingToInternalDTO, err)
	}

	return dto.toRangeSingle()
}

// Macros implements application.User.
func (u *User) Macros(ctx context.Context, r domain.DateRange) (domain.MacrosRange, error) {
	var (
		dto GetMacroIntakeDTO
		req = client.Request{
			Method:   http.MethodGet,
			Endpoint: macrosEndpoint,
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
