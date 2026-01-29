package seeder

import (
	"GEWIS-Rooster/cmd/seeder/models"
	mainModels "GEWIS-Rooster/cmd/src/pkg/models"
	"fmt"
	"gorm.io/gorm"
)

func Seeder(d *gorm.DB) {
	err := WipeAllTables(d)

	if err != nil {
		panic(err)
	}

	err = d.AutoMigrate(
		&mainModels.User{},
		&mainModels.Roster{},
		&mainModels.RosterShift{},
		&mainModels.RosterAnswer{},
		&mainModels.SavedShift{},
		&mainModels.Organ{},
		&mainModels.RosterTemplate{},
		&mainModels.RosterTemplateShift{},
	)
	if err != nil {
		return
	}

	models.OrganSeeder(d, 2)
	models.SeedUser(d, 2)
	models.SeedRosters(d, 2)
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
