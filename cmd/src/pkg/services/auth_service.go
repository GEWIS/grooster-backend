package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"io"
	"strconv"
	"strings"
)

type AuthServiceInterface interface {
	SetCallBackCookie(*gin.Context, string)
	RandString(int) (string, error)
	ProcessUserInfo(*oauth2.Token)
	GetOrgans(claims map[string]interface{}) ([]*models.Organ, error)
}

type AuthService struct {
	u  *UserService
	db *gorm.DB
}

func NewAuthService(u *UserService, db *gorm.DB) *AuthService {
	return &AuthService{u, db}
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

	// Extract the id
	idStr := strings.Split(claims["preferred_username"].(string), "m")[1]
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	username := claims["given_name"].(string)

	id := uint(idInt)
	filters := &models.UserFilterParams{
		ID: &id,
	}

	users, err := s.u.GetUsers(filters)
	var user *models.User

	if err != nil {
		log.Error().Msg("Failed to get users: " + err.Error())
		return
	}

	if len(users) > 0 {
		user = users[0]
	}

	if user == nil {
		id := uint(idInt)

		organs, organErr := s.GetOrgans(claims)
		if organErr != nil {
			log.Error().Msg("Failed to get organs: " + organErr.Error())
			return
		}

		params := models.UserCreateRequest{
			Name:    username,
			GEWISID: id,
			Organs:  organs,
		}

		user, err = s.u.Create(&params)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
	} else {
		organs, organErr := s.GetOrgans(claims)
		if organErr != nil {
			log.Error().Msg("Failed to get organs: " + organErr.Error())
			return
		}

		if err := s.db.Model(user).Association("Organs").Replace(organs); err != nil {
			log.Error().Msg("Failed to update user organs: " + err.Error())
		}
	}
}

func (s *AuthService) GetOrgans(claims map[string]interface{}) ([]models.Organ, error) {
	resourceAccess, ok := claims["resource_access"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("resource_access not found or wrong type")
	}

	OICDName, ok := resourceAccess["grooster-test"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("grooster-test not found or wrong type")
	}

	roles, ok := OICDName["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("roles not found or wrong type")
	}

	var organs []models.Organ

	for _, role := range roles {
		if roleStr, ok := role.(string); ok {
			if strings.HasPrefix(roleStr, "PRIV") {
				var organString = strings.TrimPrefix(roleStr, "PRIV - ")

				organ := models.Organ{
					Name: organString,
				}
				s.db.FirstOrCreate(&organ, models.Organ{Name: organString})

				organs = append(organs, organ)
			}
		}
	}

	return organs, nil
}
