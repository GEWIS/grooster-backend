package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"GEWIS-Rooster/cmd/src/pkg/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type OrganHandler struct {
	organService *services.OrganService
}

func NewOrganHandler(rg *gin.RouterGroup, organService *services.OrganService) *OrganHandler {
	h := &OrganHandler{organService: organService}

	g := rg.Group("/organ")

	g.GET("/:id/member/:userId", h.GetMemberSettings)
	g.PATCH("/:id/member/:userId", h.UpdateMemberSettings)

	return h
}

// GetMemberSettings
//
//	@Summary      Get settings for a user within an organ
//	@Security     BearerAuth
//	@Description  Get organ-specific settings like nickname/username
//	@Tags         Organ
//	@Accept       json
//	@Produce      json
//	@Param        id             path      uint                                true  "Organ ID"
//	@Param        userId         path      uint                                true  "User ID"
//	@Success      200            {object}  models.UserOrgan
//	@Failure      400            {string}  string
//	@Failure 404 {string} string
//	@Router       /organ/{id}/member/{userId} [get]
func (o *OrganHandler) GetMemberSettings(c *gin.Context) {
	organID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Organ ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid User ID")
		return
	}

	settings, err := o.organService.GetMemberSettings(uint(organID), uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not find Organ or User"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateMemberSettings
//
//	@Summary      Update settings for a user within an organ
//	@Security     BearerAuth
//	@Description  Update organ-specific settings like nickname/username
//	@Tags         Organ
//	@Accept       json
//	@Produce      json
//	@Param        id             path      uint                                true  "Organ ID"
//	@Param        userId         path      uint                                true  "User ID"
//	@Param        updateParams   body      models.UpdateMemberSettingsParams   true  "Settings input"
//	@Success      200            {object}  models.UserOrgan
//	@Failure      400            {string}  string
//	@Router       /organ/{id}/member/{userId} [patch]
func (o *OrganHandler) UpdateMemberSettings(c *gin.Context) {
	organID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Organ ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid User ID")
		return
	}

	var params models.UpdateMemberSettingsParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := o.organService.UpdateMemberSettings(uint(organID), uint(userID), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
