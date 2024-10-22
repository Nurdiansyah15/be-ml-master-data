package seeders

import (
	"ml-master-data/config"
	"ml-master-data/models"

	"golang.org/x/crypto/bcrypt"
)

// Fungsi untuk seeding heroes
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
