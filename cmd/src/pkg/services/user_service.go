package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"errors"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	Create(*models.UserCreateRequest) (*models.User, error)
	GetUsers(*models.UserFilterParams) ([]*models.User, error)
	Delete(uint) error
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(createParams *models.UserCreateRequest) (*models.User, error) {
	if createParams.Name == "" {
		return nil, errors.New("name is required")
	}

	user := models.User{
		Name:    createParams.Name,
		GEWISID: createParams.GEWISID,
		Organs:  createParams.Organs,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUsers(filters *models.UserFilterParams) ([]*models.User, error) {
	db := s.db.Model(&models.User{}).Preload("Organs")

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

func (s *UserService) Delete(ID uint) error {
	var user models.User
	result := s.db.Delete(&user, ID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
