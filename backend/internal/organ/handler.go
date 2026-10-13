package organ

import (
	"GEWIS-Rooster/internal/models"
	_ "GEWIS-Rooster/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	organService Service
}

func NewOrganHandler(rg *gin.RouterGroup, organService Service, db *gorm.DB) *Handler {
	h := &Handler{organService: organService}

	g := rg.Group("/organ")

	g.GET("/:id", h.GetMembersSettings)
	g.GET("/:id/member/:userId", h.GetMemberSettings)
	g.PATCH("/:id/member/:userId", h.UpdateMemberSettings)

	g.PATCH("/:id/member/:userId/role", requireRosterOrganMemberRoleParam(db, models.RoleAdmin), h.UpdateMemberRole)

	return h
}

// GetMembersSettings
//
//	@Summary      Get settings for all members within an organ
//	@Security     BearerAuth
//	@Description  Get organ-specific settings like nickname/username for all its members
//	@Tags         Organ
//	@Accept       json
//	@Produce      json
//	@Param        id             path      uint                                true  "Organ ID"
//	@Success      200            {array}  models.UserOrgan
//	@Failure      400            {string}  string
//	@Failure 404 {string} string
//	@ID	getMembersSettings
//	@Router       /organ/{id} [get]
func (o *Handler) GetMembersSettings(c *gin.Context) {
	organID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Organ ID")
		return
	}

	settings, err := o.organService.GetMembersSettings(uint(organID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not find Organ"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, settings)
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
func (o *Handler) GetMemberSettings(c *gin.Context) {
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
//	@Param        updateParams   body      organ.UpdateMemberSettingsParams   true  "Settings input"
//	@Success      200            {object}  models.UserOrgan
//	@Failure      400            {string}  string
//	@Router       /organ/{id}/member/{userId} [patch]
func (o *Handler) UpdateMemberSettings(c *gin.Context) {
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

	var params UpdateMemberSettingsParams
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

// UpdateMemberRole
//
//	@Summary      Update the role of a user within a specific organ
//	@Security     BearerAuth
//	@Description  Update a users role within a specific organ
//	@Tags         Organ
//	@Accept       json
//	@Produce      json
//	@Param        id             path      uint                                true  "Organ ID"
//	@Param        userId         path      uint                                true  "User ID"
//	@Param        updateParams   body      organ.UpdateMemberRoleParams   true  "Settings input"
//	@Success      200            {object}  models.UserOrgan
//	@Failure      400            {string}  string
//	@Failure 	  404			 {string}  string
//	@Router       /organ/{id}/member/{userId}/role [patch]
func (o *Handler) UpdateMemberRole(c *gin.Context) {
	organID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, "Invalid Organ ID")
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, "Invalid User ID")
		return
	}

	var params UpdateMemberRoleParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := o.organService.UpdateMemberRole(uint(organID), uint(userID), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
