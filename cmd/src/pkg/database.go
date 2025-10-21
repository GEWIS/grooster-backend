package database

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func ConnectDB(name string) *gorm.DB {
	devType := os.Getenv("DEV_TYPE")

	var db *gorm.DB
	var err error

	if devType == "production" {
		dsn := os.Getenv("DATABASE_DSN")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})

	}

	if err != nil {
		panic("failed to connect database")
	}

	if devType != "production" {
		db.Exec("PRAGMA foreign_keys = ON")
	}

	if err := db.AutoMigrate(&models.User{}, &models.Organ{}, &models.Roster{}, &models.RosterShift{}, &models.RosterAnswer{}, &models.SavedShift{}, &models.RosterTemplate{}); err != nil {
		panic(err)
	}

	return db
}
