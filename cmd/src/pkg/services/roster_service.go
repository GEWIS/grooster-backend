package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slices"
)

type RosterServiceInterface interface {
	CreateRoster(*models.RosterCreateRequest) (*models.Roster, error)
	GetRosters(*string) ([]*models.RosterResponse, error)
	GetRoster(uint) (*models.RosterResponse, error)
	UpdateRoster(uint, *models.RosterUpdateRequest) (*models.Roster, error)
	DeleteRoster(ID uint) error

	CreateRosterShift(*models.RosterShiftCreateRequest) (*models.RosterShift, error)
	DeleteRosterShift(ID uint) error
	CreateRosterAnswer(*models.RosterAnswerCreateRequest) (*models.RosterAnswer, error)
	UpdateRosterAnswer(uint, *models.RosterAnswerUpdateRequest) (*models.RosterAnswer, error)

	SaveRoster(uint) error
	UpdateSavedShift(uint, *models.SavedShiftUpdateRequest) (*models.SavedShift, error)
	GetSavedRoster(uint) ([]*models.SavedShift, error)
}

type RosterService struct {
	db *gorm.DB
}

func NewRosterService(db *gorm.DB) *RosterService {
	return &RosterService{db: db}
}

func (s *RosterService) CreateRoster(createParams *models.RosterCreateRequest) (*models.Roster, error) {
	var users []*models.User
	var values = models.Values{"Ja", "X", "L", "Nee"}

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	roster := models.Roster{
		Name:   createParams.Name,
		Values: values,
	}

	if err := s.db.Create(&roster).Error; err != nil {
		return nil, err
	}

	return &roster, nil
}

func (s *RosterService) GetRosters(date *string) ([]*models.RosterResponse, error) {
	var rosters []*models.Roster
	if err := s.db.Preload("RosterShift").Preload("RosterAnswer").Where("Date > ?", date).Find(&rosters).Error; err != nil {
		return nil, err
	}

	var rosterResponses []*models.RosterResponse
	for _, roster := range rosters {
		var users []*models.User
		if err := s.db.
			Joins("JOIN user_organs ON users.id = user_organs.user_id").
			Where("user_organs.organ_id = ?", roster.Organ).
			Find(&users).Error; err != nil {
			return nil, err
		}

		rosterResponse := models.RosterResponse{
			Roster: roster,
			Users:  users,
		}
		rosterResponses = append(rosterResponses, &rosterResponse)
	}

	return rosterResponses, nil
}

func (s *RosterService) GetRoster(ID uint) (*models.RosterResponse, error) {
	var roster models.Roster
	if err := s.db.Preload(clause.Associations).First(&roster, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	var users []*models.User
	if err := s.db.
		Joins("JOIN user_organs ON users.id = user_organs.user_id").
		Where("user_organs.organ_id = ?", roster.Organ).
		Find(&users).Error; err != nil {
		return nil, err
	}

	rosterResponse := models.RosterResponse{
		Roster: &roster,
		Users:  users,
	}

	return &rosterResponse, nil
}

func (s *RosterService) UpdateRoster(ID uint, updateParams *models.RosterUpdateRequest) (*models.Roster, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, ID).Error; err != nil {
		return nil, err
	}

	//if updateParams.UserIDs != nil {
	//	var users []models.User
	//	if err := s.db.Where("id IN ?", *updateParams.UserIDs).Find(&users).Error; err != nil {
	//		return nil, err
	//	}
	//	if err := s.db.Model(&roster).Association("Users").Replace(&users); err != nil {
	//		return nil, err
	//	}
	//}

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

func (s *RosterService) CreateRosterAnswer(createParams *models.RosterAnswerCreateRequest) (*models.RosterAnswer, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, createParams.RosterID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	var rosterShift *models.RosterShift
	if err := s.db.First(&rosterShift, roster.ID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	if !slices.Contains(roster.Values, createParams.Value) {
		return nil, fmt.Errorf("%s is not a valid value for this roster", createParams.Value)
	}

	rosterAnswer := models.RosterAnswer{
		UserID:        createParams.UserID,
		RosterID:      roster.ID,
		RosterShiftID: createParams.RosterShiftID,
		Value:         createParams.Value,
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

	return answer, nil
}

func (s *RosterService) SaveRoster(ID uint) error {
	var roster *models.Roster
	if err := s.db.Preload("RosterShift").First(&roster, ID).Error; err != nil {
		return err
	}

	for _, shift := range roster.RosterShift {
		if err := s.createSavedShift(roster.ID, shift); err != nil {
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
	}
	log.Print(saved)
	return saved, nil
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
