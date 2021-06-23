package random

import (
	"math/rand"
	"time"
)

func RandStr(n int) string {
	bytes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return randStr(bytes, n)
}

func RandDigitsStr(n int) string {
	bytes := []byte("0123456789")
	return randStr(bytes, n)
}

func randStr(origin []byte, n int) string {
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, origin[r.Intn(len(origin))])
	}
	return string(result)
}
