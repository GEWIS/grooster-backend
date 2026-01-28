package models

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func SeedRosters(db *gorm.DB, count int) {
	rosters := roster(db, count)
	rosterShift(db, rosters)
	rosterAnswer(db, rosters)
	rosterTemplates(db, count)
}

func roster(db *gorm.DB, count int) []*models.Roster {
	var users []*models.User
	var organs []*models.Organ
	var rosters []*models.Roster
	var values = models.Values{"Ja", "X", "L", "Nee"}

	if err := db.Find(&users).Error; err != nil {
		log.Printf("Could not get users: %v\n", err)
	}

	if err := db.Find(&organs).Error; err != nil {
		log.Printf("Could not get organs: %v\n", err)
	}

	for i := 0; i < count; i++ {
		r := &models.Roster{
			Name:    "Roster" + strconv.Itoa(i),
			Values:  values,
			OrganID: organs[i].ID,
			Organ:   *organs[i],
			Date:    time.Now(),
			Saved:   false,
		}
		if err := db.Create(r).Error; err != nil {
			log.Printf("Seeder Error: %v", err)
		} else {
			rosters = append(rosters, r)
		}
	}

	return rosters
}

func rosterShift(db *gorm.DB, roster []*models.Roster) {
	for _, roster := range roster {
		shifts := []models.RosterShift{
			{Name: "Shift A", RosterID: roster.ID, Order: 0},
			{Name: "Shift B", RosterID: roster.ID, Order: 1},
			{Name: "Shift C", RosterID: roster.ID, Order: 2},
		}

		if err := db.Create(&shifts).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create roster shifts")
		}
	}
}

func rosterAnswer(db *gorm.DB, roster []*models.Roster) {
	for _, roster := range roster {
		var users []*models.User
		var shifts []models.RosterShift

		var organ models.Organ

		err := db.First(&organ, roster.Organ).Error
		if err != nil {
			log.Error().Err(err).Msg("Could not get organ")
		}

		err = db.Model(&organ).Association("Users").Find(&users)
		if err != nil {
			log.Error().Err(err).Msg("Could not get users")
		}
		db.Where("roster_id = ?", roster.ID).Find(&shifts)

		values := roster.Values
		valueCount := len(values)

		i := 0
		var answers []models.RosterAnswer
		for _, user := range users {
			for _, shift := range shifts {
				answers = append(answers, models.RosterAnswer{
					UserID:        user.ID,
					RosterID:      roster.ID,
					RosterShiftID: shift.ID,
					Value:         values[i%valueCount],
				})
				i++
			}
		}

		if err := db.Create(&answers).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create roster answers")
		}
	}
}

func rosterTemplates(db *gorm.DB, count int) {
	var organ models.Organ
	if err := db.First(&organ).Error; err != nil {
		log.Error().Err(err).Msg("No organ found for seeding")
		return
	}

	for i := 0; i < count; i++ {
		shiftName := fmt.Sprintf("Shift %d", i)

		templateShifts := []models.RosterTemplateShift{
			{
				ShiftName: shiftName,
			},
		}

		template := models.RosterTemplate{
			OrganID: organ.ID,
			Name:    fmt.Sprintf("Template %d", i),
			Shifts:  templateShifts,
		}

		if err := db.Create(&template).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create roster templates")
		}
	}
}
