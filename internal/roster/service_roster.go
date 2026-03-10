package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"time"
)

type RosterManager interface {
	CreateRoster(*CreateRequest) (*models.Roster, error)
	GetRosters(*FilterParams) ([]*models.Roster, error)
	UpdateRoster(uint, *UpdateRequest) (*models.Roster, error)
	DeleteRoster(ID uint) error
}

func (s *service) CreateRoster(params *CreateRequest) (*models.Roster, error) {
	var users []models.User
	var values = models.Values{"J", "X", "L", "N"} //TODO Change this to input values

	err := s.db.Joins("JOIN user_organs ON user_organs.user_id = users.id").
		Where("user_organs.organ_id = ?", params.OrganID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	if !isTodayOrLater(params.Date) {
		return nil, errors.New("date must be after the current date")
	}
	if params.Name == "" {
		return nil, errors.New("name is required")
	}

	roster := models.Roster{
		Name:       params.Name,
		Date:       params.Date,
		OrganID:    params.OrganID,
		Values:     values,
		TemplateID: params.TemplateID,
	}

	if err := s.db.Create(&roster).Error; err != nil {
		return nil, err
	}

	groupMapping := make(map[string]*uint)

	if params.TemplateID != nil {
		var templateShifts []models.RosterTemplateShift
		s.db.Where("template_id = ?", params.TemplateID).Find(&templateShifts)

		for _, ts := range templateShifts {
			groupMapping[ts.ShiftName] = ts.ShiftGroupID
		}
	}

	if params.Shifts != nil && len(params.Shifts) > 0 {
		for index, shift := range params.Shifts {
			var groupID *uint
			if gID, ok := groupMapping[shift]; ok {
				groupID = gID
			}

			rosterShift := &models.RosterShift{
				Name:         shift,
				RosterID:     roster.ID,
				Order:        uint(index),
				ShiftGroupID: groupID,
			}

			if err := s.db.Create(&rosterShift).Error; err != nil {
				return nil, err
			}
		}
	}

	if err := s.db.Preload("Organ").Preload("RosterShift").First(&roster, roster.ID).Error; err != nil {
		return nil, err
	}

	return &roster, nil
}

func (s *service) GetRosters(params *FilterParams) ([]*models.Roster, error) {
	db := s.db.Model(&models.Roster{})

	if params.ID != nil {
		db = db.Where("id = ?", *params.ID)
	} else {
		oneWeekAgo := time.Now().UTC().AddDate(0, 0, -7)

		if params.Archived != nil && *params.Archived {
			db = db.Where("date < ?", oneWeekAgo)
		} else {
			db = db.Where("date >= ?", oneWeekAgo)
		}
	}

	if params.OrganID != nil {
		db = db.Where("organ_id = ?", *params.OrganID)
	}

	var rosters []*models.Roster
	db.
		Preload("RosterShift").
		Preload("Organ").
		Preload("RosterAnswer")

	if err := db.Find(&rosters).Error; err != nil {
		return nil, err
	}

	return rosters, nil
}

func (s *service) UpdateRoster(id uint, params *UpdateRequest) (*models.Roster, error) {
	var roster *models.Roster

	if err := s.db.First(&roster, id).Error; err != nil {
		return nil, err
	}

	if params.Date != nil && !isTodayOrLater(*params.Date) {
		return nil, errors.New("date must be after the current date")
	}

	if params.Date != nil {
		roster.Date = *params.Date
	}
	if params.Name != nil {
		roster.Name = *params.Name
	}
	if params.Saved != nil {
		roster.Saved = *params.Saved
	}

	if err := s.db.Save(&roster).Error; err != nil {
		return nil, err
	}

	return roster, nil
}

func (s *service) DeleteRoster(ID uint) error {
	var roster models.Roster
	if err := s.db.Delete(&roster, ID).Error; err != nil {
		return err
	}

	return nil
}
