package config

import (
	"fmt"
	"log"
	"ml-master-data/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	fmt.Println("Connected to database successfully")

	err = database.AutoMigrate(
		&models.HeroPick{},
		&models.HeroBan{},
		&models.ObjectiveResult{},
		&models.PlayerStats{},
		&models.CoachStats{},
		&models.PriorityPick{},
		&models.PriorityBan{},
		&models.FlexPick{},
		&models.MatchVideo{},
		&models.Coach{},
		&models.GameDetail{},
		&models.Hero{},
		&models.Team{},
		&models.Player{},
		&models.Game{},
		&models.Match{},
		&models.TournamentTeam{},
		&models.Tournament{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")

	DB = database
}
