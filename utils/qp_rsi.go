package utils

import (
	"math"
	"strconv"
)

// QP(Quantizer Parameter)量化参数

// AvgUp = sum(+ΔP)/Period
//
// AvgDown = sum(|-ΔP|)/Period
//
// rs = AvgUp/AvgDown
//
// rsi = 100 - (100/(1+rs))
//
// calculate RSI 计算相对强弱指标(period默认14)
func CalcRsiOnce(prices any, _period ...int) (rsi, avgGain, avgLoss string) {
	period := 14 // 默认计算窗口
	if len(_period) > 0 {
		period = _period[0]
	}

	// 解析计算集 deltas elements array
	var deltas = qpInput(prices)
	if len(deltas) < period {
		// err = errors.New("element length is insufficient")
		return
	}

	// 计算每日收盘价变化 deltas
	// 初始化 gain 和 loss 数组，用于计算平均涨跌幅
	var gainSum, lossSum string
	for i := 1; i < period; i++ {
		change := CalcFloatStr(deltas[i], deltas[i-1], '-')
		if change[0] == '-' {
			lossSum = CalcFloatStr(lossSum, change[1:])
		} else {
			gainSum = CalcFloatStr(gainSum, change)
		}
		// fmt.Printf("%v: %v (%v) - %v %v\n", i, deltas[i], change, gainSum, lossSum)
	}

	// 计算平均涨跌幅
	avgGain = CalcFloatStr(gainSum, period, '/')
	avgLoss = CalcFloatStr(lossSum, period, '/')

	// 初始化第一个period的平均涨跌幅
	rsi = strconv.FormatFloat(calcRsi(avgGain, avgLoss), 'f', -1, 64)
	// fmt.Println("one: ", rsi, avgGain, avgLoss)

	// 计算后续period的平均涨跌幅
	// fmt.Printf("len(deltas): %v\n", len(deltas))
	// fmt.Printf("period: %v %v \n", period, len(deltas) > period)
	if len(deltas) > period {
		for i := period; i < len(deltas); i++ {
			// fmt.Printf(" %d -", i)
			rsi, avgGain, avgLoss = CalcRSI(deltas[i], deltas[i-1], avgGain, avgLoss, period)
			// fmt.Println("rsi: ", rsi, avgGain, avgLoss)
		}
	}

	// 计算RSI
	// for i := period - 1; i < len(prices); i++ {
	// 	if avgLoss[i] == 0 {
	// 		rsi[i] = 100
	// 	} else {
	// 		rs := avgGain[i] / avgLoss[i]
	// 		rsi[i] = 100 - (100 / (1 + rs))
	// 	}
	// }
	return
}

// AvgUp = (AvgGain * (Period-1) + (+ΔP)) / Period
//
// AvgDown = (AvgLoss * (Period-1) + (|-ΔP|)) / Period
//
// rs = AvgUp/AvgDown
//
// rsi = 100 - (100/(1+rs))
//
// 计算增量rsi: 当前价格'price', 增量参数'pre-'
func CalcRSI(price, prePrice, preAvgGain, preAvgLoss string, _period ...int) (rsi, avgGain, avgLoss string) {
	period := 14
	if len(_period) > 0 {
		period = _period[0]
	}

	change := CalcFloat(price, prePrice, '-')
	var gain, loss float64
	if change > 0 {
		gain = change
	} else {
		loss = math.Abs(change)
	}
	avgGain = CalcFloatStr(CalcFloat(CalcFloat(preAvgGain, period-1, '*'), gain), period, '/')
	avgLoss = CalcFloatStr(CalcFloat(CalcFloat(preAvgLoss, period-1, '*'), loss), period, '/')
	rsi = strconv.FormatFloat(calcRsi(avgGain, avgLoss), 'f', -1, 64)
	return
}

// 通过平均涨跌幅计算当前rsi值
func calcRsi(avgGain, avgLoss string) float64 {
	// 校正 直降, 直涨, 不升降 问题
	var rs, rsi float64
	if IsZero(avgGain) && IsZero(avgLoss) { // 未涨跌 -> rsi=50, ri=1可计算
		rsi = 50
		//elseif IsZero(avgGain)// 直降 -> rsi=0 , 可以直接计算
	} else if IsZero(avgLoss) { // 直涨 -> rsi=100
		// rs = 999999999
		rsi = 100
		// return 100
	} else {
		rs = CalcFloat(avgGain, avgLoss, '/')
		rsi = 100 - (100 / (1 + rs))
		// fmt.Printf("rsi(%v): %v = %v/%v\n", rsi, rs, avgGain, avgLoss)
	}

	return rsi
}
