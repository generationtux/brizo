package resources

import (
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Volume as defined by Brizo.
type Volume struct {
	database.Model
	UUID      string  `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name      string  `gorm:"not null" json:"name"`
	Type      string  `gorm:"not null" json:"type"`
	Source    string  `gorm:"not null" json:"source"`
	VersionID uint    `gorm:"not null" json:"version_id"`
	Version   Version `gorm:"not null" json:"version"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (volume *Volume) BeforeCreate() (err error) {
	if volume.UUID == "" {
		volume.UUID = uuid.New()
	}
	return
}

// AllVolumes will return all of the Volumes
func AllVolumes(db *gorm.DB) ([]Volume, error) {
	var volumes []Volume
	result := db.Find(&volumes)

	return volumes, result.Error
}

// CreateVolume will create a new volume
func CreateVolume(db *gorm.DB, volume *Volume) (bool, error) {
	// @todo kube-client work
	persist := db.Create(&volume)

	return persist.RowsAffected == 1, persist.Error
}

// GetVolumesByVersionUUID will get an existing volume using a version's UUID
func GetVolumesByVersionUUID(db *gorm.DB, uuid string) (*[]Volume, error) {
	var volumes []Volume
	version, err := GetVersion(db, uuid)
	if err != nil {
		return &volumes, err
	}

	if err = db.Preload("Version.Environment").Where("version_id = ?", version.ID).Find(&volumes).Error; err != nil {
		return &volumes, err
	}

	return &volumes, err
}
