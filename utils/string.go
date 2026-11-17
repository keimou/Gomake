package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// "a, b, c" -> ["a", "b", "c"]  // StrConvSlice()
// 字符串切割成字符串切片 (可设置同时多个间隔字符进行拆分,且忽略空字符串,除去前后多余空格)
func StrSplit(s string, separators ...rune) (result []string) {
	if len(separators) == 0 {
		separators = []rune{',', ';'} // ' '?
	}
	// 定义分隔符函数
	splitFunc := func(r rune) bool {
		for _, sep := range separators {
			if r == sep {
				return true
			}
		}
		return false
	}

	// parts := strings.Split(s, ",") // 按逗号切割字符串
	// 使用 strings.FieldsFunc 切割字符串
	parts := strings.FieldsFunc(s, splitFunc)

	// 去除每个部分的多余空格
	// result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" { // 忽略空字符串
			result = append(result, trimmed)
		}
	}
	return result
}

// Combine Unite Unify Merge Join Concatenate
// 将字符串切片组拼接成字符串, 默认间隔','
func StrCombine(array []string, separators ...rune) string {
	if len(array) == 0 {
		return ""
	}

	sign := ","
	if len(separators) > 0 {
		sign = string(separators[0])
	}

	// result := array[0] // v 1.0
	var sb strings.Builder // v 2.0 strings.Builder可以提高性能, 尤其是在处理大量数据时
	sb.WriteString(array[0])
	for _, a := range array[1:] {
		// // result += sep + a
		// sb.WriteString(sign + strings.TrimSpace(a))
		sb.WriteString(sign)
		sb.WriteString(strings.TrimSpace(a))
	}
	// return result
	return sb.String()
}

// "a", "b", "," -> "a,b")
// 字符串拼接, 默认','间隔
func StrJoin(str *string, addstr any, separators ...string) {
	sign := ","
	if len(separators) > 0 {
		sign = string(separators[0])
	}
	var add string
	switch a := addstr.(type) {
	case string:
		add = a
	default:
		add = fmt.Sprintf("%v", a)
	}

	if strings.TrimSpace(*str) == "" {
		*str = add
	} else if add != "" {
		*str = *str + sign + add
	}

	// // v2.0
	// var sb strings.Builder
	// sb.WriteString(*str)
	// if add != "" {
	// 	if !(sb.Len() > 0) {
	// 		sb.WriteString(add)
	// 	} else {
	// 		sb.WriteString(s + add)
	// 	}
	// }
	// *str = sb.String()
}

// 字符切片倒叙 slice array
func StringOrderDesc(slice []string) (array []string) {
	for _, v := range slice {
		array = append([]string{v}, array...)
	}
	return
}

// 123AB 将数字和字母分离
func StringSubmatch(input string) (number int, letter string) {
	// // 示例字符串
	// input := "123A"

	// 使用正则表达式匹配字符串
	re := regexp.MustCompile(`^(\d{1,3})([A-Za-z])$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) == 3 {
		// 提取数字部分并转换为整型
		var err error
		number, err = strconv.Atoi(matches[1])
		if err != nil {
			fmt.Println("StringSubmatch Error converting number:", matches[1], err)
		}

		// 提取字母部分
		letter = matches[2]
	} else {
		fmt.Println("StringSubmatch Input does not match the expected format")
	}

	return
}

// 输出百分比 number*100 + `%`
func PerCentStr(number float64) (pct string) {
	number *= 100
	if number > 1 {
		pct = fmt.Sprintf("%.2f%%", number)
	} else {
		pct = fmt.Sprintf("%.3g%%", number)
	}
	return
}

// 字符串切割比较
func CutSliceStr(symbol string) bool {
	ok := false
	last4Chars := symbol[len(symbol)-4:]
	// parts := strings.Split(symbol, "_")
	// if len(parts) > 1 { // || len(parts) > 1
	// 	last4Chars = parts[0][len(parts[0])-4:]
	// }
	if last4Chars == "USDT" {
		ok = true
	}

	return ok
}
