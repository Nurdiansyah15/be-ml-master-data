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
		&models.User{},
		&models.Tournament{},
		&models.Match{},
		&models.PlayerMatch{},
		&models.CoachMatch{},
		&models.Game{},
		&models.Team{},
		&models.Player{},
		&models.Coach{},
		&models.Hero{},
		&models.MatchTeamDetail{},
		&models.HeroPick{},
		&models.HeroBan{},
		&models.HeroPickGame{},
		&models.HeroBanGame{},
		&models.PriorityPick{},
		&models.PriorityBan{},
		&models.FlexPick{},
		&models.GameResult{},
		&models.TrioMid{},
		&models.Goldlaner{},
		&models.Explaner{},
		&models.TurtleResult{},
		&models.LordResult{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")

	DB = database
}
