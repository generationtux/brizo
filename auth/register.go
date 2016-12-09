package auth

import (
	"github.com/jinzhu/gorm"
)

// IsFirstUser will determine if any users have registered yet
func IsFirstUser(db *gorm.DB) bool {
	var users []User
	db.Find(&users)

	return len(users) == 0
}

// GithubUserAllowed will determine if the provided username from Github is allowed
func GithubUserAllowed(db *gorm.DB, username string) bool {
	var users []User
	db.Where("github_username = ?", username).Find(&users)

	return len(users) > 0
}
