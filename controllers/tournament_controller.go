package controllers

import (
	"fmt"
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
	var tournament models.Tournament
	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(tournament)

	if err := config.DB.Create(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tournament)
}

// UpdateTournament updates an existing tournament
func UpdateTournament(c *gin.Context) {
	tournamentID := c.Param("tournamentID")
	var tournament models.Tournament

	if err := config.DB.First(&tournament, tournamentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
