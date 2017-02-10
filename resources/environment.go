package resources

import (
	"errors"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Environment as defined by Brizo.
type Environment struct {
	database.Model
	UUID          string      `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name          string      `gorm:"not null" json:"name"`
	Slug          string      `gorm:"not null" json:"slug"`
	ApplicationID uint64      `json:"application_id,string"`
	Application   Application `json:"application,array"`
	Versions      []Version   `json:"versions,array"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (environment *Environment) BeforeCreate() (err error) {
	if environment.UUID == "" {
		environment.UUID = uuid.New()
	}
	if environment.Slug == "" {
		environment.Slug = slugify.Slugify(environment.Name)
	}
	return
}

// AllEnvironments will return all of the Environments
func AllEnvironments(db *gorm.DB) ([]Environment, error) {
	var environments []Environment
	result := db.Find(&environments)

	return environments, result.Error
}

// CreateEnvironment will add a new Environment to Brizo
func CreateEnvironment(db *gorm.DB, environment *Environment) (bool, error) {
	result := db.Create(&environment)

	return result.RowsAffected == 1, result.Error
}

// UpdateEnvironment will update an existing Environment
func UpdateEnvironment(db *gorm.DB, environment *Environment) (bool, error) {
	environment.Slug = slugify.Slugify(environment.Name)
	result := db.Model(environment).Where("uuid = ?", environment.UUID).
		UpdateColumns(environment)

	return result.RowsAffected == 1, result.Error
}

// GetEnvironment will get an existing Environment by id
func GetEnvironment(db *gorm.DB, id string) (*Environment, error) {
	environment := new(Environment)
	if err := db.Preload("Application").Where("uuid = ?", id).First(&environment).Error; err != nil {
		return environment, err
	}

	if environment.ID == 0 {
		return new(Environment), errors.New("not-found")
	}

	return environment, nil
}

// DeleteEnvironment will delete an existing Environment by name
func DeleteEnvironment(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Environment{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
