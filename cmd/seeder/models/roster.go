package models

import (
	"GEWIS-Rooster/internal/organ"
	roster2 "GEWIS-Rooster/internal/roster"
	"GEWIS-Rooster/internal/user"
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

func roster(db *gorm.DB, count int) []*roster2.Roster {
	var users []*user.User
	var organs []*organ.Organ
	var rosters []*roster2.Roster
	var values = roster2.Values{"Ja", "X", "L", "Nee"}

	if err := db.Find(&users).Error; err != nil {
		log.Printf("Could not get users: %v\n", err)
	}

	if err := db.Find(&organs).Error; err != nil {
		log.Printf("Could not get organs: %v\n", err)
	}

	for i := 0; i < count; i++ {
		r := &roster2.Roster{
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

func rosterShift(db *gorm.DB, roster []*roster2.Roster) {
	for _, roster := range roster {
		shifts := []roster2.RosterShift{
			{Name: "Shift A", RosterID: roster.ID, Order: 0},
			{Name: "Shift B", RosterID: roster.ID, Order: 1},
			{Name: "Shift C", RosterID: roster.ID, Order: 2},
		}

		if err := db.Create(&shifts).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create roster shifts")
		}
	}
}

func rosterAnswer(db *gorm.DB, roster []*roster2.Roster) {
	for _, roster := range roster {
		var users []*user.User
		var shifts []roster2.RosterShift

		var organ organ.Organ

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
		var answers []roster2.RosterAnswer
		for _, user := range users {
			for _, shift := range shifts {
				answers = append(answers, roster2.RosterAnswer{
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
	var organ organ.Organ
	if err := db.First(&organ).Error; err != nil {
		log.Error().Err(err).Msg("No organ found for seeding")
		return
	}

	for i := 0; i < count; i++ {
		shiftName := fmt.Sprintf("Shift %d", i)

		templateShifts := []roster2.RosterTemplateShift{
			{
				ShiftName: shiftName,
			},
		}

		template := roster2.RosterTemplate{
			OrganID: organ.ID,
			Name:    fmt.Sprintf("Template %d", i),
			Shifts:  templateShifts,
		}

		if err := db.Create(&template).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create roster templates")
		}
	}
}
