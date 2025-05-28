package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/services"
	"context"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type AuthHandler struct {
	config   *oauth2.Config
	provider *oidc.Provider
	service  services.AuthServiceInterface
}

func NewAuthHandler(rg *gin.RouterGroup, auth services.AuthServiceInterface) *AuthHandler {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://auth.gewis.nl/realms/GEWISWG")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create OIDC provider")
		return nil
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("URI_CALLBACK"),
		Endpoint:     provider.Endpoint(),
		// TODO Add correct GEWIS scopes
		Scopes: []string{oidc.ScopeOpenID},
	}

	h := &AuthHandler{config: config, provider: provider, service: auth}

	g := rg.Group("/auth")

	log.Printf("Path %s", g.BasePath())

	g.GET("/redirect", h.AuthRedirect)
	g.GET("/callback", h.AuthCallback)

	return h
}

// AuthRedirect
//
//	@Summary		Redirect to OIDC provider
//	@Description	Generates state, sets a cookie, and redirects to Google OIDC
//	@Tags			Auth
//	@Success		302	{string}	string				"redirect"
//	@Failure		500	{object}	map[string]string	"pkg server error"
//	@Router			/auth/redirect [get]
func (h *AuthHandler) AuthRedirect(c *gin.Context) {
	state, err := h.service.RandString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.service.SetCallBackCookie(c, state)

	redirectURL := h.config.AuthCodeURL(state)
	c.JSON(http.StatusOK, gin.H{"state": state, "redirectURL": redirectURL})
}

// AuthCallback
//
//	@Summary		Handle OAuth2 Callback
//	@Description	Validates state, exchanges code for token, and returns user info
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

	userInfo, err := h.provider.UserInfo(c.Request.Context(), oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user info" + err.Error()})
		return
	}

	resp := struct {
		OAuth2Token *oauth2.Token
		UserInfo    *oidc.UserInfo
	}{oauth2Token, userInfo}

	h.service.ProcessUserInfo(oauth2Token)

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": string(data)})
}
