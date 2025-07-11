package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slices"
	"time"
)

type RosterServiceInterface interface {
	CreateRoster(*models.RosterCreateRequest) (*models.Roster, error)
	GetRosters(*models.RosterFilterParams) ([]*models.Roster, error)
	UpdateRoster(uint, *models.RosterUpdateRequest) (*models.Roster, error)
	DeleteRoster(ID uint) error

	CreateRosterShift(*models.RosterShiftCreateRequest) (*models.RosterShift, error)
	DeleteRosterShift(ID uint) error
	CreateRosterAnswer(*models.RosterAnswerCreateRequest) (*models.RosterAnswer, error)
	UpdateRosterAnswer(uint, *models.RosterAnswerUpdateRequest) (*models.RosterAnswer, error)

	SaveRoster(uint) error
	UpdateSavedShift(uint, *models.SavedShiftUpdateRequest) (*models.SavedShift, error)
	GetSavedRoster(uint) ([]*models.SavedShift, error)

	CreateRosterTemplate(*models.RosterTemplateCreateRequest) (*models.RosterTemplate, error)
}

type RosterService struct {
	db *gorm.DB
}

func NewRosterService(db *gorm.DB) *RosterService {
	return &RosterService{db: db}
}

func (s *RosterService) CreateRoster(params *models.RosterCreateRequest) (*models.Roster, error) {
	var users []models.User
	var values = models.Values{"Ja", "X", "L", "Nee"} //TODO Change this to input values

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	if err := s.db.Find(&models.Organ{}, params.OrganID).Error; err != nil {
		return nil, err
	}
	if !isAfterToday(params.Date) {
		return nil, errors.New("date must be after the current date")
	}
	if params.Name == "" {
		return nil, errors.New("name is required")
	}

	roster := models.Roster{
		Name:    params.Name,
		Date:    params.Date,
		OrganID: params.OrganID,
		Values:  values,
	}

	if err := s.db.Create(&roster).Error; err != nil {
		return nil, err
	}

	if params.Shifts != nil && len(*params.Shifts) > 0 {
		for _, shift := range *params.Shifts {
			rosterShift := &models.RosterShift{
				Name:     shift,
				RosterID: roster.ID,
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

func (s *RosterService) GetRosters(params *models.RosterFilterParams) ([]*models.Roster, error) {
	db := s.db.Model(&models.Roster{})

	if params.ID != nil {
		db = db.Where("id = ?", *params.ID)
	}

	if params.Date != nil {
		db = db.Where("date = ?", *params.Date)
	}

	if params.OrganID != nil {
		db = db.Where("organ_id = ?", *params.OrganID)
	}

	db.
		Preload("RosterShift").
		Preload("RosterAnswer").
		Preload("Organ")

	var rosters []*models.Roster
	if err := db.Find(&rosters).Error; err != nil {
		return nil, err
	}

	return rosters, nil
}

func (s *RosterService) UpdateRoster(id uint, params *models.RosterUpdateRequest) (*models.Roster, error) {
	var roster *models.Roster

	if err := s.db.First(&roster, id).Error; err != nil {
		return nil, err
	}

	if params.Date != nil && !isAfterToday(*params.Date) {
		return nil, errors.New("date must be after the current date")
	}

	if params.Date != nil {
		roster.Date = *params.Date
	}
	if params.Name != nil {
		roster.Name = *params.Name
	}

	if err := s.db.Save(&roster).Error; err != nil {
		return nil, err
	}

	return roster, nil
}

func (s *RosterService) DeleteRoster(ID uint) error {
	var roster models.Roster
	if err := s.db.Delete(&roster, ID).Error; err != nil {
		return err
	}

	return nil
}

func (s *RosterService) CreateRosterShift(createParams *models.RosterShiftCreateRequest) (*models.RosterShift, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, createParams.RosterID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	rosterShift := models.RosterShift{
		Name:     createParams.Name,
		RosterID: createParams.RosterID,
	}

	if err := s.db.Create(&rosterShift).Error; err != nil {
		return nil, err
	}

	return &rosterShift, nil
}

func (s *RosterService) DeleteRosterShift(ID uint) error {
	var rosterShift *models.RosterShift

	if err := s.db.Delete(&rosterShift, ID).Error; err != nil {
		return err
	}

	return nil
}

func (s *RosterService) CreateRosterAnswer(params *models.RosterAnswerCreateRequest) (*models.RosterAnswer, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, params.RosterID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	var rosterShift *models.RosterShift
	if err := s.db.First(&rosterShift, params.RosterShiftID).Error; err != nil {
		return nil, fmt.Errorf("roster shift not found: %w", err)
	}

	if !slices.Contains(roster.Values, params.Value) {
		return nil, fmt.Errorf("%s is not a valid value for this roster", params.Value)
	}

	rosterAnswer := models.RosterAnswer{
		UserID:        params.UserID,
		RosterID:      roster.ID,
		RosterShiftID: params.RosterShiftID,
		Value:         params.Value,
	}

	if err := s.db.Create(&rosterAnswer).Error; err != nil {
		return nil, err
	}

	return &rosterAnswer, nil
}

func (s *RosterService) UpdateRosterAnswer(ID uint, updateParams *models.RosterAnswerUpdateRequest) (*models.RosterAnswer, error) {
	var answer *models.RosterAnswer

	if err := s.db.First(&answer, ID).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&answer).Updates(updateParams).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&answer, ID).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

func (s *RosterService) SaveRoster(ID uint) error {
	var roster *models.Roster
	if err := s.db.Preload("RosterShift").First(&roster, ID).Error; err != nil {
		return err
	}

	for _, shift := range roster.RosterShift {
		if err := s.createSavedShift(roster.ID, &shift); err != nil {
			return err
		}
	}

	roster.Saved = true
	if err := s.db.Save(&roster).Error; err != nil {
		return err
	}

	return nil
}

func (s *RosterService) GetSavedRoster(ID uint) ([]*models.SavedShift, error) {
	var savedShifts []*models.SavedShift
	if err := s.db.Preload(clause.Associations).Where("roster_id = ?", ID).Find(&savedShifts).Error; err != nil {
		return nil, err
	}

	return savedShifts, nil
}

func (s *RosterService) UpdateSavedShift(ID uint, updateParams *models.SavedShiftUpdateRequest) (*models.SavedShift, error) {
	var saved *models.SavedShift
	if err := s.db.Preload("Users").First(&saved, ID).Error; err != nil {

		return nil, err
	}

	if updateParams.UserIDs != nil {
		var users []*models.User
		if err := s.db.Where("ID IN ?", updateParams.UserIDs).Find(&users).Error; err != nil {

			return nil, err
		}
		// Replace existing users with the new set
		if err := s.db.Model(&saved).Association("Users").Replace(users); err != nil {
			return nil, err
		}
		// Reload associations to get fresh data
		if err := s.db.Preload("Users").Preload("RosterShift").First(&saved, ID).Error; err != nil {
			return nil, err
		}
	}

	return saved, nil
}

func (s *RosterService) CreateRosterTemplate(params *models.RosterTemplateCreateRequest) (*models.RosterTemplate, error) {
	var organ models.Organ

	if err := s.db.First(&organ, params.OrganID).Error; err != nil {
		return nil, err
	}

	if len(params.Shifts) == 0 {
		return nil, errors.New("no shifts were given")
	}

	template := models.RosterTemplate{
		OrganID: organ.ID,
		Organ:   &organ,
		Shifts:  params.Shifts,
	}

	if err := s.db.Create(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *RosterService) createSavedShift(rID uint, shift *models.RosterShift) error {
	var savedShift = models.SavedShift{
		RosterID:    rID,
		RosterShift: shift,
		Users:       []*models.User{},
	}

	if err := s.db.Create(&savedShift).Error; err != nil {
		return err
	}
	return nil
}

func isAfterToday(date time.Time) bool {
	today := time.Now().Truncate(24 * time.Hour)
	inputDate := date.Truncate(24 * time.Hour)

	return inputDate.After(today)
}
