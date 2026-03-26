package seeder_models

import (
	"GEWIS-Rooster/internal/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	"strconv"
)

func SeedUser(d *gorm.DB, count int) {
	var organs []models.Organ

	if err := d.Find(&organs).Error; err != nil {
		log.Printf("Seeder Error")
	}

	for i := 0; i < count; i++ {
		user := models.User{
			Name:    "User" + strconv.Itoa(i),
			GEWISID: uint(1000 + i),
		}

		if err := d.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %d: %v", i, err)
		}
	}
}

func SeedUserOrgan(d *gorm.DB) {
	var organs []models.Organ
	if err := d.Find(&organs).Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to find organs")
	}

	var users []models.User
	if err := d.Find(&users).Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to find users")
	}

	roles := []models.OrganRole{
		models.RoleMember,
		models.RoleAdmin,
		models.RoleOwner,
	}

	for _, organ := range organs {
		for _, user := range users {
			// Pick a random role
			randomRole := roles[rand.Intn(len(roles))]

			userOrgan := models.UserOrgan{
				UserID:   user.ID,
				OrganID:  organ.ID,
				Username: fmt.Sprintf("User_%d_%d", organ.ID, user.ID),
				Role:     randomRole,
			}

			// 'UpdateAll' ensures that even if the link exists,
			// the Role and Username get updated to our seeded values.
			err := d.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "organ_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"role", "username"}),
			}).Create(&userOrgan).Error

			if err != nil {
				log.Error().Err(err).Msg("Failed to link User to Organ")
			}
		}
	}
}
