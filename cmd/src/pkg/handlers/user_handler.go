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
	g.GET("/", h.GetAll)
	g.GET("/:id", h.Get)
	g.PATCH("/:id", h.Update)
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
//	@Param			createParams	body		models.UserCreateOrUpdate	true	"User input"
//	@Success		200				{object}	models.User
//	@Failure		400				{string}	string
//	@Router			/user/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	var param *models.UserCreateOrUpdate

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

// GetAll
//
//	@Summary	Get all users
//	@Security	BearerAuth
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		models.User
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/user/ [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Print(users[0])
	c.JSON(http.StatusOK, users)
}

// Get
//
//	@Summary	Get user by GEWIS id
//	@Security	BearerAuth
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id				path		uint						true	"GEWIS ID"
//	@Success	200	{array}		models.User
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/user/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update
//
//	@Summary	UpdateRoster a user
//	@Security	BearerAuth
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id				path		uint						true	"User ID"
//	@Param		updateParams	body		models.UserCreateOrUpdate	true	"User input"
//	@Success	200				{object}	models.User
//	@Failure	400				{string}	string
//	@Router		/user/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateParams models.UserCreateOrUpdate
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedUser, err := h.userService.Update(uint(id), &updateParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": updatedUser,
	})
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
