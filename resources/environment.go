package resources

import (
	"errors"

	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Environment as defined by Brizo.
type Environment struct {
	database.Model
	UUID          string `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name          string `gorm:"not null;unique_index" json:"name"`
	ApplicationID uint64 `json:"application_id,string"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (a *Environment) BeforeCreate() (err error) {
	if a.UUID == "" {
		a.UUID = uuid.New()
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
	result := db.Model(environment).Where("name = ?", environment.Name).
		UpdateColumns(environment)

	return result.RowsAffected == 1, result.Error
}

// GetEnvironment will get an existing Environment by id
func GetEnvironment(db *gorm.DB, id string) (*Environment, error) {
	environment := new(Environment)
	if err := db.Where("uuid = ?", id).First(&environment).Error; err != nil {
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
