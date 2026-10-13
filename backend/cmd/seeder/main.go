package main

import (
	"GEWIS-Rooster/cmd/seeder/seeder"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		log.Print("Loading .env file")
		if err := godotenv.Load(); err != nil {
			log.Fatal().Msgf("Error loading .env file: %v", err)
		}
	}

	db := os.Getenv("DATABASE")

	seeder.Seeder(db)
}
