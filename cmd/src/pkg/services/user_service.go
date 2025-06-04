package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	Create(*models.UserCreateOrUpdate) (*models.User, error)
	GetUser(uint) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(uint, *models.UserCreateOrUpdate) (*models.User, error)
	Delete(uint) error
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(createParams *models.UserCreateOrUpdate) (*models.User, error) {
	user := models.User{
		Name:    *createParams.Name,
		GEWISID: *createParams.GEWISID,
		Organs:  createParams.Organs,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUser(gewisId uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("gewis_id = ?", gewisId).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetAll() ([]*models.User, error) {
	var users []*models.User

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) Update(ID uint, updateParams *models.UserCreateOrUpdate) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, ID).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&user).Updates(updateParams).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) Delete(ID uint) error {
	var user models.User

	if err := s.db.Delete(&user, ID).Error; err != nil {
		return err
	}

	return nil
}
