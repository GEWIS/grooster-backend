package models

import (
	"GEWIS-Rooster/internal/organ"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
)

func OrganSeeder(d *gorm.DB, count int) {
	for i := 0; i < count; i++ {
		user := organ.Organ{
			Name: "Organ" + strconv.Itoa(i),
		}

		if err := d.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %d: %v", i, err)
		}
	}
}
