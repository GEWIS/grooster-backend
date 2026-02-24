package organ

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"gorm.io/gorm"
)

type Service interface {
	GetMemberSettings(organID uint, userID uint) (*models.UserOrgan, error)
	UpdateMemberSettings(organID uint, userID uint, params *UpdateMemberSettingsParams) (*models.UserOrgan, error)
}

type service struct {
	db *gorm.DB
}

func NewOrganService(db *gorm.DB) Service {
	return &service{db: db}
}

func (o *service) GetMemberSettings(organID uint, userID uint) (*models.UserOrgan, error) {
	var userSettings models.UserOrgan
	err := o.db.Where("organ_id = ? AND user_id = ?", organID, userID).First(&userSettings).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &userSettings, nil
}

func (o *service) UpdateMemberSettings(organID uint, userID uint, params *UpdateMemberSettingsParams) (*models.UserOrgan, error) {
	updates := make(map[string]interface{})

	if params.Username != nil {
		updates["username"] = *params.Username
	}

	if len(updates) == 0 {
		var current models.UserOrgan
		err := o.db.Where("organ_id = ? AND user_id = ?", organID, userID).First(&current).Error

		if err != nil {
			return nil, err
		}

		return &current, nil
	}

	err := o.db.Model(&models.UserOrgan{}).
		Where("organ_id = ? AND user_id = ?", organID, userID).
		Updates(updates).Error

	if err != nil {
		return nil, err
	}

	var updatedRecord models.UserOrgan
	err = o.db.Where("organ_id = ? AND user_id = ?", organID, userID).
		First(&updatedRecord).Error

	if err != nil {
		return nil, err
	}

	return &updatedRecord, nil
}
