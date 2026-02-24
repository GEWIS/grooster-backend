package user

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(*CreateRequest) (*models.User, error)
	Get(*FilterParams) ([]*models.User, error)
	Delete(uint) error
}

type service struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) Create(createParams *CreateRequest) (*models.User, error) {
	if createParams.Name == "" {
		return nil, errors.New("name is required")
	}

	var userOrgans []models.Organ
	for _, id := range createParams.OrganIDs {
		userOrgans = append(userOrgans, models.Organ{
			BaseModel: models.BaseModel{ID: id},
		})
	}

	user := models.User{
		Name:    createParams.Name,
		GEWISID: createParams.GEWISID,
		Organs:  userOrgans,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) Get(filters *FilterParams) ([]*models.User, error) {
	db := s.db.Model(&models.User{}).Preload("Organs")
	log.Print(filters)
	if filters != nil {
		if filters.ID != nil {
			db = db.Where("id = ?", *filters.ID)
		}
		if filters.GEWISID != nil {
			db = db.Where("gewis_id = ?", *filters.GEWISID)
		}
		if filters.OrganID != nil {
			db = db.Joins("JOIN user_organs ON user_organs.user_id = users.id").
				Where("user_organs.organ_id = ?", *filters.OrganID)
		}
	}

	var users []*models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) Delete(ID uint) error {
	result := s.db.Unscoped().Delete(&models.User{}, ID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
