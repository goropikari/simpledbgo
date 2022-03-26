package fake

import "math/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandInt() int {
	return rand.Int()
}

func RandInt32() int32 {
	return rand.Int31()
}

// RandInt32n returns random number from [0, n).
func RandInt32n(n int32) int32 {
	return rand.Int31n(n)
}

func RandInt64() int64 {
	return rand.Int63()
}

func RandUint32() uint32 {
	return uint32(rand.Int31())
}

func RandString(length int) string {
	kindChars := int32(len(charset))
	b := make([]byte, length)

	for i := 0; i < length; i++ {
		b[i] = charset[int(RandInt32n(kindChars))]
	}

	return string(b)
}
