package organ

import "GEWIS-Rooster/internal/models"

type UpdateMemberSettingsParams struct {
	Username *string `json:"username"`
} // @name UpdateMemberSettingsParams

type UpdateMemberRoleParams struct {
	Role models.OrganRole `json:"role" binding:"required,oneof=admin member owner"`
} // @name UpdateMemberRoleParams
