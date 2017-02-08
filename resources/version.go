package resources

import (
	"errors"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Version as defined by Brizo.
type Version struct {
	database.Model
	UUID        string      `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name        string      `gorm:"not null;unique_index" json:"name"`
	Slug        string      `gorm:"not null;unique_index" json:"slug"`
	Image       string      `gorm:"not null" json:"image"`
	Replicas    int         `gorm:"not null" sql:"DEFAULT:'0'" json:"replicas"`
	Application Application `gorm:"not null" json:"application_id"`
	Environment Environment `form:"not null" json:"environment_id"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (version *Version) BeforeCreate() (err error) {
	if version.UUID == "" {
		version.UUID = uuid.New()
	}
	if version.Slug == "" {
		version.Slug = slugify.Slugify(version.Name)
	}
	return
}

// AllVersions will return all of the Versions
func AllVersions(db *gorm.DB) ([]Version, error) {
	var versions []Version
	result := db.Find(&versions)

	return versions, result.Error
}

// CreateVersion will add a new Version to Brizo
func CreateVersion(db *gorm.DB, version *Version) (bool, error) {
	result := db.Create(&version)

	return result.RowsAffected == 1, result.Error
}

// UpdateVersion will update an existing Version
func UpdateVersion(db *gorm.DB, version *Version) (bool, error) {
	version.Slug = slugify.Slugify(version.Name)
	result := db.Model(version).Where("uuid = ?", version.UUID).
		UpdateColumns(version)

	return result.RowsAffected == 1, result.Error
}

// GetVersion will get an existing Version by id
func GetVersion(db *gorm.DB, id string) (*Version, error) {
	version := new(Version)
	if err := db.Where("uuid = ?", id).First(&version).Error; err != nil {
		return version, err
	}

	if version.ID == 0 {
		return new(Version), errors.New("not-found")
	}

	return version, nil
}

// DeleteVersion will delete an existing Version by name
func DeleteVersion(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Version{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
