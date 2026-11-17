package utils

import "strconv"

// prices: []string, []float64, []int
// 将输入的有效参数 prices 转换计算集合
func qpInput(prices any) (deltas []string) {
	// 解析计算集 deltas elements array
	switch arrays := prices.(type) {
	case string:
		deltas = StrSplit(arrays)
	case []string:
		deltas = arrays
	case []float64:
		for _, v := range arrays {
			deltas = append(deltas, strconv.FormatFloat(v, 'f', -1, 64))
		}
	case []int:
		for _, v := range arrays {
			deltas = append(deltas, strconv.Itoa(v))
		}
	case []int64:
		for _, v := range arrays {
			deltas = append(deltas, strconv.Itoa(int(v)))
		}
	default:
		// err = errors.ErrUnsupported
	}
	return
}
