package ledger

import (
	"fmt"
	"math"
	"strconv"
)

func RoundOld(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func Round(val float64, roundOn float64, places int) float64 {
	result := fmt.Sprintf("%.2f", val)
	float, _ := strconv.ParseFloat(result, 64)
	return float
}
