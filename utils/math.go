package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// 重新指定无穷小数位有限小数
func Normalize(x float64) float64 {
	const SmallestNormal = 2.2250738585072014e-308 // 2**-1022
	if math.Abs(x) < SmallestNormal {
		return x * (1 << 52) // 较大的范围: x * 2^52
	}
	return x
}

// 使用线性回归计算斜率
func LinearRegressionSlope(values []float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}

	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < n; i++ {
		sumX += float64(i)
		sumY += values[i]
		sumXY += float64(i) * values[i]
		sumX2 += float64(i) * float64(i)
	}

	numerator := float64(n)*sumXY - sumX*sumY
	denominator := float64(n)*sumX2 - sumX*sumX

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// 两点 计算斜率 模拟回归
func TwoPointSlope(y1, y2 float64, xs ...float64) float64 {
	// sumY := (y1 + y2) * 5 / 2
	// sumXY := y1*4 + y2*6
	// sumX := 1 + 2 + 3 + 4
	// sumX2 = 1 + 2*2 + 3*3 + 4*4
	// numerator := 5*(y1*4+y2*6) - 10*((y1+y2)*5/2)
	// denominator := float64(5)*sumX2 - sumX*sumX

	// x := 60.0
	// if len(xs) > 0 {
	// 	x = xs[0] / 1000
	// }
	// return (y2 - y1) / x
	return (5*(y1*4+y2*6) - 10*((y1+y2)*5/2)) / 50
}

// 1717171200000
// - 1717171260000 1m
// - 1717171380000 3m
// - 1717171500000 5m
// - 1717174800000 1h
// - 1717257600000 1d
// func CalculateSlope(oldy, newy float64, oldx, newx int64) float64 {
// 	if oldx == newx {
// 		return 0
// 	}
// 	gap := int64(1000) // m-60000, s-1000
// 	oldx /= gap
// 	newx /= gap
// 	return (newy - oldy) / float64(newx-oldx)
// }

// *
// *
// Float小数位: 规范取小数点后 n 位有效, (styles - 0: 四舍五入, 1: 向上取整, -1: 向下取整)
func FloatDecimalPlaces(value float64, n int, styles ...int) float64 {
	// math.Ceil:	向上取整, 返回不小于输入值的最小整数.
	// math.Round:	四舍五入取整.
	// math.Floor:	向下取整, 返回不大于输入值的最大整数.
	const (
		round = 0
		ceil  = 1
		floor = -1
	)

	var dividend float64

	factor := math.Pow(10, float64(n)) // 因子 因素

	style := round // 取有效位数风格
	if len(styles) > 0 {
		style = styles[0]
	}
	switch style {
	case ceil:
		dividend = math.Ceil(value * factor)
	case floor:
		dividend = math.Floor(value * factor)
	default:
		dividend = math.Round(value * factor)
	}
	return dividend / factor
}

// Float小数位String: 规范取小数点后 n 位有效, (styles - 0: 四舍五入, 1: 向上取整, -1: 向下取整)
func FloatDecimalPlaceStr(value float64, n int, styles ...int) string {
	const (
		round = 0  // math.Round:	四舍五入取整.
		ceil  = 1  // math.Ceil:	向上取整, 返回不小于输入值的最小整数.
		floor = -1 // math.Floor:	向下取整, 返回不大于输入值的最大整数.
	)
	style := round // 取有效位数风格
	if len(styles) > 0 {
		style = styles[0]
	}
	num := FloatDecimalPlaces(value, n, style)

	return strconv.FormatFloat(num, 'f', n, 64)

}

// 计算小数位后 第几位 为有效数
func GetFloatDecimalPlace(value float64) (digit int) {
	compared := math.Abs(value)
	if compared > 1 {
		digit = 0
	} else if compared > 0.1 {
		digit = 1
	} else if compared > 0.01 {
		digit = 2
	} else if compared > 0.001 {
		digit = 3
	} else if compared > 0.0001 {
		digit = 4
	} else if compared > 0.00001 {
		digit = 5
	} else {
		digit = 6
	}
	return
}

// 计算小数位后 有几位数
func GetFloatDecimalPlaces(values interface{}) (decimalPlaces int) {
	var str string
	switch value := values.(type) {
	case string:
		str = convertScientificNotation(value)
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64) // 将浮点数转换为字符串
	default:
	}

	// 查找小数点的位置
	if idx := strings.IndexByte(str, '.'); idx != -1 {
		// 计算小数点后面的字符数
		decimalPlaces = len(str) - idx - 1
		// fmt.Printf("小数点后有 %d 位数\n", decimalPlaces)
	} else {
		decimalPlaces = 0
		// fmt.Println("没有小数部分")
	}

	return
}

