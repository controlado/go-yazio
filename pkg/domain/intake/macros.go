package intake

import (
	"fmt"
	"strings"
	"time"
)

type Macros struct {
	Date    time.Time
	Energy  float64
	Carb    float64
	Fat     float64
	Protein float64
}

type MacrosRange []Macros

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

type MacrosAverage struct {
	DaysLength int
	Energy     float64
	Carb       float64
	Fat        float64
	Protein    float64
}

func (ma MacrosAverage) String() string {
	if ma.DaysLength < 1 {
		return "Empty macro data to calculate the average"
	}

	stringParts := []string{
		fmt.Sprintf("Average (%d days)", ma.DaysLength),
		fmt.Sprintf("Energy: %.1f%s", ma.Energy, Energy.baseUnit),
		fmt.Sprintf("Carb: %.1f%s", ma.Carb, Carb.baseUnit),
		fmt.Sprintf("Fat: %.1f%s", ma.Fat, Fat.baseUnit),
		fmt.Sprintf("Protein: %.1f%s", ma.Protein, Protein.baseUnit),
	}
	return strings.Join(stringParts, "\n")
}
