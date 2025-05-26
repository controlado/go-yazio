package application

import (
	"context"

	"github.com/controlado/go-yazio/pkg/domain/date"
	"github.com/controlado/go-yazio/pkg/domain/food"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/meal"
	"github.com/controlado/go-yazio/pkg/domain/user"
	"github.com/controlado/go-yazio/pkg/visibility"
)

type User interface {
	Token() Token
	Data(context.Context) (user.Data, error)
	AddFood(context.Context, food.Food, visibility.Food) error
	EntryFood(context.Context, meal.Time, food.ID, food.Serving) error
	Macros(context.Context, date.Range) (intake.MacrosRange, error)
	Intake(context.Context, intake.Kind, date.Range) (intake.SingleRange, error)
}
