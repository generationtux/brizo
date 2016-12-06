package auth

import (
	"github.com/generationtux/brizo/database"
)

// User represents a Brizo user
type User struct {
	database.Model
	Username       string `gorm:"not null;unique_index"`
	FirstName      string
	LastName       string
	Email          string
	GithubUsername string
	GithubToken    string
}
