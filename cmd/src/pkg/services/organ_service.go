package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"gorm.io/gorm"
)

type OrganServiceInterface interface {
	UpdateMemberSettings(organID uint, userID uint, params *models.UpdateMemberSettingsParams) (*models.UserOrgan, error)
}

type OrganService struct {
	db *gorm.DB
}

func NewOrganService(db *gorm.DB) *OrganService {
	return &OrganService{db: db}
}

func (o *OrganService) UpdateMemberSettings(organID uint, userID uint, params *models.UpdateMemberSettingsParams) (*models.UserOrgan, error) {
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
