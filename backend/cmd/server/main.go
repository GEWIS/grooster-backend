package main

import (
	"GEWIS-Rooster/docs"
	"GEWIS-Rooster/internal/auth"
	"GEWIS-Rooster/internal/export"
	"GEWIS-Rooster/internal/organ"
	"GEWIS-Rooster/internal/platform/database"
	"GEWIS-Rooster/internal/platform/middleware"
	"GEWIS-Rooster/internal/roster"
	"GEWIS-Rooster/internal/user"
	"database/sql"
	"log/slog"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title						GRooster
// @version					0.1
// @description				A GEWIS Rooster maker for fun
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	if _, err := os.Stat(".env"); err == nil {
		logger.Info("Loading .env file")
		if err := godotenv.Load(); err != nil {
			logger.Error("Error loading .env file", "err", err)
			os.Exit(1)
		}
	}

	db := database.ConnectDB(os.Getenv("DATABASE"))
	sqlDB, _ := db.DB()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	docs.SwaggerInfo.Title = "Docs for grooster"
	docs.SwaggerInfo.BasePath = os.Getenv("BASE_PATH")

	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			logger.Error("failed to close database connection", "err", err)
		}
	}(sqlDB)

	logger.Info("Connecting to database")

	r := gin.Default()

	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if rawOrigins != "" {
		origins = strings.Split(rawOrigins, ",")
	}

	if len(origins) == 0 {
		logger.Error("No ORIGINS provided")
		os.Exit(1)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // If you need to support cookies or authentication
	}))

	api := r.Group(os.Getenv("BASE_PATH"))

	userService := user.NewUserService(db)
	organService := organ.NewOrganService(db)
	rosterService := roster.NewRosterService(db, userService)
	exportService := export.NewExportService(rosterService, organService, db)

	m := middleware.AuthMiddleware{}
	provider, config := m.SetupOIDC()

	logMiddleware := middleware.NewLoggingMiddleware()
	api.Use(logMiddleware.LogMiddleware())

	authVerifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("CLIENT_ID"), SkipClientIDCheck: true})
	authService := auth.NewAuthService(userService, db, authVerifier)

	authMiddle := middleware.NewAuthMiddleware(authService, userService)

	// Auth routes (no authentication required)
	authGroup := api.Group("/auth")
	auth.NewAuthHandler(authGroup, authService, provider, config)

	protectedGroup := api.Group("")
	protectedGroup.Use(authMiddle.AuthMiddlewareCheck())
	{
		user.NewUserHandler(protectedGroup, userService)
		roster.NewRosterHandler(rosterService, protectedGroup, db)
		export.NewExportHandler(exportService, protectedGroup)
		organ.NewOrganHandler(protectedGroup, organService, db)
	}

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		logger.Error("Server error", "err", err)
	}
}
