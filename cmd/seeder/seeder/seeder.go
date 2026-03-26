package seeder

import (
	"GEWIS-Rooster/cmd/seeder/seeder_models"
	"GEWIS-Rooster/internal/models"
	"GEWIS-Rooster/internal/platform/database"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func Seeder(name string) *gorm.DB {

	devType := os.Getenv("DATABASE_TYPE")

	var db *gorm.DB
	var err error
	var dsn string

	if devType == "mysql" {
		dsn = os.Getenv("DATABASE_DSN")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Fatal().Msgf("Failed to connect database: %v", err)
		}

		err := WipeAllTables(db)

		if err != nil {
			panic(err)
		}

		mysqlDsn := "mysql://" + dsn

		database.RunDBMigrations(mysqlDsn)
	} else {
		dsn = "sqlite3://" + name
		db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})

		if err != nil {
			log.Fatal().Msgf("Failed to connect database: %v", err)
		}

		err := WipeAllTables(db)

		if err != nil {
			panic(err)
		}

		db.Exec("PRAGMA foreign_keys = ON")

		if err := db.AutoMigrate(
			&models.User{},
			&models.Organ{},
			&models.UserOrgan{},
			&models.Roster{},
			&models.RosterShift{},
			&models.RosterAnswer{},
			&models.SavedShift{},
			&models.RosterTemplate{},
			&models.RosterTemplateShift{},
			&models.RosterTemplateShiftPreference{},
			&models.ShiftGroup{},
		); err != nil {
			panic(err)
		}
	}

	seeder_models.OrganSeeder(db, 5)
	var organCount int64
	db.Model(&models.Organ{}).Count(&organCount)
	log.Info().Msgf("Organs in DB: %d", organCount)

	// 4. Seed Users
	seeder_models.SeedUser(db, 10)
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	log.Info().Msgf("Users in DB: %d", userCount)

	seeder_models.SeedUserOrgan(db)
	var linkCount int64
	db.Model(&models.UserOrgan{}).Count(&linkCount)
	log.Info().Msgf("User-Organ Links in DB: %d", linkCount)

	if organCount > 0 && userCount > 0 {
		seeder_models.SeedRosters(db, 5)
	} else {
		log.Error().Msg("Skipping Rosters: Missing Users or Organs!")
	}

	return db
}

func WipeAllTables(db *gorm.DB) error {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return fmt.Errorf("failed to fetch tables: %w", err)
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Error().Msgf("Failed to drop table %s: %v", table, err)
		}
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	return nil
}
