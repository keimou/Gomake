package utils

import (
	"math/rand"
)

// 生成随机字符串
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 否定该用法, 采用事先用时间戳生成随机种子
// 生成随机整数
func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// 生成随机浮点数
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
