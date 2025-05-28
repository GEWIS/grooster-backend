package database

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err := db.AutoMigrate(&models.User{}, &models.Roster{}, &models.RosterShift{}, &models.RosterAnswer{}, &models.SavedShift{}, &models.Organ{}); err != nil {
		panic(err)
	}

	return db
}
