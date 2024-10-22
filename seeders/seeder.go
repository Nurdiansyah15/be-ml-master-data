package seeders

import (
	"log"
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
	// seedUsers()
	// seedHeroes()
	// seedTeams()

	log.Println("Seeding completed successfully!")
}
