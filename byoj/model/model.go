package model

import (
	"byoj/shared/database"

	"gorm.io/gorm"
)

func InitModel(db *gorm.DB) error {
	return database.AutoMigrateTable(db, &User{})
}
