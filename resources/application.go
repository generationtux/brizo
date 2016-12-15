package resources

import (
	"errors"

	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Application as defined by Brizo.
type Application struct {
	database.Model
	UUID string `gorm:"not null;unique_index" sql:"type:varchar(36)"`
	Name string `gorm:"not null;unique_index"`
	Pods []Pod  `gorm:"-"` // gorm will ignore, but we can populate
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
	result := db.Find(&apps)

	return apps, result.Error
}

// CreateApplication will add a new Application to Brizo
func CreateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Create(&app)

	return result.RowsAffected == 1, result.Error
}

// UpdateApplication will update an existing Application
func UpdateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Model(app).Where("name = ?", app.Name).UpdateColumns(app)

	return result.RowsAffected == 1, result.Error
}

// GetApplication will get an existing Application by name
func GetApplication(db *gorm.DB, id string, getPods PodRetrieval) (*Application, error) {
	if _, valid := uuid.Parse(id).Version(); valid {
		app := new(Application)
		if err := db.Where("uuid = ?", id).First(&app).Error; err != nil {
			return app, err
		}

		pods, err := getPods(app.UUID)
		app.Pods = pods

		return app, err
	}

	return new(Application), errors.New("invalid id uuid version or format")
}

// DeleteApplication will delete an existing Application by name
func DeleteApplication(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Application{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
