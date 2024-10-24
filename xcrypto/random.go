package xcrypto

import (
	"crypto/rand"
	"math"
	"math/big"
	"unsafe"
)

// Source: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandInt64n(n int64) int64 {
	r, _ := rand.Int(rand.Reader, big.NewInt((n)))
	return r.Int64()
}

func RandInt(n int) int64 {
	return RandInt64n(int64(n))
}

func RandInt63() int64 {
	return RandInt64n(math.MaxInt64)
}

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, RandInt63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = RandInt63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
