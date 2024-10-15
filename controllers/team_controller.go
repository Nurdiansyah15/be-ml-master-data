package controllers

import (
	"fmt"
	"ml-master-data/config"
	"ml-master-data/models"
	"ml-master-data/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllTeams(c *gin.Context) {
	var teams []models.Team

	if err := config.DB.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func CreateTeam(c *gin.Context) {

	name := c.Request.FormValue("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	file, err := c.FormFile("logo")
	var logoPath string

	if err != nil {
		logoPath = "https://placehold.co/400x600"
	} else {

		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		newFileName := utils.GenerateUniqueFileName("team") + ext
		logoPath = fmt.Sprintf("public/images/%s", newFileName)

		if err := c.SaveUploadedFile(file, logoPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo file"})
			return
		}

		logoPath = os.Getenv("BASE_URL") + "/" + logoPath
	}

	team := models.Team{
		Name:  name,
		Image: logoPath,
	}

	if err := config.DB.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func UpdateTeam(c *gin.Context) {
	// Ambil parameter teamID dari URL
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	// Cari tim di database
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Tangkap data form dari request (name dan logo jika ada file)
	name := c.Request.FormValue("name")

	// Jika ada perubahan name, update
	if name != "" {
		team.Name = name
	}

	// Tangani file logo baru jika ada
	file, err := c.FormFile("logo")
	if err == nil {

		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Jika ada file logo baru, hapus logo lama
		if team.Image != "" && team.Image != "https://placehold.co/400x600" {
			team.Image = strings.Replace(team.Image, os.Getenv("BASE_URL")+"/", "", 1)
			// Cek apakah file Image lama ada di sistem
			if _, err := os.Stat(team.Image); err == nil {
				// Jika file ada, hapus file Image lama dari folder images
				if err := os.Remove(team.Image); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old image"})
					return
				}
			} else if os.IsNotExist(err) {
				// Jika file tidak ada, berikan pesan peringatan (opsional)
				c.JSON(http.StatusNotFound, gin.H{"warning": "Old image not found, skipping deletion"})
			}
		}

		newFileName := utils.GenerateUniqueFileName("team") + ext
		newLogoPath := fmt.Sprintf("public/images/%s", newFileName)
		if err := c.SaveUploadedFile(file, newLogoPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new logo"})
			return
		}

		// Update path logo di database
		team.Image = os.Getenv("BASE_URL") + "/" + newLogoPath
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, team)
}

func GetTeamByID(c *gin.Context) {
	// Ambil parameter teamID dari URL
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	// Cari tim di database
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Kembalikan data tim dalam format JSON
	c.JSON(http.StatusOK, team)
}

func CreatePlayerInTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	// Tangkap data name dan role dari form-data
	name := c.Request.FormValue("name")
	role := c.Request.FormValue("role")

	if name == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Role are required"})
		return
	}

	// Cari tim di database berdasarkan teamID
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Tangani file gambar jika ada
	file, err := c.FormFile("image")
	var imagePath string
	if err != nil {
		// Jika tidak ada file yang diupload, gunakan placeholder
		imagePath = "https://placehold.co/400x600"
	} else {

		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		newFileName := utils.GenerateUniqueFileName("player") + ext
		imagePath = fmt.Sprintf("public/images/%s", newFileName)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imagePath = os.Getenv("BASE_URL") + "/" + imagePath
	}

	// Buat objek Player
	player := models.Player{
		Name:   name,
		Role:   role,
		Image:  imagePath,
		TeamID: team.TeamID,
	}

	// Simpan player ke database
	if err := config.DB.Create(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusCreated, player)
}

