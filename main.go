package main

import (
	"log"
	"ml-master-data/config"
	"ml-master-data/routes"
	"ml-master-data/seeders"

	"github.com/joho/godotenv"
)

// @title ML Master Data API
// @version 1.0
// @description API for ML Master Data
// @host localhost:8080
// @BasePath /api/

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()

	seeders.Seed()

	r := routes.SetupRouter()
	r.Static("/public", "./public")
	r.Run(":8080")
}
