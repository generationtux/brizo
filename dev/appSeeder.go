package dev

import (
	"github.com/generationtux/brizo/resources"
	"github.com/jinzhu/gorm"
)

func seedApplications(db *gorm.DB) error {
	apps := []resources.Application{
		resources.Application{Name: "product-mocker"},
	}

	for _, app := range apps {
		err := db.Create(&app).Error
		if err != nil {
			return err
		}
	}

	return nil
}