func CreateCoachInTeam(c *gin.Context) {
	// Ambil parameter teamID dari URL
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	// Cari tim di database berdasarkan teamID
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Tangkap data name dan role dari form-data
	name := c.Request.FormValue("name")
	role := c.Request.FormValue("role")

	if name == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Role are required"})
		return
	}

	// Tangani file gambar jika ada
	file, err := c.FormFile("image")
	var imagePath string
	if err != nil {
		// Jika tidak ada file yang diupload, gunakan placeholder
		imagePath = "https://placehold.co/400x600"
	} else {
		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		newFileName := utils.GenerateUniqueFileName("coach") + ext
		imagePath = fmt.Sprintf("public/images/%s", newFileName)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imagePath = os.Getenv("BASE_URL") + "/" + imagePath
	}

	// Buat objek Coach
	coach := models.Coach{
		Name:   name,
		Role:   role,
		Image:  imagePath,
		TeamID: team.TeamID,
	}

	// Simpan coach ke database
	if err := config.DB.Create(&coach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusCreated, coach)
}

func UpdatePlayerInTeam(c *gin.Context) {
	playerID := c.Param("playerID")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID is required"})
		return
	}

	// Mencari pemain berdasarkan playerID
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	// Menangkap data name dan role dari form-data
	name := c.Request.FormValue("name")
	role := c.Request.FormValue("role")

	// Update data jika ada input baru
	if name != "" {
		player.Name = name
	}
	if role != "" {
		player.Role = role
	}

	// Tangani file gambar jika ada
	file, err := c.FormFile("image")
	if err == nil {
		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Cek apakah file image lama ada di sistem dan hapus jika ada
		if player.Image != "" && player.Image != "https://placehold.co/400x600" {
			player.Image = strings.Replace(player.Image, os.Getenv("BASE_URL")+"/", "", 1)
			if _, err := os.Stat(player.Image); err == nil {
				if err := os.Remove(player.Image); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old image"})
					return
				}
			}
		}

		newFileName := utils.GenerateUniqueFileName("player") + ext
		imagePath := fmt.Sprintf("public/images/%s", newFileName)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}

		// Update path image di database
		player.Image = os.Getenv("BASE_URL") + "/" + imagePath
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&player).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, player)
}

func UpdateCoachInTeam(c *gin.Context) {
	// Ambil parameter coachID dari URL
	coachID := c.Param("coachID")
	if coachID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach ID is required"})
		return
	}

	// Cari pelatih di database berdasarkan coachID
	var coach models.Coach
	if err := config.DB.First(&coach, coachID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	// Tangkap data name dan role dari form-data
	name := c.Request.FormValue("name")
	role := c.Request.FormValue("role")

	// Update data jika ada input baru
	if name != "" {
		coach.Name = name
	}
	if role != "" {
		coach.Role = role
	}

	// Tangani file gambar jika ada
	file, err := c.FormFile("image")
	if err == nil {

		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Cek apakah file image lama ada di sistem dan hapus jika ada
		if coach.Image != "" && coach.Image != "https://placehold.co/400x600" {
			coach.Image = strings.Replace(coach.Image, os.Getenv("BASE_URL")+"/", "", 1)
			if _, err := os.Stat(coach.Image); err == nil {
				if err := os.Remove(coach.Image); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old image"})
					return
				}
			}
		}

		newFileName := utils.GenerateUniqueFileName("coach") + ext
		imagePath := fmt.Sprintf("public/images/%s", newFileName)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}

		// Update path image di database
		coach.Image = os.Getenv("BASE_URL") + "/" + imagePath
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&coach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, coach)
}

func GetAllPlayersInTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	var players []models.Player
	if err := config.DB.Where("team_id = ?", teamID).Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

func GetPlayerByID(c *gin.Context) {
	// Ambil parameter playerID dari URL
	playerID := c.Param("playerID")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID is required"})
		return
	}

	// Cari pemain di database
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	// Kembalikan data pemain dalam format JSON
	c.JSON(http.StatusOK, player)
}

func GetAllCoachesInTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	var coaches []models.Coach
	if err := config.DB.Where("team_id = ?", teamID).Find(&coaches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coaches)
}

func GetCoachByID(c *gin.Context) {
	// Ambil parameter coachID dari URL
	coachID := c.Param("coachID")
	if coachID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach ID is required"})
		return
	}

	// Cari pelatih di database
	var coach models.Coach
	if err := config.DB.First(&coach, coachID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	// Kembalikan data pelatih dalam format JSON
	c.JSON(http.StatusOK, coach)
}
