package database

import (
	"GEWIS-Rooster/internal/models"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func ConnectDB(name string) *gorm.DB {
	devType := os.Getenv("DEV_TYPE")

	var db *gorm.DB
	var err error
	var dsn string

	if devType == "production" {
		dsn = os.Getenv("DATABASE_DSN")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Fatal().Msgf("Failed to connect database: %v", err)
		}

		runDBMigrations(dsn)
	} else {
		dsn = "sqlite3://" + name
		db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})

		if err != nil {
			log.Fatal().Msgf("Failed to connect database: %v", err)
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

	return db
}

func runDBMigrations(dsn string) {
	d, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration source")
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize migration")
	}

	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close migration source")
		}
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	log.Info().Msg("Database migrations completed successfully")
}
