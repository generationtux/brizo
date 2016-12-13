package resources

import (
	"github.com/generationtux/brizo/database"
	"github.com/jinzhu/gorm"
)

type Application struct {
	database.Model
	Name string `gorm:"not null;unique_index"`
	Pods []Pod  `gorm:"-"` // gorm will ignore, but we can populate
}

func CreateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Create(&app)

	return result.RowsAffected == 1, result.Error
}

func UpdateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Model(app).Where("name = ?", app.Name).UpdateColumns(app)

	return result.RowsAffected == 1, result.Error
}

func GetApplication(db *gorm.DB, name string) (*Application, error) {
	app := new(Application)
	if err := db.Where("name = ?").First(&app).Error; err != nil {
		return app, err
	}

	return app, nil
}

func DeleteApplication(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Application{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
