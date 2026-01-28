package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type AuthServiceInterface interface {
	SetCallBackCookie(*gin.Context, string)
	RandString(int) (string, error)
	ProcessUserInfo(*oauth2.Token) (string, error)
	GetOrgans(claims map[string]interface{}) ([]models.Organ, error)
	HandleLocalAuthentication(ctx *gin.Context) (string, error)
	CreateInternalToken(user *models.User) (string, error)
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

func (s *AuthService) HandleLocalAuthentication(ctx *gin.Context) (string, error) {
	var user models.User
	if err := s.db.First(&user).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "No user in local database found, try to seed it first"})
		return "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.GEWISID,
		"name": user.Name,
		"iat":  now.Unix(),
		"exp":  now.Add(7 * 24 * time.Hour).Unix(), // 7 days expiry
	}

	secret := os.Getenv("JWT_SECRET")
	if strings.TrimSpace(secret) == "" {
		err := fmt.Errorf("JWT_SECRET environment variable is not set or is empty")
		log.Error().Err(err).Msg("failed to create internal token due to missing JWT secret")
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func (s *AuthService) ProcessUserInfo(OAuth2Token *oauth2.Token) (string, error) {
	token := OAuth2Token.AccessToken
	infoString := strings.Split(token, ".")[1]

	// Decode the payload (second part)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(infoString)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	// Unmarshal JSON into a map
	var claims map[string]interface{}
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	log.Debug().
		Interface("preferred_username", claims["preferred_username"]).
		Interface("given_name", claims["given_name"]).
		Msg("Processing user claims")

	// Extract the id
	idParts := strings.Split(claims["preferred_username"].(string), "m")
	if len(idParts) < 2 {
		log.Error().Msg("Failed to parse preferred_username: missing 'm' prefix")
		return "", errors.New("invalid username format")
	}

	idStr := idParts[1]
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Str("idStr", idStr).Msg("Failed to convert ID to int")
		return "", err // Added return to prevent nil pointer later
	}
	username := claims["given_name"].(string)

	id := uint(idInt)
	filters := &models.UserFilterParams{
		ID: &id,
	}

	// 2. Log the search attempt
	log.Info().Uint("search_id", id).Msg("Searching for existing user")

	users, err := s.u.GetUsers(filters)
	var user *models.User

	if err != nil {
		log.Error().Err(err).Uint("id", id).Msg("Database error during GetUsers")
		return "", err
	}

	if len(users) > 0 {
		user = users[0]
		log.Info().Uint("user_id", user.ID).Str("name", user.Name).Msg("Existing user found")
	}

	if user == nil {
		log.Info().Uint("gewis_id", id).Msg("User not found, attempting to create new record")

		organs, organErr := s.GetOrgans(claims)
		if organErr != nil {
			log.Error().Err(organErr).Msg("Failed to get organs for new user")
			return "", organErr
		}

		params := models.UserCreateRequest{
			Name:    username,
			GEWISID: id,
			Organs:  organs,
		}

		// 3. Log the creation parameters to see exactly what is being sent to GORM
		log.Debug().Interface("create_params", params).Msg("Sending Create request to user service")

		user, err = s.u.Create(&params)
		if err != nil {
			log.Error().Err(err).Uint("attempted_id", id).Msg("User creation failed")
			return "", err
		}
		log.Info().Uint("new_user_id", user.ID).Msg("Successfully created new user")
	} else {
		organs, organErr := s.GetOrgans(claims)
		if organErr != nil {
			log.Error().Msg("Failed to get organs: " + organErr.Error())
			return "", organErr
		}

		if err := s.db.Model(user).Association("Organs").Replace(organs); err != nil {
			log.Error().Msg("Failed to update user organs: " + err.Error())
		}
	}

	jwtToken, err := s.CreateInternalToken(user)

	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	return jwtToken, nil
}

func (s *AuthService) GetOrgans(claims map[string]interface{}) ([]models.Organ, error) {
	resourceAccess, ok := claims["resource_access"].(map[string]interface{})
	if !ok {
		log.Debug().Interface("available_claims", claims).Msg("resource_access missing")
		return nil, fmt.Errorf("resource_access not found or wrong type")
	}

	envType := os.Getenv("RESOURCE_ENV_TYPE")
	key := fmt.Sprintf("grooster-%s", envType)

	OICDName, ok := resourceAccess[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("grooster-test not found or wrong type")
	}

	roles, ok := OICDName["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("roles not found or wrong type")
	}

	var organs []models.Organ

	separator := envType + " "

	for _, role := range roles {
		if roleStr, ok := role.(string); ok {
			if strings.Contains(roleStr, separator) {

				parts := strings.SplitN(roleStr, separator, 2)

				if len(parts) > 1 && parts[1] != "" {
					organString := parts[1]

					organ := models.Organ{
						Name: organString,
					}
					s.db.FirstOrCreate(&organ, models.Organ{Name: organString})
					organs = append(organs, organ)
				}
			}
		}
	}

	return organs, nil
}

func (s *AuthService) CreateInternalToken(user *models.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.GEWISID,
		"name": user.Name,
		"iat":  now.Unix(),
		"exp":  now.Add(7 * 24 * time.Hour).Unix(), // 7 days expiry
	}

	secret := os.Getenv("JWT_SECRET")
	if strings.TrimSpace(secret) == "" {
		log.Error().Msg("JWT_SECRET environment variable is not set or empty")
		return "", fmt.Errorf("JWT_SECRET environment variable is not set or empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
