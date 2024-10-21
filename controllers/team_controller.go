package controllers

import (
	"fmt"
	"ml-master-data/config"
	"ml-master-data/models"
	"ml-master-data/services"
	"ml-master-data/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Get all teams
// @Description Get all teams
// @Produce json
// @Tags Team
// @Security Bearer
// @Success 200 {array} models.Team
// @Failure 500 {string} string "Internal server error"
// @Router /teams [get]
func GetAllTeams(c *gin.Context) {
	var teams []models.Team

	if err := config.DB.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// @Summary Create a team
// @Description Create a team and save its logo
// @Produce json
// @Tags Team
// @Security Bearer
// @Param name formData string true "Team name"
// @Param image formData file true "Team logo"
// @Success 201 {object} models.Team
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /teams [post]
func CreateTeam(c *gin.Context) {

	name := c.Request.FormValue("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	file, err := c.FormFile("image")
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

// @Summary Update a team
// @Description Update a team and save its logo
// @Produce json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param name formData string false "Team name"
// @Param image formData file false "Team logo"
// @Success 200 {object} models.Team
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /teams/{teamID} [put]
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

// DeleteTeam godoc
// @Summary Delete a team
// @Description Delete a team and all related data
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Success 200 {object} string "Team deleted successfully"
// @Failure 400 {string} string "Team ID is required" or "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /teams/{teamID} [delete]
func DeleteTeam(c *gin.Context) {
	// Ambil parameter teamID dari URL
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	team := models.Team{}
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	if err := services.DeleteTeam(config.DB, team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

// @Summary Get a team by ID
// @Description Get a team by ID
// @Produce json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Success 200 {object} models.Team
// @Failure 400 {string} string "Team ID is required"
// @Failure 404 {string} string "Team not found"
// @Router /teams/{teamID} [get]
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

// @Summary Create a player in a team
// @Description Create a player in a team by ID and save its image
// @Produce json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param name formData string true "Player name"
// @Param image formData file false "Player image"
// @Success 201 {object} models.Player
// @Failure 400 {string} string "Team ID is required" or "Name and Role are required" or "File size must not exceed 500 KB" or "Invalid file type"
// @Failure 404 {string} string "Team not found"
// @Router /teams/{teamID}/players [post]
func CreatePlayerInTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	// Tangkap data name dan role dari form-data
	name := c.Request.FormValue("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name are required"})
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

// CreateCoachInTeam godoc
// @Summary Create a coach in a team
// @Description Create a coach in a team and save its image
// @Produce json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param name formData string true "Coach name"
// @Param image formData file false "Coach image"
// @Success 201 {object} models.Coach
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /teams/{teamID}/coaches [post]
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

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name are required"})
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

// @Summary Update a player in a team
// @Description Update a player in a team and save its image
// @Produce json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Player ID"
// @Param name formData string false "Player name"
// @Param image formData file false "Player image"
// @Success 200 {object} models.Player
// @Failure 400 {string} string "Player ID is required" or "File size must not exceed 500 KB" or "Invalid file type"
// @Failure 404 {string} string "Player not found"
// @Router /players/{teamID} [put]
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

	// Update data jika ada input baru
	if name != "" {
		player.Name = name
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

// @Summary Delete a player in a team
// @Description Delete a player in a team and all its related data
// @Produce json
// @Tags Team
// @Security Bearer
// @Param playerID path string true "Player ID"
// @Success 200 {string} string "Player deleted successfully"
// @Failure 400 {string} string "Player ID is required"
// @Failure 404 {string} string "Player not found"
// @Failure 500 {string} string "Internal server error"
// @Router /players/{playerID} [delete]
func DeletePlayerInTeam(c *gin.Context) {

	playerID := c.Param("playerID")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID is required"})
		return
	}

	player := models.Player{}

	if err := config.DB.First(&player, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	if err := services.DeletePlayer(config.DB, player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})

}

// @Summary Update a coach in a team
// @Description Update a coach in a team and save its image
// @Produce json
// @Tags Team
// @Security Bearer
// @Param coachID path string true "Coach ID"
// @Param name formData string false "Coach name"
// @Param image formData file false "Coach image"
// @Success 200 {object} models.Coach
// @Failure 400 {string} string "Coach ID is required" or "File size must not exceed 500 KB" or "Invalid file type"
// @Failure 404 {string} string "Coach not found"
// @Router /coaches/{coachID} [put]
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

	// Update data jika ada input baru
	if name != "" {
		coach.Name = name
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

// @Summary Delete a coach in a team
// @Description Delete a coach in a team and all its related data
// @Produce json
// @Tags Team
// @Security Bearer
// @Param coachID path string true "Coach ID"
// @Success 200 {string} string "Coach deleted successfully"
// @Failure 400 {string} string "Coach ID is required"
// @Failure 404 {string} string "Coach not found"
// @Failure 500 {string} string "Internal server error"
// @Router /coaches/{coachID} [delete]
func DeleteCoachInTeam(c *gin.Context) {
	// Ambil parameter coachID dari URL
	coachID := c.Param("coachID")

	if coachID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach ID is required"})
		return
	}

	coach := models.Coach{}

	if err := config.DB.First(&coach, coachID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	if err := services.DeleteCoach(config.DB, coach); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coach deleted successfully"})
}

type PlayerStats struct {
	TotalMatch    int `json:"total_match"`
	TotalMatchWin int `json:"total_match_win"`
	TotalGame     int `json:"total_game"`
	TotalGameWin  int `json:"total_game_win"`
}

// @Summary Get player statistics
// @Description Get player statistics with the given player ID and tournament ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param playerID path string true "Player ID"
// @Param tournamentID path string true "Tournament ID"
// @Success 200 {object} PlayerStats
// @Failure 400 {string} string "Player ID and Tournament ID are required" or "Invalid Player ID format" or "Invalid Tournament ID format"
// @Failure 500 {string} string "Internal server error"
// @Router /players/{playerID}/tournaments/{tournamentID}/player-statistics [get]
func PlayerStatistics(c *gin.Context) {
	playerIDStr := c.Param("playerID")
	tournamentIDStr := c.Param("tournamentID")

	if playerIDStr == "" || tournamentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID and Tournament ID are required"})
		return
	}

	playerID, err := strconv.ParseUint(playerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Player ID format"})
		return
	}

	tournamentID, err := strconv.ParseUint(tournamentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID format"})
		return
	}

	matchStats := struct {
		TotalMatch    int `json:"total_match"`
		TotalMatchWin int `json:"total_match_win"`
	}{}
	// Query untuk total_match dan total_match_win
	matchQuery := `
		SELECT 
			COUNT(DISTINCT m.match_id) as total_match,
			SUM(CASE 
				WHEN (m.team_a_id = t.team_id AND m.team_a_score > m.team_b_score) OR 
					 (m.team_b_id = t.team_id AND m.team_b_score > m.team_a_score) 
				THEN 1 ELSE 0 END) as total_match_win
		FROM players p
		JOIN teams t ON p.team_id = t.team_id
		JOIN player_matches pm ON p.player_id = pm.player_id
		JOIN match_team_details mtd ON pm.match_team_detail_id = mtd.match_team_detail_id
		JOIN matches m ON mtd.match_id = m.match_id
		WHERE p.player_id = ? AND m.tournament_id = ?
	`
	if err := config.DB.Raw(matchQuery, playerID, tournamentID).Scan(&matchStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying match statistics: " + err.Error()})
		return
	}

	gameStats := struct {
		TotalGame    int `json:"total_game"`
		TotalGameWin int `json:"total_game_win"`
	}{}

	// Query untuk total_game dan total_game_win
	gameQuery := `
		SELECT 
			COUNT(DISTINCT g.game_id) as total_game,
			SUM(CASE WHEN g.winner_team_id = t.team_id THEN 1 ELSE 0 END) as total_game_win
		FROM players p
		JOIN teams t ON p.team_id = t.team_id
		JOIN player_matches pm ON p.player_id = pm.player_id
		JOIN match_team_details mtd ON pm.match_team_detail_id = mtd.match_team_detail_id
		JOIN matches m ON mtd.match_id = m.match_id
		JOIN games g ON m.match_id = g.match_id
		WHERE p.player_id = ? AND m.tournament_id = ?
	`
	if err := config.DB.Raw(gameQuery, playerID, tournamentID).Scan(&gameStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying game statistics: " + err.Error()})
		return
	}

	var stats = PlayerStats{
		TotalMatch:    matchStats.TotalMatch,
		TotalMatchWin: matchStats.TotalMatchWin,
		TotalGame:     gameStats.TotalGame,
		TotalGameWin:  gameStats.TotalGameWin,
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Get all players in a team
// @Description Get all players in a team with the given team ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Success 200 {array} models.Player
// @Failure 400 {string} string "Team ID is required"
// @Failure 404 {string} string "Team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /teams/{teamID}/players [get]
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

// @Summary Get a player by ID
// @Description Get a player by ID with the given player ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param playerID path string true "Player ID"
// @Success 200 {object} models.Player
// @Failure 400 {string} string "Player ID is required"
// @Failure 404 {string} string "Player not found"
// @Failure 500 {string} string "Internal server error"
// @Router /players/{playerID} [get]
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

type CoachStats struct {
	TotalMatch    int `json:"total_match"`
	TotalMatchWin int `json:"total_match_win"`
	TotalGame     int `json:"total_game"`
	TotalGameWin  int `json:"total_game_win"`
}

// @Summary Get coach statistics
// @Description Get coach statistics with the given coach ID and tournament ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param coachID path string true "Coach ID"
// @Param tournamentID path string true "Tournament ID"
// @Success 200 {object} CoachStats
// @Failure 400 {string} string "Coach ID and Tournament ID are required"
// @Failure 404 {string} string "Coach not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments/{tournamentID}/coachs/{coachID}/coach-statistics [get]
func CoachStatistics(c *gin.Context) {
	coachIDStr := c.Param("coachID")
	tournamentIDStr := c.Param("tournamentID")

	if coachIDStr == "" || tournamentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach ID and Tournament ID are required"})
		return
	}

	coachID, err := strconv.ParseUint(coachIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Coach ID format"})
		return
	}

	tournamentID, err := strconv.ParseUint(tournamentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID format"})
		return
	}

	matchStats := struct {
		TotalMatch    int `json:"total_match"`
		TotalMatchWin int `json:"total_match_win"`
	}{}
	// Query untuk total_match dan total_match_win
	matchQuery := `
	SELECT 
		COUNT(DISTINCT m.match_id) as total_match,
		SUM(CASE 
			WHEN (m.team_a_id = t.team_id AND m.team_a_score > m.team_b_score) OR 
				 (m.team_b_id = t.team_id AND m.team_b_score > m.team_a_score) 
			THEN 1 ELSE 0 END) as total_match_win
	FROM coaches p
	JOIN teams t ON p.team_id = t.team_id
	JOIN coach_matches pm ON p.coach_id = pm.coach_id
	JOIN match_team_details mtd ON pm.match_team_detail_id = mtd.match_team_detail_id
	JOIN matches m ON mtd.match_id = m.match_id
	WHERE p.coach_id = ? AND m.tournament_id = ?
`
	if err := config.DB.Raw(matchQuery, coachID, tournamentID).Scan(&matchStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying match statistics: " + err.Error()})
		return
	}

	gameStats := struct {
		TotalGame    int `json:"total_game"`
		TotalGameWin int `json:"total_game_win"`
	}{}

	// Query untuk total_game dan total_game_win
	gameQuery := `
			SELECT 
				COUNT(DISTINCT g.game_id) as total_game,
				SUM(CASE WHEN g.winner_team_id = t.team_id THEN 1 ELSE 0 END) as total_game_win
			FROM coaches p
			JOIN teams t ON p.team_id = t.team_id
			JOIN coach_matches pm ON p.coach_id = pm.coach_id
			JOIN match_team_details mtd ON pm.match_team_detail_id = mtd.match_team_detail_id
			JOIN matches m ON mtd.match_id = m.match_id
			JOIN games g ON m.match_id = g.match_id
			WHERE p.coach_id = ? AND m.tournament_id = ?
		`
	if err := config.DB.Raw(gameQuery, coachID, tournamentID).Scan(&gameStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying game statistics: " + err.Error()})
		return
	}

	var stats = CoachStats{
		TotalMatch:    matchStats.TotalMatch,
		TotalMatchWin: matchStats.TotalMatchWin,
		TotalGame:     gameStats.TotalGame,
		TotalGameWin:  gameStats.TotalGameWin,
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Get all coaches in a team
// @Description Get all coaches in a team with the given team ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Success 200 {array} models.Coach
// @Failure 400 {string} string "Team ID is required"
// @Failure 404 {string} string "Team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /teams/{teamID}/coaches [get]
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

// @Summary Get a coach by ID
// @Description Get a coach by ID with the given coach ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param coachID path string true "Coach ID"
// @Success 200 {object} models.Coach
// @Failure 400 {string} string "Coach ID is required"
// @Failure 404 {string} string "Coach not found"
// @Failure 500 {string} string "Internal server error"
// @Router /coaches/{coachID} [get]
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

type TeamStatisticsDto struct {
	TeamID                 uint `json:"teamID"`
	TotalMatch             int  `json:"totalMatch"`
	TotalMatchAndWin       int  `json:"totalMatchAndWin"`
	TotalMatchAndLose      int  `json:"totalMatchAndLose"`
	TotalGame              int  `json:"totalGame"`
	TotalGameAndWin        int  `json:"totalGameAndWin"`
	TotalGameAndLose       int  `json:"totalGameAndLose"`
	TotalFirstPick         int  `json:"totalFirstPick"`
	TotalFirstPickAndWin   int  `json:"totalFirstPickAndWin"`
	TotalFirstPickAndLose  int  `json:"totalFirstPickAndLose"`
	TotalSecondPick        int  `json:"totalSecondPick"`
	TotalSecondPickAndWin  int  `json:"totalSecondPickAndWin"`
	TotalSecondPickAndLose int  `json:"totalSecondPickAndLose"`
}

// @Summary Get team statistics
// @Description Get team statistics with the given team ID
// @Accept  json
// @Produce  json
// @Tags Team
// @Security Bearer
// @Param teamID path string true "Team ID"
// @Param tournamentID path string true "Tournament ID"
// @Success 200 {object} TeamStatisticsDto
// @Failure 400 {string} string "Team ID is required" or "Invalid Team ID format"
// @Failure 404 {string} string "Team not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments/{tournamentID}/teams/{teamID}/team-statistics [get]
func GetTeamStatistics(c *gin.Context) {
	teamIDStr := c.Param("teamID")
	tournamentIDStr := c.Param("tournamentID")

	if teamIDStr == "" || tournamentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID and Tournament ID are required"})
		return
	}

	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Team ID format"})
		return
	}

	tournamentID, err := strconv.ParseUint(tournamentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID format"})
		return
	}

	stats := TeamStatisticsDto{}
	stats.TeamID = uint(teamID)

	// Query matches for the specific team and tournament
	var matches []models.Match
	if err := config.DB.Where("tournament_id = ? AND (team_a_id = ? OR team_b_id = ?)", tournamentID, teamID, teamID).Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying matches: " + err.Error()})
		return
	}

	stats.TotalMatch = len(matches)

	// Calculate match wins and losses from match scores
	matchIDs := make([]uint, 0, len(matches))
	for _, match := range matches {
		matchIDs = append(matchIDs, match.MatchID)
		if match.TeamAID == uint(teamID) {
			if match.TeamAScore > match.TeamBScore {
				stats.TotalMatchAndWin++
			} else if match.TeamAScore < match.TeamBScore {
				stats.TotalMatchAndLose++
			}
		} else if match.TeamBID == uint(teamID) {
			if match.TeamBScore > match.TeamAScore {
				stats.TotalMatchAndWin++
			} else if match.TeamBScore < match.TeamAScore {
				stats.TotalMatchAndLose++
			}
		}
	}

	// Query games for the filtered matches and specific team
	var games []models.Game
	if err := config.DB.Where("match_id IN ? AND (first_pick_team_id = ? OR second_pick_team_id = ?)", matchIDs, teamID, teamID).Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying games: " + err.Error()})
		return
	}

	stats.TotalGame = len(games)

	// Calculate game-specific statistics
	for _, game := range games {
		if game.FirstPickTeamID == uint(teamID) {
			stats.TotalFirstPick++
			if game.WinnerTeamID == uint(teamID) {
				stats.TotalGameAndWin++
				stats.TotalFirstPickAndWin++
			} else {
				stats.TotalGameAndLose++
				stats.TotalFirstPickAndLose++
			}
		} else if game.SecondPickTeamID == uint(teamID) {
			stats.TotalSecondPick++
			if game.WinnerTeamID == uint(teamID) {
				stats.TotalGameAndWin++
				stats.TotalSecondPickAndWin++
			} else {
				stats.TotalGameAndLose++
				stats.TotalSecondPickAndLose++
			}
		}
	}

	c.JSON(http.StatusOK, stats)
}
