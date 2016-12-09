package auth

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringIsGenerated(t *testing.T) {
	generatedString := generateRandomString(32)
	matches, _ := regexp.Match("[a-zA-Z0-9]", []byte(generatedString))

	assert.Len(t, generatedString, 32)
	assert.True(t, true, matches)
}
