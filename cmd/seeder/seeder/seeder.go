package seeder

import (
	"GEWIS-Rooster/cmd/seeder/seeder_models"
	"GEWIS-Rooster/internal/models"
	"fmt"
	"gorm.io/gorm"
)

func Seeder(d *gorm.DB) {
	err := WipeAllTables(d)

	if err != nil {
		panic(err)
	}

	err = d.AutoMigrate(
		&models.User{},
		&models.Roster{},
		&models.RosterShift{},
		&models.RosterAnswer{},
		&models.SavedShift{},
		&models.Organ{},
		&models.RosterTemplate{},
		&models.RosterTemplateShift{},
	)
	if err != nil {
		return
	}

	seeder_models.OrganSeeder(d, 2)
	seeder_models.SeedUser(d, 2)
	seeder_models.SeedRosters(d, 2)
}

func WipeAllTables(db *gorm.DB) error {
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
