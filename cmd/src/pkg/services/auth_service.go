package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"io"
	"strings"
)

type AuthServiceInterface interface {
	SetCallBackCookie(*gin.Context, string)
	RandString(int) (string, error)
	ProcessUserInfo(*oauth2.Token)
}

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) SetCallBackCookie(c *gin.Context, value string) {
	secure := c.Request.TLS != nil // true if HTTPS

	c.SetCookie(
		"auth_state",
		value,
		3600,
		"/", // path
		"",  // domain (default)
		secure,
		true, // HttpOnly
	)
}

func (s *AuthService) RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *AuthService) ProcessUserInfo(OAuth2Token *oauth2.Token) {
	token := OAuth2Token.AccessToken
	infoString := strings.Split(token, ".")[1]

	// Decode the payload (second part)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(infoString)
	if err != nil {
		log.Error().Msgf(err.Error())
		return
	}

	// Unmarshal JSON into a map
	var claims map[string]interface{}
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	// Extract the claim
	gewisId, ok := claims["preferred_username"].(string)
	if !ok {
		log.Error().Msg("Error getting gewisId")
		return
	}

	log.Print(gewisId)
}
