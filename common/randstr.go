package common

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TestRand(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(RandString(10))
	}
}
