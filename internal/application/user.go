package application

import (
	"context"

	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
)

type User interface {
	IsExpired() bool
	Data(context.Context) (user.User, error)
	AddFood(context.Context, food.Food, visibility.Food) error
	Macros(context.Context, date.Range) (intake.MacrosRange, error)
	Intake(context.Context, intake.Kind, date.Range) (intake.SingleRange, error)
}
