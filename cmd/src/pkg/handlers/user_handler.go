package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"GEWIS-Rooster/cmd/src/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(rg *gin.RouterGroup, userService services.UserServiceInterface) *UserHandler {
	h := &UserHandler{userService: userService}

	g := rg.Group("/user")

	log.Printf("Path %s", g.BasePath())

	g.POST("/create", h.Create)
	g.GET("/", h.GetUsers)
	g.GET("/:id", h.GetUsers)
	g.DELETE("/:id", h.Delete)

	return h
}

// Create
//
//	@Summary		CreateRoster a new user
//	@Security		BearerAuth
//	@Description	create user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			createParams	body		models.UserCreateRequest	true	"User input"
//	@Success		200				{object}	models.User
//	@Failure		400				{string}	string
//	@Router			/user/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	var param *models.UserCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdUser, err := h.userService.Create(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": createdUser,
	})
}

// GetUsers
//
//	@Summary	Get users optionally filtered by parameters
//	@Security	BearerAuth
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id			path		uint	true	"ID"
//	@Param		filter	query		models.UserFilterParams	false	"Filter parameters"
//	@Success	200			{array}		models.User
//	@Failure	400			{string}	string
//	@Router		/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	idParam := c.Param("id")

	var filters *models.UserFilterParams

	if idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
			return
		}
		idUint := uint(id)
		filters = &models.UserFilterParams{ID: &idUint}
	} else {
		f := models.UserFilterParams{}
		if err := c.ShouldBindQuery(&f); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if f.ID == nil && f.OrganID == nil && f.GEWISID == nil {
			filters = &f
		}
	}

	users, err := h.userService.GetUsers(filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Delete
//
//	@Summary	DeleteRoster a user
//	@Security	BearerAuth
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"User ID"
//	@Success	200	{string}	string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/user/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.Delete(uint(userId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
	})
}
