package util

import (
	"time"

	"golang.org/x/exp/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

// GenerateRandomString generates a random string in an efficient way.
// @see https://stackoverflow.com/a/31832326/6335286
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
