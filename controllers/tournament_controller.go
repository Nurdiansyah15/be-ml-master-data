package controllers

import (
	"net/http"

	"ml-master-data/config" // Ganti dengan path yang sesuai untuk package database Anda
	"ml-master-data/models"

	"github.com/gin-gonic/gin"
)

// GetAllTournaments retrieves all tournaments from the database
func GetAllTournaments(c *gin.Context) {
	var tournaments []models.Tournament

	if err := config.DB.Find(&tournaments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// CreateTournament creates a new tournament
func CreateTournament(c *gin.Context) {
	input := struct {
		Name   string `json:"name" binding:"required"`
		Season string `json:"season" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tournament models.Tournament

	tournament.Name = input.Name
	tournament.Season = input.Season

	if err := config.DB.Create(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tournament)
}

// UpdateTournament updates an existing tournament
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

	input := struct {
		Name   string `json:"name"`
		Season string `json:"season"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		tournament.Name = input.Name
	}
	if input.Season != "" {
		tournament.Season = input.Season
	}

	// Update the tournament's name
	if err := config.DB.Save(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// DeleteTournament deletes a tournament
func DeleteTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	var tournament models.Tournament

	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	if err := config.DB.Delete(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tournament deleted successfully"})
}
