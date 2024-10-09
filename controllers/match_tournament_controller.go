package controllers

import (
	"ml-master-data/config"
	"ml-master-data/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTeamMatchesinTournament(c *gin.Context) {

	tournamentID := c.Param("tournamentID")
	teamID := c.Param("teamID")

	tournamentTeam := models.TournamentTeam{}

	if err := config.DB.Where("tournament_id = ? AND team_id = ?", tournamentID, teamID).First(&tournamentTeam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var matches []models.Match
	if err := config.DB.Model(&models.Match{}).Where("tournament_team_id = ?", tournamentTeam.TournamentID).Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matches)
}

func CreateTeamMatchinTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	teamID := c.Param("teamID")

	tournamentTeam := models.TournamentTeam{}

	if err := config.DB.Where("tournament_id = ? AND team_id = ?", tournamentID, teamID).First(&tournamentTeam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var match models.Match

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match.TournamentTeamID = tournamentTeam.TournamentID

	if err := config.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

func UpdateTeamMatchinTournament(c *gin.Context) {
	matchID := c.Param("matchID")

	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}
