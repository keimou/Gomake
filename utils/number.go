package utils

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// string 转 []int  e.g. "2,1,1"->[2,1,1]
func StrToInts(input string) []int {
	var intSlice []int

	// 使用正则表达式去除所有空白字符 (包括空格,制表符,换行符等)
	re := regexp.MustCompile(`\s+`)
	noSpacesMethods := re.ReplaceAllString(input, "")

	stringSlice := strings.Split(noSpacesMethods, ",")
	for _, v := range stringSlice {
		if num, err := strconv.Atoi(v); err == nil {
			intSlice = append(intSlice, num)
		}
	}

	return intSlice
}

// string 转 []int (已去重)  e.g. "2,3,2,1"->[2,3,1]
func StrConvInt(input string) []int {
	var intSlice []int

	// 使用正则表达式去除所有空白字符 (包括空格,制表符,换行符等)
	re := regexp.MustCompile(`\s+`)
	noSpacesMethods := re.ReplaceAllString(input, "")

	// 创建一个 map 用于跟踪已经遇到的元素
	seen := make(map[int]bool)

	stringSlice := strings.Split(noSpacesMethods, ",")
	for _, v := range stringSlice {
		num, err := strconv.Atoi(v)
		// 为整数且元素尚未出现过, 添加到 intSlice 并标记为已出现
		if err == nil && !seen[num] {
			intSlice = append(intSlice, num)
			seen[num] = true
		}
	}

	return intSlice
}

// string 转 []uint (已排序+去重) e.g. "3,2,1"->[1,2,3]
func StrConvUint(input string) []uint {

	intSlice := StrConvInt(input)
	uintSlice := make([]uint, len(intSlice))
	sort.Ints(intSlice)
	for i, v := range intSlice {
		uintSlice[i] = uint(v)
	}

	return uintSlice
}

// 	sort.Ints(output) // 使用 sort.Ints 对(output []int)进行排序

// 将数值型字符串翻译为数值, 以供计算或比较
func Number(value string) float64 {
	str := strings.TrimSpace(value) // 过滤字符串前后多余空格
	num, _ := strconv.ParseFloat(str, 64)
	return num
}

func Max(nums ...float64) float64 {
	if len(nums) == 0 {
		// fmt.Println("警告：切片为空，返回 0")
		return 0
	}
	max := nums[0]               // 将第一个元素作为初始最大值
	for _, v := range nums[1:] { // 从第二个元素开始遍历
		if v > max {
			max = v
		}
	}
	return max
}
