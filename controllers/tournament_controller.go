package controllers

import (
	"net/http"

	"ml-master-data/config" // Ganti dengan path yang sesuai untuk package database Anda
	"ml-master-data/dto"
	"ml-master-data/models"
	"ml-master-data/services"

	"github.com/gin-gonic/gin"
)

// GetAllTournaments gets all tournaments
// @Summary Get all tournaments
// @Description Get all tournaments
// @Tags Tournament
// @Security Bearer
// @Produce json
// @Success 200 {array} models.Tournament
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments [get]
func GetAllTournaments(c *gin.Context) {
	var tournaments []models.Tournament

	if err := config.DB.Find(&tournaments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// GetTournamentByID gets a tournament by ID
// @Summary Get a tournament by ID
// @Description Get a tournament by ID
// @Tags Tournament
// @Security Bearer
// @Produce json
// @Param tournamentID path string true "Tournament ID"
// @Success 200 {object} models.Tournament
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Tournament not found"
// @Router /tournaments/{tournamentID} [get]
func GetTournamentByID(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tournament ID is required"})
		return
	}

	var tournament models.Tournament

	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// CreateTournament creates a new tournament
// @Summary Create a new tournament
// @Description Create a new tournament with the given name
// @Tags Tournament
// @Security Bearer
// @Produce json
// @Param dto body dto.TournamentRequestDto true "Tournament request"
// @Success 201 {object} models.Tournament
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments [post]
func CreateTournament(c *gin.Context) {

	input := dto.TournamentRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tournament models.Tournament

	tournament.Name = input.Name

	if err := config.DB.Create(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tournament)
}

// UpdateTournament updates a tournament
// @Summary Update a tournament
// @Description Update a tournament with the given name
// @Tags Tournament
// @Security Bearer
// @Produce json
// @Param tournamentID path string true "Tournament ID"
// @Param dto body dto.TournamentRequestDto true "Tournament request"
// @Success 200 {object} models.Tournament
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Tournament not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments/{tournamentID} [put]
func UpdateTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	if tournamentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tournament ID is required"})
		return
	}

	var tournament models.Tournament

	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	input := dto.TournamentRequestDto{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		tournament.Name = input.Name
	}

	// Update the tournament's name
	if err := config.DB.Save(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// @Summary Delete a tournament
// @Description Delete a tournament with the given tournament ID
// @Tags Tournament
// @Security Bearer
// @Produce json
// @Param tournamentID path string true "Tournament ID"
// @Success 200 {string} string "Tournament deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Tournament not found"
// @Failure 500 {string} string "Internal server error"
// @Router /tournaments/{tournamentID} [delete]
func DeleteTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	var tournament models.Tournament

	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	tournament = models.Tournament{}
	if err := config.DB.Where("tournament_id = ?", tournamentID).First(&tournament).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	if err := services.DeleteTournament(config.DB, uint(tournament.TournamentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tournament deleted successfully"})
}
