package mymath

import "math"

func Round(x float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(x*factor) / factor
}
