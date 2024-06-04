package mymath

import "math"

func Max[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](nums ...T) T {
	var maxNum T
	for i, num := range nums {
		if i == 0 {
			maxNum = num
		} else {
			if num > maxNum {
				maxNum = num
			}
		}

	}
	return maxNum
}

func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](nums ...T) T {
	var minNum T
	for i, num := range nums {
		if i == 0 {
			minNum = num
		} else {
			if num < minNum {
				minNum = num
			}
		}
	}
	return minNum
}

func Sum[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](nums ...T) T {
	var sumNum T
	for _, num := range nums {
		sumNum += num
	}
	return sumNum
}

func Float64ByPrecision(num float64, precision int, roundDown bool) float64 {
	multiplier := math.Pow10(precision)
	truncated := num * multiplier
	if roundDown {
		truncated = math.Floor(truncated)
	} else {
		truncated = math.Round(truncated)
	}
	return truncated / multiplier
}
