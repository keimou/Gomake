package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"os"
	"time"
)

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("the file of the same name")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 生成随机数 (性能不好, 待优化)
func Intn(n int) (num int) {
	// rand.Seed(time.Now().UnixNano()) // Crucial: seed the generator
	// num = rand.Intn(n)
	num = rand.Intn(int(time.Now().UnixMicro())) % n
	// fmt.Printf("num: %v\n", num)
	return
}