// 取 value2 的值, 取小于 value1 小数位的位数
func SetBothFloat(value1, value2 string) string {
	var str string

	value := convertScientificNotation(value2)

	n1 := GetFloatDecimalPlaces(value1)
	n2 := GetFloatDecimalPlaces(value2)
	if n1 < n2 {
		nx := n2 - n1
		str = value[:len(value)-nx]
	} else if n1 > n2 {
		nx := n1 - n2
		str = value + strings.Repeat("0", nx)
	} else {
		str = value
	}

	return str
}

// // 定义匹配科学计数法的正则表达式
// scientificNotationPattern := `^[+-]?\d+(\.\d+)?([eE][+-]?\d+)?$`
// re := regexp.MustCompile(scientificNotationPattern)
// re.MatchString(s)
// 判断字符串是否为科学计数法表示
func isScientificNotation(str string) bool {
	return strings.ContainsAny(str, "Ee")
}

// 如果是科学计数法表示的字符串则格式化为非科学计数法的字符串, 否则保留原字符串
func convertScientificNotation(strFloat string) string {
	// 如果字符串匹配科学计数法模式
	if isScientificNotation(strFloat) {
		// 将字符串转换为浮点数
		floatVal, err := strconv.ParseFloat(strFloat, 64)
		if err != nil {
			// 如果转换失败，返回原始字符串
			return strFloat
		}
		// 将浮点数格式化为非科学计数法的字符串
		return strconv.FormatFloat(floatVal, 'f', -1, 64)
	}

	// 如果不是科学计数法表示，返回原始字符串
	return strFloat
}

// *
// 浮点数计算精度处理
// *

// *
// 处理两浮点数字相加精度问题，字符结果 x: +(默认) - * /
func CalcFloatStr(a, b interface{}, x ...rune) (results string) {
	sign := xSign(x)
	if decimal, ok := calcFloat(a, b, sign); !ok {
		results = decimal.String()
	} else {
		results = "NaN"
	}
	return
}

// 处理两浮点数字相加精度问题，数字结果 x: +(默认) - * /
func CalcFloat(a, b interface{}, x ...rune) (results float64) {
	sign := xSign(x)

	if decimal, ok := calcFloat(a, b, sign); !ok {
		results, _ = decimal.Float64()
	} else {
		results = math.NaN()
	}
	return
}

// 浮点精度计算主体
func calcFloat(a, b interface{}, sign rune) (result decimal.Decimal, ok bool) {

	floatA := transFloat(a)
	floatB := transFloat(b)
	// // 将字符串转换为浮点数
	// floatA, _ := strconv.ParseFloat(a, 64)
	// floatB, _ := strconv.ParseFloat(b, 64)

	if math.IsNaN(floatA) || math.IsNaN(floatB) {
		fmt.Println(" * CalcFloat is NaN: ", a, string(sign), b)
		ok = true
		return
	}

	// // 进行加法操作
	// result := floatA + floatB
	// 0.1+0.2 0.1+0.7 0.2+0.4 0.2+0.7 0.3+0.6 1.2+1.9 2.2+9.9 3.2+9.9
	// 处理精度问题 (1590.9+49.2=1640.1000000000001) (0.1+0.2=0.30000000000000004)
	decimalA := decimal.NewFromFloat(floatA)
	decimalB := decimal.NewFromFloat(floatB)

	switch sign {
	case '+':
		// result = decimalA.Add(decimalB).String()
		result = decimalA.Add(decimalB)
	case '-':
		result = decimalA.Sub(decimalB)
	case '*':
		result = decimalA.Mul(decimalB)
	case '/':
		result = decimalA.Div(decimalB)
	case '%':
		result = decimalA.Mod(decimalB)
	default:
		// result = decimal.NewFromFloat(404) // 未定义
		ok = true
	}
	// fmt.Printf("%v %s %v = %v\n", decimalA, string(sign), decimalB, result)
	// 将结果转换回字符串
	// return strconv.FormatFloat(result, 'f', -1, 64)
	// decimalA.Add(decimalB).Float64()
	return
}

// 浮点数转换
func transFloat(value interface{}) (float float64) {
	switch v := value.(type) {
	case float64:
		float = v
	case float32:
		float = float64(v)
	case int:
		float = float64(v)
	case int64:
		float = float64(v)
	case int32:
		float = float64(v)
	case string:
		float, _ = strconv.ParseFloat(v, 64)
		// case uint:
		// 	float = float64(v)
	}
	return
}

// 计算符号识别
func xSign(x []rune) rune {
	sign := '+'
	if len(x) > 0 {
		sign = x[0]
		if sign == 'X' || sign == 'x' {
			sign = '*'
		}
		if sign != '+' && sign != '-' && sign != '*' && sign != '/' && sign != '%' {
			return ' '
		}
	}
	return sign
}
