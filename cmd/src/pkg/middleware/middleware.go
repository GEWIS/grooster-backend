package middleware

import (
	"GEWIS-Rooster/cmd/src/pkg/services"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
)

var (
	verifier *oidc.IDTokenVerifier
)

type AuthMiddlewareInterface interface {
	AuthMiddlewareCheck() gin.HandlerFunc
	SetupOIDC() (*oidc.IDToken, *oidc.Config)
}

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(auth *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: auth}
}

// AuthMiddlewareCheck creates a middleware that validates OIDC tokens
func (a *AuthMiddleware) AuthMiddlewareCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("DEV_TYPE") == "local" {
			token, err := a.authService.HandleLocalAuthentication(c)

			if err != nil {
				return
			}

			c.Header("Authorization", "Bearer "+token)

			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		if verifier == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OIDC verifier not initialized",
			})
			c.Abort()
			return
		}

		ctx := context.Background()
		_, err := verifier.Verify(ctx, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Invalid token: %v", err),
			})
			c.Abort()
			return
		}

		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			log.Error().Msg("invalid token format")
			return
		}

		// parts[1] is the payload
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		// Unmarshal JSON into a map
		var claims map[string]interface{}
		err = json.Unmarshal(payloadBytes, &claims)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		organs, err := a.authService.GetOrgans(claims)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		c.Set("organs", organs)
		c.Next()
	}
}

func (a *AuthMiddleware) SetupOIDC() (*oidc.Provider, *oauth2.Config) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://auth.gewis.nl/realms/GEWISWG")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create OIDC provider")
		return nil, nil
	}

	clientID := os.Getenv("CLIENT_ID")

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("URI_CALLBACK"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	err = setupOIDCVerifier(provider, clientID)
	if err != nil {
		return nil, nil
	}

	return provider, config
}

// setupOIDCVerifier creates and stores the OIDC verifier globally
func setupOIDCVerifier(provider *oidc.Provider, clientID string) error {
	verifier = provider.Verifier(&oidc.Config{
		ClientID:          clientID,
		SkipClientIDCheck: true,
	})

	if verifier == nil {
		return fmt.Errorf("failed to create OIDC verifier")
	}

	return nil
}
