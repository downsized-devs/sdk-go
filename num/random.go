package num

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var now = time.Now

// Generate random string with n length
func RandomString(n int) string {
	seed := rand.New(rand.NewSource(now().UnixNano())) //nolint:gosec
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[seed.Intn(len(letterBytes))]
	}
	return string(b)
}
