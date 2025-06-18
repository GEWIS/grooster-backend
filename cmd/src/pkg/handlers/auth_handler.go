package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/middleware"
	"GEWIS-Rooster/cmd/src/pkg/services"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
)

type AuthHandler struct {
	config     *oauth2.Config
	provider   *oidc.Provider
	service    services.AuthServiceInterface
	authMiddle *middleware.AuthMiddleware
}

func NewAuthHandler(rg *gin.RouterGroup, auth services.AuthServiceInterface, authMiddle *middleware.AuthMiddleware) *AuthHandler {
	provider, config := authMiddle.SetupOIDC()

	h := &AuthHandler{config: config, provider: provider, service: auth}

	log.Printf("Path %s", rg.BasePath())

	rg.GET("/redirect", h.AuthRedirect)
	rg.GET("/callback", h.AuthCallback)

	return h
}

// AuthRedirect
//
//	@Summary		Redirect to OIDC provider
//	@Description	Generates state, sets a cookie, and redirects to Google OIDC
//	@Security		BasicAuth
//	@Tags			Auth
//	@Param			state	query		string				true	"State returned from provider"
//	@Success		200		{string}	string				"redirect"
//	@Failure		500		{object}	map[string]string	"pkg server error"
//	@Router			/auth/redirect [get]
func (h *AuthHandler) AuthRedirect(c *gin.Context) {
	state, err := h.service.RandString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.service.SetCallBackCookie(c, state)

	redirectURL := h.config.AuthCodeURL(state)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// AuthCallback
//
//	@Summary		Handle OAuth2 Callback
//	@Description	Validates state, exchanges code for token, and returns user info
//	@Security		BasicAuth
//	@Tags			Auth
//	@Param			state	query		string				true	"State returned from provider"
//	@Param			code	query		string				true	"Authorization code from provider"
//	@Success		200		{object}	map[string]string	"User info and token"
//	@Failure		400		{object}	map[string]string	"Bad request: missing or invalid state"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/auth/callback [get]
func (h *AuthHandler) AuthCallback(c *gin.Context) {
	state, err := c.Cookie("auth_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state not found"})
		return
	}

	if c.Query("state") != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state did not match"})
		return
	}

	oauth2Token, err := h.config.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange token"})
	}

	h.service.ProcessUserInfo(oauth2Token)

	redirectUrl := fmt.Sprintf("http://localhost:5173/callback?token=%s", oauth2Token.AccessToken)
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
