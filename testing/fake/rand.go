package fake

import "math/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandInt returns random int.
func RandInt() int {
	return rand.Int()
}

// RandInt32 returns random int32.
func RandInt32() int32 {
	return rand.Int31()
}

// RandInt32n returns random number from [0, n).
func RandInt32n(n int32) int32 {
	return rand.Int31n(n)
}

// RandInt64 returns random int64.
func RandInt64() int64 {
	return rand.Int63()
}

// RandUint32 returns random uint32.
func RandUint32() uint32 {
	return uint32(rand.Int31())
}

// RandString returns random string.
func RandString() string {
	length := 10
	kindChars := int32(len(charset))
	b := make([]byte, length)

	for i := 0; i < length; i++ {
		b[i] = charset[int(RandInt32n(kindChars))]
	}

	return string(b)
}
