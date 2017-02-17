package resources

import (
	"errors"

	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Application as defined by Brizo.
type Application struct {
	database.Model
	UUID         string        `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name         string        `gorm:"not null;unique_index" json:"name"`
	Slug         string        `gorm:"not null;unique_index" json:"slug"`
	Environments []Environment `json:"environments,array"`
}

// BeforeCreate is a hook that runs before inserting a new record into the database
func (a *Application) BeforeCreate() (err error) {
	if a.UUID == "" {
		a.UUID = uuid.New()
	}
	return
}

// AllApplications will return all of the Applications
func AllApplications(db *gorm.DB) ([]Application, error) {
	var apps []Application
	result := db.Preload("Environments").Find(&apps)

	return apps, result.Error
}

// CreateApplication will add a new Application to Brizo
func CreateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Create(&app)

	return result.RowsAffected == 1, result.Error
}

// UpdateApplication will update an existing Application
func UpdateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Model(app).Where("uuid = ?", app.UUID).UpdateColumns(app)

	return result.RowsAffected == 1, result.Error
}

// GetApplication will get an existing Application by uuid
func GetApplication(db *gorm.DB, client kube.APIInterface, id string) (*Application, error) {
	app := new(Application)
	if err := db.Where("uuid = ?", id).Preload("Environments").First(&app).Error; err != nil {
		return app, err
	}

	if app.ID == 0 {
		return new(Application), errors.New("not-found")
	}

	return app, nil
}

// GetApplicationByID will get an existing Application by id
func GetApplicationByID(db *gorm.DB, id string) (*Application, error) {
	app := new(Application)
	if err := db.Where("id = ?", id).Preload("Environments").First(&app).Error; err != nil {
		return app, err
	}

	if app.ID == 0 {
		return new(Application), errors.New("not-found")
	}

	return app, nil
}

// DeleteApplication will delete an existing Application by name
func DeleteApplication(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Application{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
