package main

import (
	"GEWIS-Rooster/cmd/seeder/seeder"
	"GEWIS-Rooster/internal/models"
	"GEWIS-Rooster/internal/platform/database"
	"database/sql"
	"github.com/rs/zerolog/log"
)

func main() {
	db := database.ConnectDB("local.db")
	sqlDB, _ := db.DB()

	err := seeder.WipeAllTables(db)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Print("failed to close database connection", err)
		}
	}(sqlDB)

	err = db.AutoMigrate(&models.User{}, &models.Roster{}, &models.RosterShift{}, &models.RosterAnswer{}, &models.SavedShift{}, &models.Organ{}, &models.RosterTemplate{}, &models.RosterTemplateShift{})
	if err != nil {
		return
	}

	seeder.Seeder(db)
}
