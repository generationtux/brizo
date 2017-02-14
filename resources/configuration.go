package resources

import (
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

// Configuration as defined by Brizo.
type Configuration struct {
	database.Model
	Name          string      `gorm:"not null" json:"name"`
	Value         string      `gorm:"not null" json:"slug"`
	EnvironmentID uint        `gorm:"not null" json:"environment_id"`
	Environment   Environment `gorm:"not null" json:"environment"`
}

// CreateConfiguration will persist a new configuration into database
func CreateConfiguration(db *gorm.DB, config *Configuration) (bool, error) {
	persist := db.Create(&config)
	return persist.RowsAffected == 1, persist.Error
}

// GetConfiguration will return the configuation values for specified environment UUID
func GetConfiguration(db *gorm.DB, uuid string) (*[]Configuration, error) {
	var config []Configuration

	if err := db.Preload("Configuaration.Environment").Where("environment_uuid = ?", uuid).Find(&config).Error; err != nil {
		return &config, err
	}

	return &config, nil
}
