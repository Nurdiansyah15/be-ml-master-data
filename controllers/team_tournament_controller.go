package controllers

import (
	"net/http"
	"strconv"

	"ml-master-data/config"
	"ml-master-data/models"

	"github.com/gin-gonic/gin"
)

// CreateTeamInTournament adds a team to a tournament
func CreateTeamInTournament(c *gin.Context) {
	tournamentID, _ := strconv.Atoi(c.Param("tournamentID"))
	var requestBody struct {
		TeamID uint `json:"team_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if tournament exists
	var tournament models.Tournament
	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	match := models.Match{
		TournamentID: uint(tournamentID),
		HomeTeamID:   requestBody.TeamID,
	}
	err := config.DB.Create(&match).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team added to tournament", "tournamentID": tournamentID, "teamID": requestBody.TeamID})
}

func GetAllTeamsInTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	var teams []models.Team

	// Menggunakan query untuk mendapatkan tim yang terlibat dalam turnamen tertentu
	if err := config.DB.Model(&models.Team{}).
		Joins("JOIN matches ON teams.team_id = matches.home_team_id OR teams.team_id = matches.away_team_id").
		Where("matches.tournament_id = ?", tournamentID).
		Distinct().
		Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mengembalikan daftar tim dalam format JSON
	c.JSON(http.StatusOK, teams)
}
