package auth

import (
	"github.com/jinzhu/gorm"
)

// isFirstUser will determine if any users have registered yet
func isFirstUser(db *gorm.DB) bool {
	var users []User
	db.Find(&users)

	return len(users) == 0
}
