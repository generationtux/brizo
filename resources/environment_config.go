package resources

import (
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// EnvironmentConfig as defined by Brizo.
type EnvironmentConfig struct {
	database.Model
	UUID            string      `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name            string      `gorm:"not null" json:"name"`
	Value           string      `gorm:"not null" json:"value"`
	EnvironmentUUID string      `gorm:"not null" json:"environment_uuid"`
	Environment     Environment `gorm:"not null" json:"environment"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (config *EnvironmentConfig) BeforeCreate() (err error) {
	if config.UUID == "" {
		config.UUID = uuid.New()
	}

	return
}

// CreateEnvironmentConfig will persist a new configuration into database
func CreateEnvironmentConfig(db *gorm.DB, config *EnvironmentConfig) (bool, error) {
	// @TODO We need to handle duplicates. Adding a unique_index between the name,
	// value, and environment_uuid columns or doing a lookup before insert are two
	// possibilities.
	persist := db.Create(&config)
	return persist.RowsAffected == 1, persist.Error
}

// GetEnvironmentConfig will return the configuation values for specified environment UUID
func GetEnvironmentConfig(db *gorm.DB, uuid string) (*[]EnvironmentConfig, error) {
	var config []EnvironmentConfig

	if err := db.Where("environment_uuid = ?", uuid).Find(&config).Error; err != nil {
		return &config, err
	}

	return &config, nil
}

// DeleteEnvironmentConfig will delete specified configuration
func DeleteEnvironmentConfig(db *gorm.DB, uuid string) (bool, error) {
	var config EnvironmentConfig
	result := db.Model(config).Where("uuid = ?", uuid).First(&config)

	if result.RowsAffected == 1 {
		result = db.Delete(&config)
	}

	return result.RowsAffected == 1, result.Error
}
