package main

import (
	"GEWIS-Rooster/cmd/seeder/seeder"
	"GEWIS-Rooster/internal/organ"
	"GEWIS-Rooster/internal/platform/database"
	"GEWIS-Rooster/internal/roster"
	"GEWIS-Rooster/internal/user"
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

	err = db.AutoMigrate(&user.User{}, &roster.Roster{}, &roster.RosterShift{}, &roster.RosterAnswer{}, &roster.SavedShift{}, &organ.Organ{}, &roster.RosterTemplate{}, &roster.RosterTemplateShift{})
	if err != nil {
		return
	}

	seeder.Seeder(db)
}
