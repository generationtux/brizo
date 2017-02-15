package resources

import (
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

// EnvironmentConfig as defined by Brizo.
type EnvironmentConfig struct {
	database.Model
	Name          string      `gorm:"not null" json:"name"`
	Value         string      `gorm:"not null" json:"value"`
	EnvironmentID uint        `gorm:"not null" json:"environment_id"`
	Environment   Environment `gorm:"not null" json:"environment"`
}

// CreateEnvironmentConfig will persist a new configuration into database
func CreateEnvironmentConfig(db *gorm.DB, config *EnvironmentConfig) (bool, error) {
	persist := db.Create(&config)
	return persist.RowsAffected == 1, persist.Error
}

// GetEnvironmentConfig will return the configuation values for specified environment UUID
func GetEnvironmentConfig(db *gorm.DB, uuid string) (*[]EnvironmentConfig, error) {
	var config []EnvironmentConfig

	if err := db.Preload("Configuaration.Environment").Where("environment_uuid = ?", uuid).Find(&config).Error; err != nil {
		return &config, err
	}

	return &config, nil
}
