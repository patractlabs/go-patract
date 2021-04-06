//nolint:gosec
package utils

import "math/rand"

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"

const SaltMaxLen = 64

func GenerateSalt() string {
	return RandStringBytesRmndr(SaltMaxLen)
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
