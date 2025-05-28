package main

import (
	"GEWIS-Rooster/cmd/seeder/models"
	"gorm.io/gorm"
)

func Seeder(d *gorm.DB) {
	models.OrganSeeder(d, 2)
	models.SeedUser(d, 2)
	models.SeedRosters(d, 2)
}
