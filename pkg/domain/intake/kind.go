package intake

import (
	"github.com/controlado/go-yazio/pkg/domain/unit"
)

type Kind struct {
	id       string
	baseUnit unit.Base
}

func (k Kind) ID() string {
	return k.id
}

func (k Kind) Unit() string {
	return k.baseUnit.String()
}

var (
	Energy            = Kind{"energy.energy", unit.Kilocalorie}
	Fat               = Kind{"nutrient.fat", unit.Gram}
	Saturated         = Kind{"nutrient.saturated", unit.Gram}
	TransFat          = Kind{"nutrient.transfat", unit.Gram}
	Cholesterol       = Kind{"nutrient.cholesterol", unit.Microgram}
	Sodium            = Kind{"nutrient.sodium", unit.Microgram}
	Carb              = Kind{"nutrient.carb", unit.Gram}
	Fiber             = Kind{"nutrient.dietaryfiber", unit.Gram}
	Sugar             = Kind{"nutrient.sugar", unit.Gram}
	AddedSugar        = Kind{"nutrient.addedsugar", unit.Gram}
	Protein           = Kind{"nutrient.protein", unit.Gram}
	Salt              = Kind{"nutrient.salt", unit.Gram}
	VitaminD          = Kind{"vitamin.d", unit.Microgram}
	Calcium           = Kind{"mineral.calcium", unit.Milligram}
	Iron              = Kind{"mineral.iron", unit.Milligram}
	Potassium         = Kind{"mineral.potassium", unit.Milligram}
	Monounsaturated   = Kind{"nutrient.monounsaturated", unit.Gram}
	Polyunsaturated   = Kind{"nutrient.polyunsaturated", unit.Gram}
	VitaminA          = Kind{"vitamin.a", unit.Microgram}
	VitaminB1         = Kind{"vitamin.b1", unit.Milligram}
	VitaminB2         = Kind{"vitamin.b2", unit.Milligram}
	VitaminB3         = Kind{"vitamin.b3", unit.Milligram}
	VitaminB5         = Kind{"vitamin.b5", unit.Milligram}
	VitaminB6         = Kind{"vitamin.b6", unit.Milligram}
	VitaminB7         = Kind{"vitamin.b7", unit.Microgram}
	VitaminB11        = Kind{"vitamin.b11", unit.Microgram}
	VitaminB12        = Kind{"vitamin.b12", unit.Microgram}
	VitaminC          = Kind{"vitamin.c", unit.Milligram}
	VitaminE          = Kind{"vitamin.e", unit.Milligram}
	VitaminK          = Kind{"vitamin.k", unit.Microgram}
	MineralArsenic    = Kind{"mineral.arsenic", unit.Microgram}
	MineralBoron      = Kind{"mineral.boron", unit.Milligram}
	MineralBiotin     = Kind{"mineral.biotin", unit.Microgram}
	MineralCholine    = Kind{"mineral.choline", unit.Milligram}
	MineralChlorine   = Kind{"mineral.chlorine", unit.Milligram}
	MineralChrome     = Kind{"mineral.chrome", unit.Milligram}
	MineralCobalt     = Kind{"mineral.cobalt", unit.Microgram}
	MineralCopper     = Kind{"mineral.copper", unit.Milligram}
	MineralFluoride   = Kind{"mineral.fluoride", unit.Milligram}
	MineralFluorine   = Kind{"mineral.fluorine", unit.Milligram}
	MineralIodine     = Kind{"mineral.iodine", unit.Microgram}
	MineralMagnesium  = Kind{"mineral.magnesium", unit.Milligram}
	MineralManganese  = Kind{"mineral.manganese", unit.Milligram}
	MineralMolybdenum = Kind{"mineral.molybdenum", unit.Microgram}
	MineralPhosphorus = Kind{"mineral.phosphorus", unit.Milligram}
	MineralRubidium   = Kind{"mineral.rubidium", unit.Microgram}
	MineralSelenium   = Kind{"mineral.selenium", unit.Microgram}
	MineralSilicon    = Kind{"mineral.silicon", unit.Milligram}
	MineralSulfur     = Kind{"mineral.sulfur", unit.Milligram}
	MineralTin        = Kind{"mineral.tin", unit.Milligram}
	MineralVanadium   = Kind{"mineral.vanadium", unit.Microgram}
	MineralZinc       = Kind{"mineral.zinc", unit.Milligram}
	Water             = Kind{"nutrient.water", unit.Milliliter}
	Alcohol           = Kind{"nutrient.alcohol", unit.Milliliter}
)
