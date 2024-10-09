package main

import (
	"log"
	"ml-master-data/config"
	"ml-master-data/routes"
	"ml-master-data/seeders"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()

	seeders.Seed()

	r := routes.SetupRouter()
	r.Run(":8080")
}
