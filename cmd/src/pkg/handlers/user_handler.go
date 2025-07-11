package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"GEWIS-Rooster/cmd/src/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(rg *gin.RouterGroup, userService services.UserServiceInterface) *UserHandler {
	h := &UserHandler{userService: userService}

	g := rg.Group("/user")

	g.POST("/create", h.Create)
	g.GET("/", h.GetAllUsers)
	g.GET("/:id", h.GetUserByID)
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

// GetAllUsers
//
//		@Summary      Get all users with optional filtering
//		@Description  Retrieve a list of users with optional query parameter filtering
//		@Security     BearerAuth
//		@Tags         User
//		@Accept       json
//		@Produce      json
//	 @Param        organId    query     uint    false  "Organ ID"
//	 @Param        gewisId    query     uint    false  "GEWIS ID"
//		@Success      200         {array}   models.User
//		@Failure      400         {object}  map[string]string
//		@Router       /user/ [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	f := models.UserFilterParams{}
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.userService.GetUsers(&f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID
//
//	@Summary      Get user by ID
//	@Description  Retrieve a specific user by their unique ID
//	@Security     BearerAuth
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        id          path      uint    true   "User ID"
//	@Success      200         {object}  models.User
//	@Failure      400         {object}  map[string]string
//	@Failure      404         {object}  map[string]string
//	@Router       /user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	idUint := uint(id)
	filters := &models.UserFilterParams{ID: &idUint}

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
