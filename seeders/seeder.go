package seeders

import (
	"log"
)

func Seed() {

	// config.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	// config.DB.Exec("DELETE FROM users")
	// config.DB.Exec("DELETE FROM heros")
	// config.DB.Exec("DELETE FROM teams")
	// config.DB.Exec("DELETE FROM tournaments")
	// config.DB.Exec("DELETE FROM teams")
	// config.DB.Exec("DELETE FROM matches")
	// config.DB.Exec("DELETE FROM tournament_teams")
	// config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Seed Users
	seedUsers()
	seedHeroes()
	seedTeams()

	log.Println("Seeding completed successfully!")
}
