package seeder_models

import (
	"GEWIS-Rooster/internal/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

func SeedRosters(db *gorm.DB, count int) {
	shiftGroups := seedShiftGroups(db)
	seedGroupPriorities(db, shiftGroups)

	templates := rosterTemplates(db, count, shiftGroups)
	seedTemplatePreferences(db, templates)

	rosters := roster(db, count)

	shifts := rosterShift(db, rosters, shiftGroups)

	rosterAnswer(db, rosters)
	seedSavedShifts(db, rosters, shifts)
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

func rosterShift(db *gorm.DB, rosters []*models.Roster, groups []models.ShiftGroup) []models.RosterShift {
	var allShifts []models.RosterShift
	for _, r := range rosters {
		var groupID *uint
		if len(groups) > 0 {
			id := groups[rand.Intn(len(groups))].ID
			groupID = &id
		}

		shifts := []models.RosterShift{
			{Name: "Shift A", RosterID: r.ID, Order: 0, ShiftGroupID: groupID},
			{Name: "Shift B", RosterID: r.ID, Order: 1, ShiftGroupID: groupID},
		}
		db.Create(&shifts)
		allShifts = append(allShifts, shifts...)
	}
	return allShifts
}

func rosterAnswer(db *gorm.DB, roster []*models.Roster) {
	for _, roster := range roster {
		var users []*models.User
		var shifts []models.RosterShift

		var newOrgan models.Organ

		err := db.First(&newOrgan, roster.Organ).Error
		if err != nil {
			log.Error().Err(err).Msg("Could not get organ")
		}

		err = db.Model(&newOrgan).Association("Users").Find(&users)
		if err != nil {
			log.Error().Err(err).Msg("Could not get users")
		}
		db.Where("roster_id = ?", roster.ID).Find(&shifts)

		values := roster.Values
		valueCount := len(values)

		i := 0
		var answers []models.RosterAnswer
		for _, newUser := range users {
			for _, shift := range shifts {
				answers = append(answers, models.RosterAnswer{
					UserID:        newUser.ID,
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

func rosterTemplates(db *gorm.DB, count int, groups []models.ShiftGroup) []models.RosterTemplate {
	var organ models.Organ
	db.First(&organ)

	var createdTemplates []models.RosterTemplate

	for i := 0; i < count; i++ {
		var groupID *uint
		if len(groups) > 0 {
			id := groups[rand.Intn(len(groups))].ID
			groupID = &id
		}

		template := models.RosterTemplate{
			OrganID: organ.ID,
			Name:    fmt.Sprintf("Template %d", i),
			Shifts: []models.RosterTemplateShift{
				{ShiftName: "Morning", ShiftGroupID: groupID},
				{ShiftName: "Evening", ShiftGroupID: groupID},
			},
		}
		db.Create(&template)
		createdTemplates = append(createdTemplates, template)
	}
	return createdTemplates
}

func seedShiftGroups(db *gorm.DB) []models.ShiftGroup {
	var organs []models.Organ
	db.Find(&organs)

	var groups []models.ShiftGroup
	for _, organ := range organs {
		for i := 1; i <= 2; i++ {
			group := models.ShiftGroup{
				OrganID: organ.ID,
				Name:    fmt.Sprintf("%s Group %d", organ.Name, i),
			}
			db.FirstOrCreate(&group, models.ShiftGroup{OrganID: organ.ID, Name: group.Name})
			groups = append(groups, group)
		}
	}
	return groups
}

func seedGroupPriorities(db *gorm.DB, groups []models.ShiftGroup) {
	var users []models.User
	db.Find(&users)

	priorities := []models.GroupPriority{models.Low, models.Default, models.High}

	var entries []models.ShiftGroupPriority
	for _, group := range groups {
		for _, user := range users {
			entries = append(entries, models.ShiftGroupPriority{
				ShiftGroupID: group.ID,
				UserID:       user.ID,
				Priority:     priorities[rand.Intn(len(priorities))],
			})
		}
	}
	db.Create(&entries)
}

func seedTemplatePreferences(db *gorm.DB, templates []models.RosterTemplate) {
	var users []models.User
	db.Find(&users)
	prefs := []string{"High", "Medium", "Low", "None"}

	var allPrefs []models.RosterTemplateShiftPreference
	for _, t := range templates {
		for _, shift := range t.Shifts {
			for _, user := range users {
				allPrefs = append(allPrefs, models.RosterTemplateShiftPreference{
					RosterTemplateShiftID: shift.ID,
					UserID:                user.ID,
					Preference:            prefs[rand.Intn(len(prefs))],
				})
			}
		}
	}
	db.CreateInBatches(&allPrefs, 100)
}

func seedSavedShifts(db *gorm.DB, rosters []*models.Roster, shifts []models.RosterShift) {
	for _, r := range rosters {
		for _, s := range shifts {
			if s.RosterID == r.ID {
				var users []models.User
				db.Limit(2).Find(&users)

				saved := models.SavedShift{
					RosterID:      r.ID,
					RosterShiftID: s.ID,
					Users:         nil,
				}
				db.Create(&saved)
				db.Model(&saved).Association("Users").Append(users)

			}
		}
	}
}
