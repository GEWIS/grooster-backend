package models

import (
	"GEWIS-Rooster/internal/organ"
	"GEWIS-Rooster/internal/user"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
)

func SeedUser(d *gorm.DB, count int) {
	var organs []organ.Organ

	if err := d.Find(&organs).Error; err != nil {
		log.Printf("Seeder Error")
	}

	for i := 0; i < count; i++ {
		user := user.User{
			Name:    "User" + strconv.Itoa(i),
			GEWISID: uint(1000 + i),
			Organs:  organs,
		}

		if err := d.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %d: %v", i, err)
		}
	}
}
