package utils

// 切片倒叙 slice array
func SliceOrderDesc(slice interface{}) (array interface{}) {

	switch arr := slice.(type) {
	case []string:
		var memberStr []string
		for _, v := range arr {
			memberStr = append([]string{v}, memberStr...)
		}
		array = memberStr
	case []int:
		var memberInt []int
		for _, v := range arr {
			memberInt = append([]int{v}, memberInt...)
		}
		array = memberInt
	default:
	}

	return
}
