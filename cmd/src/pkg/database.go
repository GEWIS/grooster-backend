package database

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(name string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err := db.AutoMigrate(&models.User{}, &models.Organ{}, &models.Roster{}, &models.RosterShift{}, &models.RosterAnswer{}, &models.SavedShift{}, &models.RosterTemplate{}); err != nil {
		panic(err)
	}

	return db
}
