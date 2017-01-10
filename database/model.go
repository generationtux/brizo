package database

import (
	"time"
)

// Model base ORM struct
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id,string"`
	CreatedAt time.Time `json:"created_at,string"`
	UpdatedAt time.Time `json:"updated_at,string"`
}
