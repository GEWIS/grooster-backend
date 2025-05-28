package models

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
)

func SeedUser(d *gorm.DB, count int) {
	var organs []*models.Organ

	if err := d.Find(&organs).Error; err != nil {
		log.Printf("Seeder Error")
	}

	for i := 0; i < count; i++ {
		user := models.User{
			Name:   "User" + strconv.Itoa(i),
			Organs: organs,
		}

		if err := d.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %d: %v", i, err)
		}
	}
}
