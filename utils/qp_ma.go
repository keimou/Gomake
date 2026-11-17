package utils

import (
	"strconv"
)

//	*公式: SMA = (P1 + P2 + ... + Pn) / n
//	prices: string, []string, []float64, []int64, []int
//
// 如何计算集足够, 将计算出最后一个参数对应的sma
func CalcSMA(prices any, period int) (sma string) {
	if !(period > 3) {
		// err = errors.New("period is insufficient")
		return
	}

	// 解析计算集
	deltas := qpInput(prices)
	if len(deltas) < period {
		// err = errors.New("element length is insufficient")
		return "NaN"
	}
	// 计算初始值（SMA）
	var sum float64
	for i := len(deltas) - period; i < len(deltas); i++ {
		// fmt.Printf("%d: %v\n", i, deltas[i])
		sum = CalcFloat(sum, deltas[i])
	}
	// sma := sum / float64(period)
	// sma = CalcFloatStr(sum, period, '/')
	sma = strconv.FormatFloat(sum/float64(period), 'f', -1, 64)
	// fmt.Printf("sma: %v\n", sma)
	return
}

//	α = 2/(Period+1);
//	EMA[t] = Price[t] x α + EMA[t-1] x (1−α)
//	First EMA is SMA(prices, period)
//
// 递归 增量方式计算 EMA:
func CalcEMA(price, preEma string, period int) (ema string) {
	if !(period > 3) {
		// err = errors.New("period is insufficient")
		return
	}

	// 计算平滑系数 α
	alpha := 2.0 / (float64(period) + 1.0)

	ema = CalcFloatStr(CalcFloat(price, alpha, '*'), CalcFloat(preEma, 1.0-alpha, '*'))

	return
}

// * EMA 递归公式:
// α = 2/(N+1)
// EMA(t) = 当前Price(t)×α + 前一日EMA(t-1)×(1−α) 	->	EMA[t] = C*2/(N+1) + EMA[t-1]*(N-1)/(N+1)
// 第一个EMA(1)为SMA(1)
// EMA = SMA
