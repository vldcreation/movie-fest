package util

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericBytes  = "0123456789"
	specialBytes  = "~=+%^*/()[]{}/?>;:<"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandString(n int) string {
	return Rand(n, letterBytes)
}

func Rand(n int, letter string) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)

	// Use int63() to generate random numbers more efficiently
	for i, cache, remain := n-1, src.Int63(), letterIdxBits; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxBits
		}
		b[i] = letter[int(cache&letterIdxMask)]
		i--
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandOwnersName() string {
	return RandString(6)
}

func RandAmount() int64 {
	return RandInt(0, 1000)
}

func RandEmail() string {
	return RandString(6) + "@gmail.com"
}

func RandSecret(n int) string {
	return Rand(n, letterBytes+numericBytes+specialBytes)
}
