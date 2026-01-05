package utils

import "math/rand"

import "time"

// RandomString returns an alphanumeric random string of given length.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	// crypto/rand for better randomness
	if _, err := rand.Read(b); err == nil {
		for i := 0; i < n; i++ {
			b[i] = letters[int(b[i])%len(letters)]
		}
		return string(b)
	}
	// fallback
	for i := 0; i < n; i++ {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
