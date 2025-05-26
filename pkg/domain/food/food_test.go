package food

import (
	"fmt"
	"testing"

	"github.com/controlado/go-yazio/internal/testutil/assert"
	"github.com/controlado/go-yazio/pkg/domain/intake"
	"github.com/controlado/go-yazio/pkg/domain/unit"
	"github.com/google/uuid"
)

func mockNewUUID() (undo func()) {
	originalFn := newUUID
	staticUUID := originalFn()
	newUUID = func() uuid.UUID { return staticUUID }
	return func() { newUUID = originalFn }
}

func TestNew(t *testing.T) {
	// TODO: refactor
	// t.Parallel() @ global mocking uuid.New

	undoMock := mockNewUUID()
	defer undoMock()

	type args struct {
		name string
		cat  Category
		nut  Nutrients
		opts []Option
	}

	testBlocks := []struct {
		name    string
		args    args
		want    Food
		wantErr bool
	}{
		{name: "invalid name", wantErr: true},
		{
			name: "using default unit/serving",
			args: args{
				name: "Banana",
				cat:  Miscellaneous,
				nut:  Nutrients{intake.AddedSugar: 0.1},
			},
			want: Food{
				ID:        newUUID(),
				Name:      "Banana",
				BaseUnit:  unit.Gram,
				Category:  Miscellaneous,
				Nutrients: Nutrients{intake.AddedSugar: 0.1},
				Servings:  []Serving{defaultServing},
			},
		},
		{
			name: "using options: defining unit/serving",
			args: args{
				name: "Liquid",
				cat:  Miscellaneous,
				nut:  Nutrients{intake.Water: 0.1},
				opts: []Option{
					WithBaseUnit(unit.Milliliter),
					WithNewServing(Pack, 50),
				},
			},
			want: Food{
				ID:        newUUID(),
				Name:      "Liquid",
				BaseUnit:  unit.Milliliter,
				Category:  Miscellaneous,
				Nutrients: Nutrients{intake.Water: 0.1},
				Servings: []Serving{
					{
						Kind:   Pack,
						Amount: 50,
					},
				},
			},
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			// TODO: refactor
			// t.Parallel() @ global mocking uuid.New

			got, err := New(
				tb.args.name,
				tb.args.cat,
				tb.args.nut,
				tb.args.opts...,
			)
			assert.WantErr(t, tb.wantErr, err)
			assert.DeepEqual(t, got, tb.want)
		})
	}
}

func TestFood_String(t *testing.T) {
	// TODO: refactor
	// t.Parallel() @ global mocking uuid.New

	undoMock := mockNewUUID()
	defer undoMock()

	testBlocks := []struct {
		name string
		f    Food
		want string
	}{
		{
			name: "valid food",
			f: Food{
				ID:       newUUID(),
				Name:     "Sadness",
				BaseUnit: unit.Microgram,
				Category: Meat,
				Nutrients: Nutrients{
					intake.Salt:    50,
					intake.Protein: 10,
					intake.Water:   20,
				},
			},
			want: fmt.Sprintf("Food(%q, %s, %s)", "Sadness", newUUID(), Meat),
		},
	}

	for _, tb := range testBlocks {
		t.Run(tb.name, func(t *testing.T) {
			// TODO: refactor
			// t.Parallel() @ global mocking uuid.New

			got := tb.f.String()
			assert.Equal(t, got, tb.want)
		})
	}
}
