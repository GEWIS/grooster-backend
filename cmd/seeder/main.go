package main

import (
	database "GEWIS-Rooster/cmd/src/pkg"
	"GEWIS-Rooster/cmd/src/pkg/models"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	db := database.ConnectDB()
	sqlDB, _ := db.DB()

	err := wipeAllTables(db)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Print("failed to close database connection", err)
		}
	}(sqlDB)

	err = db.AutoMigrate(&models.User{}, &models.Roster{}, &models.RosterShift{}, &models.RosterAnswer{}, &models.SavedShift{}, &models.Organ{})
	if err != nil {
		return
	}

	Seeder(db)
}

func wipeAllTables(db *gorm.DB) error {
	var tableNames []string

	// List all tables to completely wipe them
	err := db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';").Scan(&tableNames).Error
	if err != nil {
		return err
	}

	// Drop each table
	for _, table := range tableNames {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
