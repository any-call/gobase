package mymath

import (
	"fmt"
	"math/rand"
)

func HBMath(totalMoney float64, minAmount float64, num int) (result []float64, err error) {
	if minAmount*float64(num) > totalMoney {
		return nil, fmt.Errorf("total money is not enough")
	}

	// 每个红包先填充最小值
	var remainAmount float64 = totalMoney
	var remainNumber int = num
	result = make([]float64, num)
	for i, _ := range result {
		result[i] = minAmount
		remainAmount -= minAmount
	}

	//2倍剩余均值金额
	var averageMax float64
	// 单个红包金额
	var amount float64
	for index, _ := range result {
		if remainAmount <= 0 || remainNumber <= 0 {
			break
		}

		// 剩余红包的2倍平均值
		averageMax = 2 * remainAmount / float64(remainNumber)

		// 最后一个得到全部剩余
		if remainNumber == 1 {
			amount = remainAmount
		} else {
			amount = Round(rand.Float64()*averageMax, 2)
		}

		result[index] += amount
		remainAmount -= amount
		remainNumber--
	}

	return result, nil
}
