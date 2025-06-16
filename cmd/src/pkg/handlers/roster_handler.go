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

type RosterHandler struct {
	rosterService services.RosterServiceInterface
}

func NewRosterHandler(rosterService services.RosterServiceInterface, rg *gin.RouterGroup) *RosterHandler {
	h := &RosterHandler{rosterService: rosterService}

	g := rg.Group("/roster")

	g.POST("/", h.CreateRoster)
	g.GET("/", h.GetRosters)
	g.GET("/:id", h.GetOrganRosters)
	g.PATCH("/:id", h.UpdateRoster)
	g.DELETE("/:id", h.DeleteRoster)

	g.POST("/shift", h.CreateRosterShift)
	g.DELETE("/shift/:id", h.DeleteRosterShift)
	g.POST("/answer", h.CreateRosterAnswer)
	g.PATCH("/answer/:id", h.UpdateRosterAnswer)

	g.POST("/:id/save", h.SaveRoster)
	g.PATCH("/saved-shift/:id", h.UpdateSavedShift)
	g.GET("/saved-shift/:id", h.GetSavedRoster)

	return h
}

// CreateRoster
//
//	@Summary	CreateRoster a new roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		models.RosterCreateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			createRoster
//	@Router		/roster/ [post]
func (h *RosterHandler) CreateRoster(c *gin.Context) {
	var param *models.RosterCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
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
//	@Summary	Get all rosters
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		afterDate	query		string	false	"Roster after this date"
//	@Success	200			{array}		models.RosterResponse
//	@Failure	400			{string}	string
//	@Failure	404			{string}	string
//	@ID			getRosters
//	@Router		/roster/ [get]
func (h *RosterHandler) GetRosters(c *gin.Context) {
	date := c.Query("afterDate")

	rosters, err := h.rosterService.GetRosters(&date)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rosters)
}

// GetOrganRosters
//
//	@Summary	Get the organs rosters
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		uint	true	"Organ ID"
//	@Success	200	{array}	models.RosterResponse
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@ID			getOrganRosters
//	@Router		/roster/{id} [get]
func (h *RosterHandler) GetOrganRosters(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	organsContext, ext := c.Get("organs")
	if !ext {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find any organs"})
		return
	}

	organs, ok := organsContext.([]*models.Organ)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get organs from context"})
		return
	}

	found := false
	for _, organ := range organs {
		if uint64(organ.ID) == id {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not in organ"})
		return
	}

	rosters, err := h.rosterService.GetOrganRosters(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, rosters)
}

// UpdateRoster
//
//	@Summary	Update a roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id				path		uint						true	"Roster ID"
//	@Param		updateParams	body		models.RosterUpdateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			updateRoster
//	@Router		/roster/{id} [patch]
func (h *RosterHandler) UpdateRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	var updateParams models.RosterUpdateRequest
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
func (h *RosterHandler) DeleteRoster(c *gin.Context) {
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

// CreateRosterShift
//
//	@Summary	Create a new roster shift
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		models.RosterShiftCreateRequest	true	"Roster shift input"
//	@Success	200				{object}	models.RosterShift
//	@Failure	400				{string}	string
//	@ID			createRosterShift
//	@Router		/roster/shift [post]
func (h *RosterHandler) CreateRosterShift(c *gin.Context) {
	var param *models.RosterShiftCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdRosterShift, err := h.rosterService.CreateRosterShift(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdRosterShift)
}

// DeleteRosterShift
//
//	@Summary	Deletes a roster shift
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster Answer ID"
//	@Success	200	{string}	string
//	@Failure	400	{string}	string
//	@ID			deleteRosterShift
//	@Router		/roster/shift/{id} [delete]
func (h *RosterHandler) DeleteRosterShift(c *gin.Context) {
	rosterId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	if err := h.rosterService.DeleteRosterShift(uint(rosterId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roster deleted",
	})
}

// CreateRosterAnswer
//
//	@Summary	Create a new roster shift answer
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		models.RosterAnswerCreateRequest	true	"Roster answer input"
//	@Success	200				{object}	models.RosterAnswer
//	@Failure	400				{string}	string
//	@ID			createRosterAnswer
//	@Router		/roster/answer [post]
func (h *RosterHandler) CreateRosterAnswer(c *gin.Context) {
	var param *models.RosterAnswerCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdAnswer, err := h.rosterService.CreateRosterAnswer(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAnswer)
}

// UpdateRosterAnswer
//
//	@Summary	Updates a roster answer with the new value
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id				path		int									true	"Roster Answer ID"
//	@Param		updateParams	body		models.RosterAnswerUpdateRequest	true	"New answer value"
//	@Success	200				{object}	models.RosterAnswer
//	@Failure	400				{string}	json
//	@Failure	404				{string}	json
//	@ID			updateRosterAnswer
//	@Router		/roster/answer/{id} [patch]
func (h *RosterHandler) UpdateRosterAnswer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster answer ID"})
		return
	}

	var updateParams models.RosterAnswerUpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedAnswer, err := h.rosterService.UpdateRosterAnswer(uint(id), &updateParams)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAnswer)
}

// SaveRoster
//
//	@Summary	Save a specific roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster ID"
//	@Success	200	{string}	json
//	@Failure	400	{string}	json
//	@Failure	404	{string}	json
//	@ID			rosterSave
//	@Router		/roster/{id}/save [post]
func (h *RosterHandler) SaveRoster(c *gin.Context) {
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
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id				path		int								true	"SavedShift ID"
//	@Param		updateParams	body		models.SavedShiftUpdateRequest	true	"Update data"
//	@Success	200				{object}	models.SavedShift
//	@Failure	400				{string}	string	"Invalid request"
//	@Failure	404				{string}	string	"SavedShift not found"
//	@ID			updateSavedShift
//	@Router		/roster/saved-shift/{id} [patch]
func (h *RosterHandler) UpdateSavedShift(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saved shift ID"})
		return
	}

	var updateParams models.SavedShiftUpdateRequest
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
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Roster ID"
//	@Success	200	{array}		models.SavedShift	"Saved Shifts"
//	@Failure	400	{string}	string				"Invalid request"
//	@Failure	404	{string}	string				"SavedShift not found"
//	@ID			getSavedRoster
//	@Router		/roster/saved-shift/{id} [get]
func (h *RosterHandler) GetSavedRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	savedShifts, err := h.rosterService.GetSavedRoster(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, savedShifts)
}
