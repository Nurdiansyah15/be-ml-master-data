package seeders

import (
	"log"
	"ml-master-data/config"
	"ml-master-data/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	// Clear existing data
	// config.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	// config.DB.Unscoped().Delete(&models.User{})
	// config.DB.Unscoped().Delete(&models.Tournament{})
	// config.DB.Unscoped().Delete(&models.Team{})
	// config.DB.Unscoped().Delete(&models.Match{})
	// config.DB.Unscoped().Delete(&models.TournamentTeam{})
	// config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// config.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	// config.DB.Unscoped().Where("1 = 1").Delete(&models.User{})
	// config.DB.Unscoped().Where("1 = 1").Delete(&models.Tournament{})
	// config.DB.Unscoped().Where("1 = 1").Delete(&models.Team{})
	// config.DB.Unscoped().Where("1 = 1").Delete(&models.Match{})
	// config.DB.Unscoped().Where("1 = 1").Delete(&models.TournamentTeam{})
	// config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	config.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	config.DB.Exec("DELETE FROM users")
	// config.DB.Exec("DELETE FROM tournaments")
	// config.DB.Exec("DELETE FROM teams")
	// config.DB.Exec("DELETE FROM matches")
	// config.DB.Exec("DELETE FROM tournament_teams")
	config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Seed Users
	seedUsers()

	log.Println("Seeding completed successfully!")
}

func seedUsers() []models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []models.User{
		{Username: "admin1", Password: string(hashedPassword)},
		{Username: "admin2", Password: string(hashedPassword)},
		{Username: "admin3", Password: string(hashedPassword)},
		{Username: "admin4", Password: string(hashedPassword)},
		{Username: "admin5", Password: string(hashedPassword)},
	}

	for i := range users {
		config.DB.Create(&users[i])
	}

	return users
}
