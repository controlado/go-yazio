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
	stringParts := []string{
		fmt.Sprintf("Average %d days", ma.DaysLength),
		fmt.Sprintf("Kcal: %.2f", ma.Energy),
		fmt.Sprintf("Carb: %.2f", ma.Carb),
		fmt.Sprintf("Fat: %.2f", ma.Fat),
		fmt.Sprintf("Protein: %.2f", ma.Protein),
	}
	return strings.Join(stringParts, "\n")
}
