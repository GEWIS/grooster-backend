package main

import (
	"GEWIS-Rooster/cmd/src/docs"
	"GEWIS-Rooster/cmd/src/pkg"
	"GEWIS-Rooster/cmd/src/pkg/handlers"
	"GEWIS-Rooster/cmd/src/pkg/middleware"
	"GEWIS-Rooster/cmd/src/pkg/services"
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
// @description				A GEWIS Rooster maker
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	log.Print("Starting server")

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msgf("Error loading .env file %s", err.Error())
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
		AllowOrigins:     []string{os.Getenv("ALLOWED_ORIGINS")}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // If you need to support cookies or authentication
	}))

	api := r.Group(os.Getenv("BASE_PATH"))

	userService := services.NewUserService(db)
	rosterService := services.NewRosterService(db)
	authService := services.NewAuthService(userService, db)
	authMiddle := middleware.NewAuthMiddleware(authService)

	// Auth routes (no authentication required)
	authGroup := api.Group("/auth")
	handlers.NewAuthHandler(authGroup, authService, authMiddle)

	protectedGroup := api.Group("")
	protectedGroup.Use(authMiddle.AuthMiddlewareCheck())
	{
		handlers.NewUserHandler(protectedGroup, userService)
		handlers.NewRosterHandler(rosterService, protectedGroup)
	}

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Print("Server error", err)
	}
}
