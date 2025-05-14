package domain

import (
	"fmt"
	"strings"
	"time"
)

const ( // required
	Energy      IntakeKind = "energy.energy"
	Fat         IntakeKind = "nutrient.fat"
	Saturated   IntakeKind = "nutrient.saturated"
	TransFat    IntakeKind = "nutrient.transfat"
	Cholesterol IntakeKind = "nutrient.cholesterol"
	Sodium      IntakeKind = "nutrient.sodium"
	Carb        IntakeKind = "nutrient.carb"
	Fiber       IntakeKind = "nutrient.dietaryfiber"
	Sugar       IntakeKind = "nutrient.sugar"
	AddedSugar  IntakeKind = "nutrient.addedsugar"
	Protein     IntakeKind = "nutrient.protein"
	VitaminD    IntakeKind = "vitamin.d"
	Calcium     IntakeKind = "mineral.calcium"
	Iron        IntakeKind = "mineral.iron"
	Potassium   IntakeKind = "mineral.potassium"
)

const ( // optional
	Monounsaturated IntakeKind = "nutrient.monounsaturated"
	Polyunsaturated IntakeKind = "nutrient.polyunsaturated"

	Water   IntakeKind = "nutrient.water"
	Alcohol IntakeKind = "nutrient.alcohol"
	Salt    IntakeKind = "nutrient.salt"

	VitaminA   IntakeKind = "vitamin.a"
	VitaminB1  IntakeKind = "vitamin.b1"
	VitaminB2  IntakeKind = "vitamin.b2"
	VitaminB3  IntakeKind = "vitamin.b3"
	VitaminB5  IntakeKind = "vitamin.b5"
	VitaminB6  IntakeKind = "vitamin.b6"
	VitaminB7  IntakeKind = "vitamin.b7"
	VitaminB11 IntakeKind = "vitamin.b11"
	VitaminB12 IntakeKind = "vitamin.b12"
	VitaminC   IntakeKind = "vitamin.c"
	VitaminE   IntakeKind = "vitamin.e"
	VitaminK   IntakeKind = "vitamin.k"

	MineralArsenic    IntakeKind = "mineral.arsenic"
	MineralBoron      IntakeKind = "mineral.boron"
	MineralBiotin     IntakeKind = "mineral.biotin"
	MineralCholine    IntakeKind = "mineral.choline"
	MineralChlorine   IntakeKind = "mineral.chlorine"
	MineralChrome     IntakeKind = "mineral.chrome"
	MineralCobalt     IntakeKind = "mineral.cobalt"
	MineralCopper     IntakeKind = "mineral.copper"
	MineralFluoride   IntakeKind = "mineral.fluoride"
	MineralIodine     IntakeKind = "mineral.iodine"
	MineralFluorine   IntakeKind = "mineral.fluorine"
	MineralManganese  IntakeKind = "mineral.manganese"
	MineralMagnesium  IntakeKind = "mineral.magnesium"
	MineralMolybdenum IntakeKind = "mineral.molybdenum"
	MineralPhosphorus IntakeKind = "mineral.phosphorus"
	MineralSelenium   IntakeKind = "mineral.selenium"
	MineralRubidium   IntakeKind = "mineral.rubidium"
	MineralSilicon    IntakeKind = "mineral.silicon"
	MineralTin        IntakeKind = "mineral.tin"
	MineralSulfur     IntakeKind = "mineral.sulfur"
	MineralZinc       IntakeKind = "mineral.zinc"
	MineralVanadium   IntakeKind = "mineral.vanadium"
)

type DateRange struct {
	Start time.Time
	End   time.Time
}

func (dr *DateRange) String() string {
	humanLayout := "2 January 2006"
	return fmt.Sprintf("%s - %s",
		dr.Start.Format(humanLayout),
		dr.End.Format(humanLayout),
	)
}

type IntakeKind string

func (ik IntakeKind) String() string {
	return string(ik)
}

type (
	SingleIntake struct {
		Date  time.Time
		Value float64
	}
	MacrosIntake struct {
		Date    time.Time
		Energy  float64
		Carb    float64
		Fat     float64
		Protein float64
	}
)

type SingleRange []SingleIntake

type SingleAverage struct {
	DaysLength int
	Average    float64
}

func (sa SingleAverage) String() string {
	return fmt.Sprintf("%d days: %.2f", sa.DaysLength, sa.Average)
}

func (sr SingleRange) Average() SingleAverage {
	var (
		totalValues float64
		rangeLength = len(sr)
	)

	if rangeLength == 0 {
		return SingleAverage{}
	}

	for _, intake := range sr {
		totalValues += intake.Value
	}

	return SingleAverage{
		DaysLength: rangeLength,
		Average:    totalValues / float64(rangeLength),
	}
}

type MacrosRange []MacrosIntake

type MacrosAverage struct {
	DaysLength int
	Energy     float64
	Carb       float64
	Fat        float64
	Protein    float64
}

func (ma MacrosAverage) String() string {
	stringParts := []string{
		fmt.Sprintf("Average %d days", ma.DaysLength),
		fmt.Sprintf("Kcal: %.2f", ma.Energy),
		fmt.Sprintf("Carb: %.2f", ma.Carb),
		fmt.Sprintf("Fat: %.2f", ma.Fat),
		fmt.Sprintf("Protein: %.2f", ma.Protein),
	}
	return strings.Join(stringParts, "\n")
}

func (mr MacrosRange) Average() MacrosAverage {
	var (
		totalEnergy,
		totalCarb,
		totalFat,
		totalProtein float64
		rangeLength = len(mr)
	)

	if rangeLength == 0 {
		return MacrosAverage{}
	}

	for _, intake := range mr {
		totalEnergy += intake.Energy
		totalCarb += intake.Carb
		totalFat += intake.Fat
		totalProtein += intake.Protein
	}

	return MacrosAverage{
		DaysLength: rangeLength,
		Energy:     totalEnergy / float64(rangeLength),
		Carb:       totalCarb / float64(rangeLength),
		Fat:        totalFat / float64(rangeLength),
		Protein:    totalProtein / float64(rangeLength),
	}
}
