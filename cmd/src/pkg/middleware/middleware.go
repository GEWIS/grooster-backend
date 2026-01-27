package middleware

import (
	"GEWIS-Rooster/cmd/src/pkg/services"
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing Authorization header"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := strings.TrimSpace(authHeader[len(bearerPrefix):])
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Empty bearer token"})
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if strings.TrimSpace(secret) == "" {
			log.Error().Msg("INTERNAL_JWT_SECRET is not set or empty")
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token's signing method is HS256
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

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
