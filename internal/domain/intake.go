package domain

import (
	"fmt"
	"strings"
	"time"
)

const (
	Salt  IntakeKind = "nutrient.salt"
	Sugar IntakeKind = "nutrient.sugar"
	Fiber IntakeKind = "nutrient.dietaryfiber"
	Water IntakeKind = "water"
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
