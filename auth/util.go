package auth

import (
	"math/rand"
	"time"
)

func generateRandomString(length int) string {
	const characters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[r.Intn(len(characters))]
	}

	return string(bytes)
}
