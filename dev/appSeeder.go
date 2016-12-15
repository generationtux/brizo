package dev

import (
	"github.com/generationtux/brizo/resources"
	"github.com/jinzhu/gorm"
)

func seedApplications(db *gorm.DB) error {
	apps := []resources.Application{
		Application{Name: "product-mocker"},
	}

	for _, app := range apps {
		err := db.Create(&app)
		if err != nil {
			return err
		}
	}

	return nil
}
