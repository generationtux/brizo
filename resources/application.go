package resources

import (
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

// Application as defined by Brizo.
type Application struct {
	database.Model
	Name string `gorm:"not null;unique_index"`
	Pods []Pod  `gorm:"-"` // gorm will ignore, but we can populate
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
func GetApplication(db *gorm.DB, name string) (*Application, error) {
	app := new(Application)
	if err := db.Where("name = ?").First(&app).Error; err != nil {
		return app, err
	}

	return app, nil
}

// DeleteApplication will delete an existing Application by name
func DeleteApplication(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Application{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
