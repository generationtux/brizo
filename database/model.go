package database

import (
	"time"
)

// Model base ORM struct
type Model struct {
	ID        uint `gorm:"primary_key" json:"ID,string"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
