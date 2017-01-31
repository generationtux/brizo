package auth

import (
	"math/rand"
	"time"
)

// GenerateRandomString will provide a [0-9a-Z] string of a specified length.
func GenerateRandomString(length int) string {
	const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[r.Intn(len(characters))]
	}

	return string(bytes)
}
