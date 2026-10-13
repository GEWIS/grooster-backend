package roster

import (
	"GEWIS-Rooster/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *Handler) registerRosterRoutes(g *gin.RouterGroup, db *gorm.DB) {
	g.POST("", requireRosterOrganRoleBody(db, models.RoleAdmin), h.CreateRoster)
	g.GET("", requireRosterOrganRoleQuery(db, "organId", models.RoleMember), h.GetRosters)
	g.GET(":id", requireRosterOrganRoleParam(db, "id", models.RoleMember), h.GetRoster)
	g.PATCH("/:id", requireRosterOrganRoleParam(db, "id", models.RoleAdmin), h.UpdateRoster)
	g.DELETE("/:id", requireRosterOrganRoleParam(db, "id", models.RoleAdmin), h.DeleteRoster)
}

// CreateRoster
//
//	@Summary	CreateRoster a new roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		RosterCreateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			createRoster
//	@Router		/roster [post]
func (h *Handler) CreateRoster(c *gin.Context) {
	var param *CreateRequest

	// Due to the middleware checking the body, we also use ShouldBindBodyWith here
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdRoster, err := h.rosterService.CreateRoster(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdRoster)
}

// GetRosters
//
//		@Summary	Get all rosters or query by date and organ
//		@Security	BearerAuth
//		@Tags		Roster
//		@Accept		json
//		@Produce	json
//	 @Param     filter  query     RosterFilterParams  false  "Filter parameters"
//		@Success	200		{array}		models.Roster
//		@Failure	400		{string}	string
//		@ID			getRosters
//		@Router		/roster [get]
func (h *Handler) GetRosters(c *gin.Context) {
	var params FilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	rosters, err := h.rosterService.GetRosters(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, rosters)
}

// GetRoster
//
//	@Summary	Get a specific roster by id
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//
//	@Param		id	path		uint	true	"Roster ID"
//
//	@Success	200	{object}	models.Roster
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@ID			getRoster
//	@Router		/roster/{id} [get]
func (h *Handler) GetRoster(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	params := &FilterParams{
		ID: &id,
	}
	roster, err := h.rosterService.GetRosters(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roster)
}

// UpdateRoster
//
//	@Summary	Update a roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id				path		uint						true	"Roster ID"
//	@Param		updateParams	body		RosterUpdateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			updateRoster
//	@Router		/roster/{id} [patch]
func (h *Handler) UpdateRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	var updateParams UpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedRoster, err := h.rosterService.UpdateRoster(uint(id), &updateParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRoster)
}

// DeleteRoster
//
//	@Summary	DeleteRoster a roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster ID"
//	@Success	200	{string}	string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@ID			deleteRoster
//	@Router		/roster/{id} [delete]
func (h *Handler) DeleteRoster(c *gin.Context) {
	rosterId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	if err := h.rosterService.DeleteRoster(uint(rosterId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roster deleted",
	})
}

func requireRosterOrganRoleQuery(db *gorm.DB, queryStr string, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Query(queryStr)
		if val == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": queryStr + " is required"})
			return
		}
		checkAccess(c, db, val, minRole)
	}
}

// requireRosterOrganRoleParam validates the existence of a roster by its ID
// from the URL parameters and ensures the current user has the required
// minimum role within that roster's organization.
func requireRosterOrganRoleParam(db *gorm.DB, paramStr string, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		rosterID := c.Param(paramStr)
		if rosterID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": paramStr + " is required"})
			return
		}

		var roster models.Roster
		if err := db.First(&roster, "id = ?", rosterID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}

		checkAccess(c, db, roster.OrganID, minRole)
	}
}

func requireRosterOrganRoleBody(db *gorm.DB, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			OrganID uint `json:"organId"`
		}
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil || body.OrganID == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Valid organId is required in body"})
			return
		}
		checkAccess(c, db, body.OrganID, minRole)
	}
}

func requireShiftGroupOrganRoleParams(db *gorm.DB, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		if groupID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
			return
		}

		var shiftGroup models.ShiftGroup
		if err := db.First(&shiftGroup, "id = ?", groupID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Shift group not found"})
			return
		}

		checkAccess(c, db, shiftGroup.OrganID, minRole)
	}
}
