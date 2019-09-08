package lib

import (
	"math/rand"
	"time"
)

func RandStr(length int) string {
	var res string
	var source = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		idx := rand.Int() % 36
		res += string(source[idx])
	}
	return res
}

func RandNumberStr(length int) string {
	var res string
	var source = "0123456789"
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		idx := rand.Int() % 10
		res += string(source[idx])
	}
	return res
}
