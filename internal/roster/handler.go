package roster

import (
	"GEWIS-Rooster/internal/models"
	_ "GEWIS-Rooster/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	rosterService Service
}

func NewRosterHandler(rosterService Service, rg *gin.RouterGroup, db *gorm.DB) *Handler {
	h := &Handler{rosterService: rosterService}
	g := rg.Group("/roster")

	h.registerRosterRoutes(g, db)
	h.registerShiftRoutes(g, db)
	h.registerTemplateRoutes(g, db)

	g.POST("/:id/fill", requireRosterOrganRoleParam(db, "id", models.RoleAdmin), h.FillRosterPreferences)

	g.POST("/:id/save", requireRosterOrganRoleParam(db, "id", models.RoleAdmin), h.SaveRoster)
	g.PATCH("/saved-shift/:id", h.UpdateSavedShift)
	g.GET("/saved-shift/:id", h.GetSavedRoster)

	g.POST("/shift-groups", requireRosterOrganRoleBody(db, models.RoleAdmin), h.CreateShiftGroup)
	g.GET("/shift-groups", h.GetShiftGroups)
	g.GET("/shift-groups/:id", h.GetShiftGroup)

	g.PUT("/shift-groups/:id/priority", requireShiftGroupOrganRoleParams(db, models.RoleAdmin))

	return h
}

// FillRosterPreferences
//
//	@Summary	Fills a roster with the linked user template preferences
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster ID"
//	@Success	200	{array}	models.RosterAnswer
//	@Failure	400	{string}	string
//	@Failure 404 {string} string
//	@ID			fillRoster
//	@Router		/roster/{id}/fill [post]
func (h *Handler) FillRosterPreferences(c *gin.Context) {
	rosterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	answers, err := h.rosterService.FillRosterPreferences(uint(rosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, answers)
}

// SaveRoster
//
//	@Summary	Save a specific roster
//	@Security	BearerAuth
//	@Tags		Saved Shift
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster ID"
//	@Success	200	{string}	json
//	@Failure	400	{string}	json
//	@Failure	404	{string}	json
//	@ID			rosterSave
//	@Router		/roster/{id}/save [post]
func (h *Handler) SaveRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	err = h.rosterService.SaveRoster(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SavedShift not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shift"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// UpdateSavedShift
//
//	@Summary	Update a specific saved shift
//	@Security	BearerAuth
//	@Tags		Saved Shift
//	@Accept		json
//	@Produce	json
//	@Param		id				path		int								true	"SavedShift ID"
//	@Param		updateParams	body		SavedShiftUpdateRequest	true	"Update data"
//	@Success	200				{object}	models.SavedShift
//	@Failure	400				{string}	string	"Invalid request"
//	@Failure	404				{string}	string	"SavedShift not found"
//	@ID			updateSavedShift
//	@Router		/roster/saved-shift/{id} [patch]
func (h *Handler) UpdateSavedShift(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saved shift ID"})
		return
	}

	var updateParams SavedShiftUpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	updatedShift, err := h.rosterService.UpdateSavedShift(uint(id), &updateParams)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SavedShift not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shift"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedShift)
}

// GetSavedRoster
//
//	@Summary	Get all saved shifts for a specific roster
//	@Security	BearerAuth
//	@Tags		Saved Shift
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Roster ID"
//	@Success	200	{object}		SavedShiftResponse	"Saved Shift Response"
//	@Failure	400	{string}	string				"Invalid request"
//	@Failure	404	{string}	string				"SavedShift not found"
//	@ID			getSavedRoster
//	@Router		/roster/saved-shift/{id} [get]
func (h *Handler) GetSavedRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	savedShifts, savedShiftOrdering, err := h.rosterService.GetSavedRoster(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := SavedShiftResponse{
		SavedShifts:        savedShifts,
		SavedShiftOrdering: savedShiftOrdering,
	}
	// Log the entire struct as a field called "response"
	log.Debug().Interface("response", response).Msg("Sending saved roster response")

	c.JSON(http.StatusOK, response)
}

// CreateShiftGroup
//
//	@Summary   Create a new shift group
//	@Security  BearerAuth
//	@Tags      ShiftGroup
//	@Accept    json
//	@Produce   json
//	@Param     params body      ShiftGroupCreateRequest  true "Shift Group Details"
//	@Success   201    {object}   models.ShiftGroup
//	@Failure   400    {string}   string
//	@ID        createShiftGroup
//	@Router    /roster/shift-groups [post]
func (h *Handler) CreateShiftGroup(c *gin.Context) {
	var params ShiftGroupCreateRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.rosterService.CreateShiftGroup(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// GetShiftGroups
//
//	@Summary   Get all shift groups for an organ
//	@Security  BearerAuth
//	@Tags      ShiftGroup
//	@Produce   json
//	@Param     organ_id query    int  true "Organ ID"
//	@Success   200      {array}   models.ShiftGroup
//	@Failure   400      {string}  string
//	@ID        getShiftGroups
//	@Router    /roster/shift-groups [get]
func (h *Handler) GetShiftGroups(c *gin.Context) {
	var filters ShiftGroupFilterParams

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.rosterService.GetShiftGroups(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetShiftGroup
//
//	@Summary   Get a specific shift group by ID
//	@Security  BearerAuth
//	@Tags      ShiftGroup
//	@Produce   json
//	@Param     id path      int  true "Shift Group ID"
//	@Success   200 {object} models.ShiftGroup
//	@Failure   404 {string}   string
//	@ID        getShiftGroup
//	@Router    /roster/shift-groups/{id} [get]
func (h *Handler) GetShiftGroup(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	group, err := h.rosterService.GetShiftGroup(uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shift group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// UpdateShiftGroupPriority
//
//	@Summary	Update a shift group priority
//	@Security	BearerAuth
//	@Tags		ShiftGroup
//	@Accept		json
//	@Produce	json
//	@Param		id				path		int								true	"ShiftGroup ID"
//	@Param		updateParams	body		GroupUpdatePriorityParam	true	"Update parameters"
//	@Success	200				{object}	models.ShiftGroupPriority
//	@Failure	400				{string}	string	"Invalid request"
//	@Failure	404				{string}	string	"SavedShift not found"
//	@ID			updateShiftGroupPriority
//	@Router    /roster/shift-groups/{id}/priority [put]
func (h *Handler) UpdateShiftGroupPriority(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var params GroupUpdatePriorityParam
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupPriority, err := h.rosterService.UpdateShiftGroupPriority(uint(id), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groupPriority)
}
