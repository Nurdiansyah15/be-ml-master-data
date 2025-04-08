package seeders

import (
	"ml-master-data/config"
	"ml-master-data/models"

	"golang.org/x/crypto/bcrypt"
)

// Fungsi untuk seeding heroes
func seedUsers() []models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("uw928820wjsnnw"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("amalnudros"), bcrypt.DefaultCost)

	users := []models.User{
		{Username: "heroisgod", Password: string(hashedPassword)},
		{Username: "amalnudros", Password: string(hashedPassword2)},
	}

	for i := range users {
		config.DB.Create(&users[i])
	}

	return users
}
