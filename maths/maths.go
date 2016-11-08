package maths

import "math/big"

func Round(val float64, places uint) float64 {
	tmp := big.NewFloat(val)
	tmp.SetPrec(places)
	result, _ := tmp.Float64()
	return result
}
