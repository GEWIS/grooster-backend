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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// @title						GRooster
// @version					0.1
// @description				A GEWIS Rooster maker for fun
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	log.Print("Starting server")

	if _, err := os.Stat(".env"); err == nil {
		log.Print("Loading .env file")
		if err := godotenv.Load(); err != nil {
			log.Fatal().Msgf("Error loading .env file: %v", err)
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
			log.Print("failed to close database connection", err)
		}
	}(sqlDB)

	log.Print("Connecting to database")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // If you need to support cookies or authentication
	}))

	api := r.Group(os.Getenv("BASE_PATH"))

	userService := user.NewUserService(db)
	rosterService := roster.NewRosterService(db)
	exportService := export.NewExportService(rosterService, db)
	organService := organ.NewOrganService(db)

	m := middleware.AuthMiddleware{}
	provider, config := m.SetupOIDC()

	authService := auth.NewAuthService(userService, db)
	authMiddle := middleware.NewAuthMiddleware(authService)

	// Auth routes (no authentication required)
	authGroup := api.Group("/auth")
	auth.NewAuthHandler(authGroup, authService, provider, config)

	protectedGroup := api.Group("")
	protectedGroup.Use(authMiddle.AuthMiddlewareCheck())
	{
		user.NewUserHandler(protectedGroup, userService)
		roster.NewRosterHandler(rosterService, protectedGroup)
		export.NewExportHandler(exportService, protectedGroup)
		organ.NewOrganHandler(protectedGroup, organService)
	}

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Print("Server error", err)
	}
}
