package application

import (
	"context"

	"github.com/controlado/go-yazio/internal/domain"
)

type User interface {
	IsExpired() bool
	Data(context.Context) (domain.User, error)
	Macros(context.Context, domain.DateRange) (domain.MacrosRange, error)
	Intake(context.Context, domain.IntakeKind, domain.DateRange) (domain.SingleRange, error)
}
